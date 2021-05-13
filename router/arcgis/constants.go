package arcgis

import (
	"path"

	"github.com/mocheer/charon/global"
)

// BaseDirectory arcgis图标的根目录
var BaseDirectory = path.Join(global.AssetsDataDir, "dmap/arcgis")

// BundlxMaxidx 紧凑型切片v2
const BundlxMaxidx = 128

// CompactCacheHeaderLength 紧凑型切片v2头信息长度
const CompactCacheHeaderLength = 64

// EsriMapCacheStorageModeCompact 紧凑型切片
const EsriMapCacheStorageModeCompact = "esriMapCacheStorageModeCompact"

// EsriMapCacheStorageModeCompactV2 紧凑型切片V2版本
const EsriMapCacheStorageModeCompactV2 = "esriMapCacheStorageModeCompactV2"
