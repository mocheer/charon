# charon

服务端

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