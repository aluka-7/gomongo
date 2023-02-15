# gomongo引擎

## 配置信息
mongo权限地址:`/system/base/mongo/privileges`
```json
{
  "1000":[
    "1800", "2000", "2010"
  ]
}
```
数据源配置地址:`/system/base/mongo/{systemId}`
```json
{
  "uri":"mongodb://localhost:27017"
}
```

通用数据源配置地址：`/system/base/mongo/common`
```json
{
  "timeOut":2
}
```

```go
// 单个实例配置
type Config struct {
    Uri             string `json:"uri"`      // 连接uri
    Database        string `json:"database"` // 数据库
    MaxPoolSize     uint64 `json:"maxPoolSize"`
    MinPoolSize     uint64 `json:"minPoolSize"`     //
    MaxConnecting   uint64 `json:"maxConnecting"`   // 连接池最大连接数
    MaxConnIdleTime uint64 `json:"maxConnIdleTime"` // 连接最大空闲时间
    TimeOut         int64  `json:"timeOut"`
}
```