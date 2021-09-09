package main

import (
	"context"
	"github.com/carter115/gslog"
	"github.com/gin-gonic/gin"
)

func main() {
	// test log
	var ctx = context.Background()
	gslog.Info(ctx, "first log...")
	// output:
	//{"app_id":"IM-default","level":"info","msg":"first log...","time":"2021-09-09T14:27:33+08:00","trace_id":"71fdcc5fc94f416e98f58afc1ff6393c","type":"applog"}

	// 1. 使用默认的gslog.init() 生成的Logger
	// 2. 自定义日志配置
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

	r := gin.New()
	r.Use(gslog.GinLogger(), gin.Recovery())

	r.GET("/", func(c *gin.Context) {
		gslog.Info(c, "home page")
	})

	r.GET("/error", func(c *gin.Context) {
		gslog.Error(c, "error msg")
	})

	// output:
	//{"app_id":"liveshow-account","client_ip":"::1","content_type":"","host":"localhost:8080","level":"info","method":"GET","msg":"","time":"2021-09-09T14:36:37+08:00","trace_id":"1e6f6b64c1af4f079242fc4ed39bc98b","type":"gin-accesslog","uri":"/","user_agent":"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.66 Safari/537.36"}
	//{"app_id":"liveshow-account","level":"info","msg":"home page","time":"2021-09-09T14:36:37+08:00","trace_id":"1e6f6b64c1af4f079242fc4ed39bc98b","type":"applog"}
	//{"app_id":"liveshow-account","client_ip":"::1","content_type":"","host":"localhost:8080","level":"info","method":"GET","msg":"","time":"2021-09-09T14:36:37+08:00","trace_id":"1e6f6b64c1af4f079242fc4ed39bc98b","type":"gin-accesslog","uri":"/","user_agent":"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.66 Safari/537.36"}

	r.Run(":8080")
}
