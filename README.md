# charon

卡戎，冥王星的卫星。主要用于提供地图方面的数据服务。

- 提供类似GraphQL风格或者说apijson风格的服务接口
- 提供地图服务接口
- 提供pipal框架的服务接口
- 提供网络爬虫、代理等服务

### app.json
```json
{
  "Name": "charon",
  "Version": "1.0.0",
  "DisplayName":"应用服务",
  "Mode":"dev",
  "Port": ":9912",
  "Static":{
    "/":{
      "mode":"history",
      "dir":"./public"
    },
    "/web":{
      "mode":"history",
      "dir":"./web"
    }
  },
  "DbDSN":"xxx"
}
```