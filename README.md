# gslog

简单的日志组件：
- 对logrus进行包装
- 支持多种输出(标准输出，文件输出等)
- 支持es和stash的Hook
- json格式日志
- 打印跟踪ID
- 错误日志打印堆栈信息

## 初始化

- 默认初始化，可以直接使用

```go
var ctx = context.Background()
gslog.Info(ctx, "first log...")
```

- 手动初始化

```go
	logConfig := gslog.Config{
		ProjectName: "liveshow",
		AppName:     "account",
		Level:       "info",
		EsServer:    []string{"http://192.168.100.20:9200"},
		StashServer: "192.168.100.20:4560",
		Hooks:       []string{"stash"},
		Outputs:     []string{"stdout"},
	}
	if err := gslog.InitLogger(logConfig); err != nil {
		panic(err)
	}
```

## gin logger插件

```go
	r := gin.New()
	r.Use(gslog.GinLogger(), gin.Recovery())

	r.GET("/", func(c *gin.Context) {
		gslog.Info(c, "home page")
	})


	r.Run(":8080")
```

```shell
	// output:
	//{"app_id":"liveshow-account","client_ip":"::1","content_type":"","host":"localhost:8080","level":"info","method":"GET","msg":"","time":"2021-09-09T14:36:37+08:00","trace_id":"1e6f6b64c1af4f079242fc4ed39bc98b","type":"gin-accesslog","uri":"/","user_agent":"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.66 Safari/537.36"}
	//{"app_id":"liveshow-account","level":"info","msg":"home page","time":"2021-09-09T14:36:37+08:00","trace_id":"1e6f6b64c1af4f079242fc4ed39bc98b","type":"applog"}
	//{"app_id":"liveshow-account","client_ip":"::1","content_type":"","host":"localhost:8080","level":"info","method":"GET","msg":"","time":"2021-09-09T14:36:37+08:00","trace_id":"1e6f6b64c1af4f079242fc4ed39bc98b","type":"gin-accesslog","uri":"/","user_agent":"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.66 Safari/537.36"}
```

## 使用Hook输出

`gslog.Info(c, "home page")`

- es

```json
{
	"_index": "liveshow-account-20210909",
	"_type": "log",
	"_id": "AXvJPc6y8uY57suJ3sM8",
	"_version": 1,
	"_score": 1,
	"_source": {
		"Host": "192.168.100.66",
		"@timestamp": "2021-09-09T06:27:55.8927138Z",
		"Message": "",
		"Data": {
			"app_id": "liveshow-account",
			"client_ip": "::1",
			"content_type": "",
			"host": "localhost:8080",
			"method": "GET",
			"trace_id": "4f08a766d10e438abad75fa4bf7b4e5d",
			"type": "gin-accesslog",
			"uri": "/",
			"user_agent": "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.66 Safari/537.36"
		},
		"Level": "INFO"
	}
}

{
	"_index": "liveshow-account-20210909",
	"_type": "log",
	"_id": "AXvJPc618uY57suJ3sM9",
	"_version": 1,
	"_score": 1,
	"_source": {
		"Host": "192.168.100.66",
		"@timestamp": "2021-09-09T06:27:55.8927138Z",
		"Message": "",
		"Data": {
			"app_id": "liveshow-account",
			"client_ip": "::1",
			"content_type": "",
			"host": "localhost:8080",
			"method": "GET",
			"trace_id": "4f08a766d10e438abad75fa4bf7b4e5d",
			"type": "gin-accesslog",
			"uri": "/",
			"user_agent": "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.66 Safari/537.36"
		},
		"Level": "INFO"
	}
}

{
	"_index": "liveshow-account-20210909",
	"_type": "log",
	"_id": "AXvJPc618uY57suJ3sM-",
	"_version": 1,
	"_score": 1,
	"_source": {
		"Host": "192.168.100.66",
		"@timestamp": "2021-09-09T06:27:55.8927138Z",
		"Message": "home page",
		"Data": {
			"app_id": "liveshow-account",
			"trace_id": "4f08a766d10e438abad75fa4bf7b4e5d",
			"type": "applog"
		},
		"Level": "INFO"
	}
}
```

- stash

