package arcgis

import (
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mocheer/charon/src/constants"
	"github.com/mocheer/charon/src/models/types"
)

//NewTileServer returns a new Esri, based on a conf.xml path
func NewTileServer(confPath string) (*TileServer, error) {
	server := &TileServer{}
	confXML, err := ioutil.ReadFile(confPath)
	if err != nil {
		return nil, err
	}
	var config TileLayerConfig
	err = xml.Unmarshal(confXML, &config)
	if err != nil {
		return nil, err
	}
	server.BaseDirectory = filepath.Dir(confPath)
	server.MinLevel, server.MaxLevel = calcMinMaxLevels(&config, server.BaseDirectory)
	server.FileFormat = config.TileImageInfo.CacheTileFormat
	server.CacheFormat = config.CacheStorageInfo.StorageFormat
	server.HasTransparency = (server.FileFormat == "PNG" || server.FileFormat == "PNG32" || server.FileFormat == "MIXED")
	server.EpsgCode = config.TileCacheInfo.SpatialReference.WKID
	server.TileColumnSize = config.TileCacheInfo.TileCols
	server.TileRowSize = config.TileCacheInfo.TileRows
	packetSize := config.CacheStorageInfo.PacketSize
	if packetSize != nil {
		server.ColsPerFile, server.RowsPerFile = *packetSize, *packetSize
	} else {
		server.ColsPerFile, server.RowsPerFile = 1, 1
	}
	return server, nil
}

//calcMinMaxLevels is called by NewEsri to return min and max levels
func calcMinMaxLevels(cache *TileLayerConfig, baseDir string) (int, int) {
	minLevel := int(^uint(0) >> 1)
	maxLevel := 0
	for _, li := range cache.TileCacheInfo.LODInfos.LODInfo {
		levelPath := filepath.Join(baseDir, "_alllayers", fmt.Sprintf("L%02d", li.LevelID))
		if _, err := os.Stat(levelPath); err != nil {
			continue
		}
		if li.LevelID > maxLevel {
			maxLevel = li.LevelID
		}
		if li.LevelID < minLevel {
			minLevel = li.LevelID
		}
	}
	if minLevel > maxLevel {
		minLevel = maxLevel
	}
	return minLevel, maxLevel
}

// ReadTile returns a 256x256 tile
func (server *TileServer) ReadTile(tile types.Tile) ([]byte, error) {
	switch server.CacheFormat {
	case "esriMapCacheStorageModeCompact":
		return server.ReadCompactTile(tile)
	case "esriMapCacheStorageModeCompactV2":
		return server.ReadCompactTileV2(tile)
	default:
		return server.ReadExplodedTile(tile)
	}
}

// ReadCompactTile returns a bundled 256x256 tile
func (server *TileServer) ReadCompactTile(tile types.Tile) ([]byte, error) {
	bundlxPath, bundlePath, imgDataIndex := server.GetFileInfo(tile)
	bundlx, err := os.Open(bundlxPath)
	if err != nil {
		return nil, err
	}
	defer bundlx.Close()
	bundlx.Seek((16 + (5 * imgDataIndex)), io.SeekStart)
	bOffset := make([]byte, 5, 5)
	bundlx.Read(bOffset)
	offset := int64(binary.LittleEndian.Uint64(bOffset))
	bundle, err := os.Open(bundlePath)
	if err != nil {
		return nil, err
	}
	defer bundle.Close()
	bundle.Seek(offset, io.SeekStart)
	bLength := make([]byte, 4, 4)
	bundle.Read(bLength)
	length := binary.LittleEndian.Uint64(bLength)
	imgBytes := make([]byte, length, length)
	bundle.Read(imgBytes)
	return imgBytes, nil
}

//ReadCompactTileV2 returns a bundled 256x256 tile
func (server *TileServer) ReadCompactTileV2(tile types.Tile) ([]byte, error) {
	_, bundlePath, _ := server.GetFileInfo(tile)
	BundlxMaxidx := constants.BundlxMaxidx
	CompactCacheHeaderLength := constants.CompactCacheHeaderLength

	// col and row are inverted for 10.3 caches
	index := BundlxMaxidx*(tile.Row%BundlxMaxidx) + (tile.Column % BundlxMaxidx)
	offset := (index * 8) + CompactCacheHeaderLength

	bundle, err := os.Open(bundlePath)
	if err != nil {
		return nil, err
	}
	defer bundle.Close()
	bundle.Seek(int64(offset), io.SeekStart)

	offsetBytes := make([]byte, 5, 8) //4,4
	sizeBytes := make([]byte, 3, 4)   //4,4

	bundle.Read(offsetBytes)
	bundle.Read(sizeBytes)

	offsetBytes = offsetBytes[:8]
	sizeBytes = sizeBytes[:4]

	dataOffset := binary.LittleEndian.Uint64(offsetBytes)

	size := binary.LittleEndian.Uint32(sizeBytes)

	imgBytes := make([]byte, size, size)
	bundle.Seek(int64(dataOffset), io.SeekStart)
	bundle.Read(imgBytes)
	return imgBytes, nil
}

//GetFileInfo returns file paths and indexes into those files
func (server *TileServer) GetFileInfo(tile types.Tile) (bundlxPath, bundlePath string, imgDataIndex int64) {
	internalRow := tile.Row % server.RowsPerFile
	internalCol := tile.Column % server.ColsPerFile
	bundleRow := tile.Row - internalRow
	bundleCol := tile.Column - internalCol
	bundleBasePath := filepath.Join(server.BaseDirectory, "_alllayers", fmt.Sprintf("L%02d", tile.Level), fmt.Sprintf("R%04xC%04x", bundleRow, bundleCol))
	bundlxPath = bundleBasePath + ".bundlx"
	bundlePath = bundleBasePath + ".bundle"
	imgDataIndex = int64((server.ColsPerFile * internalCol) + internalRow)
	return bundlxPath, bundlePath, imgDataIndex
}

//ReadExplodedTile returns a standalone 256x256 tile
func (server *TileServer) ReadExplodedTile(tile types.Tile) ([]byte, error) {
	return ioutil.ReadFile(server.GetFilePath(tile))
}

//GetFilePath return the primary file path, sans extension
func (server *TileServer) GetFilePath(tile types.Tile) string {
	level := fmt.Sprintf("L%02d", tile.Level)
	row := fmt.Sprintf("R%08x", tile.Row)
	column := fmt.Sprintf("C%08x", tile.Column)
	filePath := filepath.Join(server.BaseDirectory, level, row, column)

	if server.FileFormat == "JPEG" {
		return filePath + ".jpg" //JPEG
	}
	if server.FileFormat != "MIXED" {
		return filePath + ".png" //PNG, PNG8, PNG24, PNG32
	}
	if _, err := os.Stat(filePath + ".jpg"); err == nil {
		return filePath + ".jpg" //MIXED...
	}
	if _, err := os.Stat(filePath + ".png"); err == nil {
		return filePath + ".png"
	}
	return filePath
}
