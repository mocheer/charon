package agent

import (
	"github.com/gofiber/fiber/v2"
)

// ProxyImageCma 综合气象数据
// @see //data.cma.cn/dataGis/gis.html
// @see //image.data.cma.cn/vis/IMG_SURF_TEM_1H_AVG_nbg_mct/20190720/Guip_nmic_mct_cpas_T1_achn_nnn_1400_201907201400.png
func ProxyImageCma(c *fiber.Ctx) error {
	return proxyURL(c, "http://image.data.cma.cn/"+c.Params("*"))
}
