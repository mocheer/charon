package global

import "path"

const Assets = "assets"

// AssetsConfigDir 配置文件目录
var AssetsConfigDir = path.Join(Assets, "config")

// AssetsDataDir 数据文件目录
var AssetsDataDir = path.Join(Assets, "data")

// AssetsLogDir 日志文件目录
var AssetsLogDir = path.Join(Assets, "log")

// AppConfigPath 应用配置文件路径
var AppConfigPath = path.Join(AssetsConfigDir, "app.json")

// rsa pem
var RSA_Dir = path.Join(Assets, "rsa")
var RSA_PrivatePemPath = path.Join(RSA_Dir, "private.pem")
var RSA_PublicPemPath = path.Join(RSA_Dir, "public.pem")

// DevMode 开发模式
const DevMode = "dev"

// ProdMode 生产模式
const ProdMode = "production"
