package gis

import (
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mocheer/charon/src/models/types"
)

// MapServer implements TileCache for ESRI local files
// @see https://github.com/wthorp/AGES/tree/master/pkg/sources/tilecache
// @see https://github.com/fuzhenn/tiler-arcgis-bundle/blob/master/index.js
type MapServer struct {
	CacheFormat   string
	BaseDirectory string
	FileFormat    string
	types.TileCache
}

//CacheInfo corresponds to an ESRI conf.xml document
type CacheInfo struct {
	TileCacheInfo struct {
		LODInfos struct {
			LODInfo []struct {
				LevelID int
			}
		}
		SpatialReference struct {
			WKID int
		}
		TileCols int
		TileRows int
	}
	TileImageInfo struct {
		CacheTileFormat string
	}
	CacheStorageInfo struct {
		StorageFormat string
		PacketSize    *int
	}
}

//NewMapServer returns a new Esri, based on a conf.xml path
func NewMapServer(confPath string) (*MapServer, error) {
	tc := &MapServer{}
	confXML, err := ioutil.ReadFile(confPath)
	if err != nil {
		return nil, err
	}
	var cache CacheInfo
	err = xml.Unmarshal(confXML, &cache)
	if err != nil {
		return nil, err
	}
	tc.BaseDirectory = filepath.Dir(confPath)
	tc.MinLevel, tc.MaxLevel = calcMinMaxLevels(&cache, tc.BaseDirectory)
	tc.FileFormat = cache.TileImageInfo.CacheTileFormat
	tc.CacheFormat = cache.CacheStorageInfo.StorageFormat
	packetSize := cache.CacheStorageInfo.PacketSize
	tc.HasTransparency = (tc.FileFormat == "PNG" || tc.FileFormat == "PNG32" || tc.FileFormat == "MIXED")
	tc.EpsgCode = cache.TileCacheInfo.SpatialReference.WKID
	tc.TileColumnSize = cache.TileCacheInfo.TileCols
	tc.TileRowSize = cache.TileCacheInfo.TileRows
	if packetSize != nil {
		tc.ColsPerFile, tc.RowsPerFile = *packetSize, *packetSize
	} else {
		tc.ColsPerFile, tc.RowsPerFile = 1, 1
	}
	return tc, nil
}

//calcMinMaxLevels is called by NewEsri to return min and max levels
func calcMinMaxLevels(cache *CacheInfo, baseDir string) (int, int) {
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

//ReadTile returns a 256x256 tile
func (tc *MapServer) ReadTile(tile types.Tile) ([]byte, error) {
	switch tc.CacheFormat {
	case "esriMapCacheStorageModeCompact":
		return tc.ReadCompactTile(tile)
	case "esriMapCacheStorageModeCompactV2":
		return tc.ReadCompactTileV2(tile)
	default:
		return tc.ReadExplodedTile(tile)
	}
}

//ReadCompactTile returns a bundled 256x256 tile
func (tc *MapServer) ReadCompactTile(tile types.Tile) ([]byte, error) {
	bundlxPath, bundlePath, imgDataIndex := tc.GetFileInfo(tile)
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
func (tc *MapServer) ReadCompactTileV2(tile types.Tile) ([]byte, error) {
	_, bundlePath, _ := tc.GetFileInfo(tile)

	var BUNDLX_MAXIDX = 128
	var COMPACT_CACHE_HEADER_LENGTH = 64
	// col and row are inverted for 10.3 caches
	var index = BUNDLX_MAXIDX*(tile.Row%BUNDLX_MAXIDX) + (tile.Column % BUNDLX_MAXIDX)

	var offset = (index * 8) + COMPACT_CACHE_HEADER_LENGTH

	fmt.Println(index, offset, int64(offset))

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
	fmt.Println(offsetBytes)

	size := binary.LittleEndian.Uint32(sizeBytes)

	imgBytes := make([]byte, size, size)
	bundle.Seek(int64(dataOffset), io.SeekStart)
	bundle.Read(imgBytes)
	return imgBytes, nil
}

//GetFileInfo returns file paths and indexes into those files
func (tc *MapServer) GetFileInfo(tile types.Tile) (bundlxPath, bundlePath string, imgDataIndex int64) {
	internalRow := tile.Row % tc.RowsPerFile
	internalCol := tile.Column % tc.ColsPerFile
	bundleRow := tile.Row - internalRow
	bundleCol := tile.Column - internalCol
	bundleBasePath := filepath.Join(tc.BaseDirectory, "_alllayers", fmt.Sprintf("L%02d", tile.Level), fmt.Sprintf("R%04xC%04x", bundleRow, bundleCol))
	bundlxPath = bundleBasePath + ".bundlx"
	bundlePath = bundleBasePath + ".bundle"
	imgDataIndex = int64((tc.ColsPerFile * internalCol) + internalRow)
	return bundlxPath, bundlePath, imgDataIndex
}

//ReadExplodedTile returns a standalone 256x256 tile
func (tc *MapServer) ReadExplodedTile(tile types.Tile) ([]byte, error) {
	return ioutil.ReadFile(tc.GetFilePath(tile))
}

//GetFilePath return the primary file path, sans extension
func (tc *MapServer) GetFilePath(tile types.Tile) string {
	level := fmt.Sprintf("L%02d", tile.Level)
	row := fmt.Sprintf("R%08x", tile.Row)
	column := fmt.Sprintf("C%08x", tile.Column)
	filePath := filepath.Join(tc.BaseDirectory, level, row, column)
	fmt.Println(tc.FileFormat)
	if tc.FileFormat == "JPEG" {
		return filePath + ".jpg" //JPEG
	}
	if tc.FileFormat != "MIXED" {
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
