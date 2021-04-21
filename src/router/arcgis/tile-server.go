package arcgis

import (
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mocheer/charon/src/models/types"
	"github.com/mocheer/pluto/ts"
)

// NewTileServer 根据config.xml实例化服务
func NewTileServer(confPath string) (*TileServer, error) {
	server := &TileServer{}
	confXML, err := os.ReadFile(confPath)
	if err != nil {
		return nil, err
	}
	var config ts.ArcgisTileLayerConfig
	err = xml.Unmarshal(confXML, &config)
	if err != nil {
		return nil, err
	}
	server.BaseDirectory = filepath.Dir(confPath)
	server.TileFormat = strings.ToLower(config.TileImageInfo.CacheTileFormat) //转成小写是因为fiber目前只支持小写
	server.CacheFormat = config.CacheStorageInfo.StorageFormat
	server.WKID = config.TileCacheInfo.SpatialReference.WKID
	server.TileColSize = config.TileCacheInfo.TileCols
	server.TileRowSize = config.TileCacheInfo.TileRows
	packetSize := config.CacheStorageInfo.PacketSize
	if packetSize != nil {
		server.ColsPerFile, server.RowsPerFile = *packetSize, *packetSize
	} else {
		server.ColsPerFile, server.RowsPerFile = 1, 1
	}
	return server, nil
}

// ReadTile 返回瓦片数据
func (server *TileServer) ReadTile(tile types.Tile) ([]byte, error) {
	switch server.CacheFormat {
	case EsriMapCacheStorageModeCompactV2:
		return server.ReadCompactTileV2(tile)
	case EsriMapCacheStorageModeCompact:
		return server.ReadCompactTile(tile)
	default:
		return server.ReadExplodedTile(tile)
	}
}

// ReadCompactTile 返回紧凑型的切片数据
func (server *TileServer) ReadCompactTile(tile types.Tile) ([]byte, error) {
	bundlxPath, bundlePath, imgDataIndex := server.GetFileInfo(tile)
	bundlx, err := os.Open(bundlxPath)
	if err != nil {
		return nil, err
	}
	defer bundlx.Close()
	bundlx.Seek((16 + (5 * imgDataIndex)), io.SeekStart)
	bOffset := make([]byte, 5)
	bundlx.Read(bOffset)
	offset := int64(binary.LittleEndian.Uint64(bOffset))
	bundle, err := os.Open(bundlePath)
	if err != nil {
		return nil, err
	}
	defer bundle.Close()
	bundle.Seek(offset, io.SeekStart)
	bLength := make([]byte, 4)
	bundle.Read(bLength)
	length := binary.LittleEndian.Uint64(bLength)
	imgBytes := make([]byte, length)
	bundle.Read(imgBytes)
	return imgBytes, nil
}

// ReadCompactTileV2 返回紧凑型V2的切片数据
func (server *TileServer) ReadCompactTileV2(tile types.Tile) ([]byte, error) {
	_, bundlePath, _ := server.GetFileInfo(tile)
	BundlxMaxidx := BundlxMaxidx

	// col and row are inverted for 10.3 caches
	index := BundlxMaxidx*(tile.Y%BundlxMaxidx) + (tile.X % BundlxMaxidx)
	offset := (index * 8) + CompactCacheHeaderLength

	bundle, err := os.Open(bundlePath)
	if err != nil {
		return nil, err
	}
	defer bundle.Close()
	bundle.Seek(int64(offset), io.SeekStart)

	offsetBytes := make([]byte, 5, 8)
	sizeBytes := make([]byte, 3, 4)

	bundle.Read(offsetBytes)
	bundle.Read(sizeBytes)

	offsetBytes = offsetBytes[:8]
	sizeBytes = sizeBytes[:4]

	dataOffset := binary.LittleEndian.Uint64(offsetBytes)
	size := binary.LittleEndian.Uint32(sizeBytes)

	imgBytes := make([]byte, size)
	bundle.Seek(int64(dataOffset), io.SeekStart)
	bundle.Read(imgBytes)

	return imgBytes, nil
}

// GetFileInfo 返回文件路径和数据索引
func (server *TileServer) GetFileInfo(tile types.Tile) (bundlxPath, bundlePath string, imgDataIndex int64) {
	internalRow := tile.Y % server.RowsPerFile
	internalCol := tile.X % server.ColsPerFile
	bundleRow := tile.Y - internalRow
	bundleCol := tile.X - internalCol
	bundleBasePath := filepath.Join(server.BaseDirectory, "_alllayers", fmt.Sprintf("L%02d", tile.Z), fmt.Sprintf("R%04xC%04x", bundleRow, bundleCol))
	bundlxPath = bundleBasePath + ".bundlx"
	bundlePath = bundleBasePath + ".bundle"
	imgDataIndex = int64((server.ColsPerFile * internalCol) + internalRow)
	return bundlxPath, bundlePath, imgDataIndex
}

// ReadExplodedTile 返回单张瓦片
func (server *TileServer) ReadExplodedTile(tile types.Tile) ([]byte, error) {
	return os.ReadFile(server.GetFilePath(tile))
}

// GetFilePath 返回图片路径
func (server *TileServer) GetFilePath(tile types.Tile) string {
	level := fmt.Sprintf("L%02d", tile.Z)
	row := fmt.Sprintf("R%08x", tile.Y)
	column := fmt.Sprintf("C%08x", tile.X)
	filePath := filepath.Join(server.BaseDirectory, level, row, column)

	if server.TileFormat == "JPEG" {
		return filePath + ".jpg" //JPEG
	}
	if server.TileFormat != "MIXED" {
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
