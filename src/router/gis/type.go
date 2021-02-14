package gis

import "github.com/mocheer/charon/src/models/types"

// MapServer implements TileCache for ESRI local files
// @see https://github.com/wthorp/AGES/tree/master/pkg/sources/tilecache
// @see https://github.com/fuzhenn/tiler-arcgis-bundle/blob/master/index.js
type MapServer struct {
	CacheFormat   string
	BaseDirectory string
	FileFormat    string
	types.TileCache
}

//LayerConfig corresponds to an ESRI conf.xml document
type LayerConfig struct {
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
