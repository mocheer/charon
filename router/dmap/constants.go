package dmap

import (
	"path"

	"github.com/mocheer/charon/global"
)

// 动态图层瓦片格式化字符串
var ImageTilePathFormat = path.Join(global.AssetsDataDir, "dmap/image/%v/%v/%v/%v.png")
