# charon

卡戎，冥王星的卫星


```json
{
  "Name": "charon", // 服务名称，将注册到注册表中支持自启动
  "Version": "1.0.0",
  "DisplayName":"应用服务",
  "Mode":"dev",// 开发模式，分成dev、production
  "Port": ":9912",
  "Static":{ // 静态资源服务器
    "/":{
      "mode":"history",
      "dir":"./public",
      "baseName":"v"
    },
    "/web":{
      "mode":"history",
      "dir":"./web",
      "baseName":"web"
    }
  },
  "DbDSN":"xxx" // 数据库连接串，支持 pg、mysql 等
}
```