```json
{
	"_index": "test-2021.09.09",
	"_type": "gin-accesslog",
	"_id": "AXvJRcLN8uY57suJ3sNB",
	"_version": 1,
	"_score": null,
	"_source": {
		"trace_id": "1e6f6b64c1af4f079242fc4ed39bc98b",
		"method": "GET",
		"level": "info",
		"message": "",
		"type": "gin-accesslog",
		"uri": "/",
		"@timestamp": "2021-09-09T06:36:37.000Z",
		"content_type": "",
		"port": 64645,
		"@version": "1",
		"host": "localhost:8080",
		"client_ip": "::1",
		"app_id": "liveshow-account",
		"user_agent": "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.66 Safari/537.36"
	},
	"sort": [1631169397000]
}

{
	"_index": "test-2021.09.09",
	"_type": "gin-accesslog",
	"_id": "AXvJRcLN8uY57suJ3sND",
	"_version": 1,
	"_score": null,
	"_source": {
		"trace_id": "1e6f6b64c1af4f079242fc4ed39bc98b",
		"method": "GET",
		"level": "info",
		"message": "",
		"type": "gin-accesslog",
		"uri": "/",
		"@timestamp": "2021-09-09T06:36:37.000Z",
		"content_type": "",
		"port": 64645,
		"@version": "1",
		"host": "localhost:8080",
		"client_ip": "::1",
		"app_id": "liveshow-account",
		"user_agent": "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.66 Safari/537.36"
	},
	"sort": [1631169397000]
}

{
	"_index": "test-2021.09.09",
	"_type": "applog",
	"_id": "AXvJRcLN8uY57suJ3sNC",
	"_version": 1,
	"_score": null,
	"_source": {
		"trace_id": "1e6f6b64c1af4f079242fc4ed39bc98b",
		"@timestamp": "2021-09-09T06:36:37.000Z",
		"level": "info",
		"port": 64645,
		"@version": "1",
		"host": "192.168.100.66",
		"message": "home page",
		"type": "applog",
		"app_id": "liveshow-account"
	},
	"sort": [1631169397000]
}
```

## 打印错误日志

`gslog.Error(c, "error msg")`

```json
{
	"_index": "test-2021.09.09",
	"_type": "applog",
	"_id": "AXvJWPMG8uY57suJ3sNI",
	"_version": 1,
	"_score": null,
	"_source": {
		"stack": "goroutine 34 [running]: runtime/debug.Stack(0xc0137d6639, 0x985598, 0xc000406300) C:/go1.16.5/src/runtime/debug/stack.go:24 +0xa5 github.com/carter115/gslog.(*MyLogger).WithStack(0xc0003946c0, 0x985598, 0xc000406300, 0x8e26e6) C:/Users/Administrator/Desktop/gslog/logger.go:100 +0xd1 github.com/carter115/gslog.Error(0x985598, 0xc000406300, 0xc000431960, 0x1, 0x1) C:/Users/Administrator/Desktop/gslog/logger.go:167 +0x48 main.main.func2(0xc000406300) C:/Users/Administrator/Desktop/gslog/example/example.go:38 +0x7a github.com/gin-gonic/gin.(*Context).Next(...) C:/Users/Administrator/gohome/pkg/mod/github.com/gin-gonic/gin@v1.7.4/context.go:165 github.com/gin-gonic/gin.CustomRecoveryWithWriter.func1(0xc000406300) C:/Users/Administrator/gohome/pkg/mod/github.com/gin-gonic/gin@v1.7.4/recovery.go:99 +0x82 github.com/gin-gonic/gin.(*Context).Next(...) C:/Users/Administrator/gohome/pkg/mod/github.com/gin-gonic/gin@v1.7.4/context.go:165 github.com/carter115/gslog.GinLogger.func1(0xc000406300) C:/Users/Administrator/Desktop/gslog/gin_logger.go:12 +0x63 github.com/gin-gonic/gin.(*Context).Next(...) C:/Users/Administrator/gohome/pkg/mod/github.com/gin-gonic/gin@v1.7.4/context.go:165 github.com/gin-gonic/gin.(*Engine).handleHTTPRequest(0xc0003d6680, 0xc000406300) C:/Users/Administrator/gohome/pkg/mod/github.com/gin-gonic/gin@v1.7.4/gin.go:489 +0x2b0 github.com/gin-gonic/gin.(*Engine).ServeHTTP(0xc0003d6680, 0x984750, 0xc00041c1c0, 0xc000314100) C:/Users/Administrator/gohome/pkg/mod/github.com/gin-gonic/gin@v1.7.4/gin.go:445 +0x165 net/http.serverHandler.ServeHTTP(0xc00041c0e0, 0x984750, 0xc00041c1c0, 0xc000314100) C:/go1.16.5/src/net/http/server.go:2887 +0xaa net/http.(*conn).serve(0xc00031a000, 0x985480, 0xc00030c080) C:/go1.16.5/src/net/http/server.go:1952 +0x8cd created by net/http.(*Server).Serve C:/go1.16.5/src/net/http/server.go:3013 +0x3b8 ",
		"trace_id": "8d15850914f14848b2770de444afe94e",
		"@timestamp": "2021-09-09T06:57:34.000Z",
		"level": "error",
		"port": 64907,
		"@version": "1",
		"host": "192.168.100.66",
		"message": "error msg",
		"type": "applog",
		"app_id": "liveshow-account"
	},
	"sort": [1631170654000]
}
```