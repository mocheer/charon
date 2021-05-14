package dmap

import (
	"path"

	"github.com/mocheer/charon/cts"
)

// 动态图层瓦片格式化字符串
var ImageTilePathFormat = path.Join(cts.DataDir, "dmap/image/%v/%v/%v/%v.png")
