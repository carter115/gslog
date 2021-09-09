package gslog

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// gin日志插件
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		GinWithAccessInfo(Logger, c)
		c.Next()
		GinWithAccessInfo(Logger, c)
	}
}

// 打印gin访问日志
func GinWithAccessInfo(l *MyLogger, c *gin.Context) {
	fields := logrus.Fields{
		traceIdKey:     getTraceIdFromContext(c),
		"type":         "gin-accesslog",
		"uri":          c.Request.RequestURI,
		"method":       c.Request.Method,
		"host":         c.Request.Host,
		"user_agent":   c.Request.UserAgent(),
		"client_ip":    c.ClientIP(),
		"content_type": c.ContentType(),
	}

	l.WithFields(fields).Info()
}

func getTraceIdFromContext(c *gin.Context) (traceId string) {
	if v := c.Value(traceIdKey); v != nil {
		traceId = v.(string)
		return
	}

	if traceId = c.Query(traceIdKey); traceId != "" {
		return
	}
	if traceId = c.PostForm(traceIdKey); traceId != "" {
		return
	}
	if traceId = c.Request.Header.Get(traceIdKey); traceId != "" {
		return
	}

	traceId = getUUID()
	c.Set(traceIdKey, traceId)
	return
}
