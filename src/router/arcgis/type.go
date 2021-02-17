package arcgis

// TileServer implements TileCache for ESRI local files
// @see https://github.com/wthorp/AGES/tree/master/pkg/sources/tilecache
// @see https://github.com/fuzhenn/tiler-arcgis-bundle/blob/master/index.js
type TileServer struct {
	BaseDirectory string // 根目录
	CacheFormat   string // 切片缓存模式
	TileFormat    string // 切片数据类型
	TileColSize   int    // 切片column大小
	TileRowSize   int    // 切片 row 大小
	ColsPerFile   int    //
	RowsPerFile   int    //
	WKID          int    // 切片服务wkid
}

//TileLayerConfig corresponds to an ESRI conf.xml document
type TileLayerConfig struct {
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
