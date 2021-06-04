package pano

import (
	"path/filepath"

	"github.com/mocheer/charon/cts"
)

// 全景图配置文件-配置项
// @see https://krpano.com/tools/kmakemultires/config/?version=119#askforxmloverwrite
// askforxmloverwrite=false

/**
* krpano 工具根目录
 */
var KRPANO_ROOT = filepath.Join(cts.Assets, `scripts/krpano-1.19-pr16`)

/**
 * 注册码，目前已知支持 krpano-1.19-xx
 */
var REGISTER_CODE_PATH = filepath.Join(KRPANO_ROOT, `key.txt`)

/**
* 可执行文件的名称
 */
var KRPANO_EXE_NAME = filepath.Join(KRPANO_ROOT, `krpanotools64.exe`)

/**
* 配置文件的路径(单场景)
 */
var KRPANO_CONFIG_PATH = filepath.Join(KRPANO_ROOT, `templates/vtour-normal.config`)

/**
* 配置文件的路径(多场景)
 */
var KRPANO_MULTIRES_CONFIG_PATH = filepath.Join(KRPANO_ROOT, `templates/vtour-multires.config`)

/**
* 注册脚本，todo 修改
* krpanolicense
* 手动注册后如果在程序部署启动后仍有水印，首先用代码调用工具的命令执行一次注册就可以
* 造成的原因是因为Krpano的注册码读取的是当前的User Context下的注册信息（换句话说，KrPano的注册信息是存储在所执行用户的相应位置的）;
 */
// var KRPANO_REGISTER_SCRIPT = fmt.Sprintf(`%s register %s`, KRPANO_EXE_NAME, REGISTER_CODE)
