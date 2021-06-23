package image

// ImageArgs
type ImageArgs struct {
	W     int    // width  宽度
	H     int    // height 高度
	F     string // format 格式 img.PNG | img.JPEG | img.GIF
	Cache bool   // cache 是否缓存
}
