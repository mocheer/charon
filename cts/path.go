package cts

const Assets = "assets"

// ConfigDir 配置文件目录
const ConfigDir = Assets + "/config"

// DataDir 数据文件目录
const DataDir = Assets + "/data"

// LogDir 日志文件目录
const LogDir = Assets + "/log"

// AppConfigPath 应用配置文件路径
const AppConfigPath = ConfigDir + "/app.json"

// cert 用于https
const Cert = Assets + "/cert"
const Cert_PEM = Cert + "/cert.pem"
const Cert_KEY = Cert + "/cert.key"

// rsa pem
const RSA_Dir = Assets + "/rsa"
const RSA_PrivatePemPath = RSA_Dir + "/private.pem"
const RSA_PublicPemPath = RSA_Dir + "/public.pem"
