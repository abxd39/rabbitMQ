###### 配置文件config.json
`{
  "Listen":":5000",
  "runtimePath":"runtime",
  "maxWork":10,
  "maxQueueSize":200,
  "db": {
    "host":"10.20.20.82",
    "name":"typhoon",
    "user":"root",
    "password":"",
    "maxOpenConns":100,
    "maxIdleConns":10
  },
  "log": {
    "traceLevel": 1,
    "logFile":"th-platform-gateway.log"
  }
}`


###### 启动方式
nohup th-platform-gateway &


