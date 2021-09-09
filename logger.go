package gslog

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"runtime/debug"
	"strings"
)

var (
	Logger            *MyLogger
	projectNameKey    = "IM"
	appDefaultNameKey = "default"

	appKey     = "app_id"
	traceIdKey = "trace_id"
	stackKey   = "stack"
)

// 日志组件配置
type Config struct {
	ProjectName string   // im
	AppName     string   // account,...
	FileName    string   // output.log
	Level       string   // debug,info,...
	EsServer    []string // http://192.168.100.20:9200
	StashServer string   // 192.168.100.20:4560
	Hooks       []string // es,stash
	Outputs     []string // stdout,file
}

type MyLogger struct {
	log *logrus.Logger
	*logrus.Entry
}

// 初始化默认Logger
func init() {
	conf := Config{ProjectName: projectNameKey, AppName: appDefaultNameKey, Level: "info", Outputs: []string{"stdout"}}
	if err := InitLogger(conf); err != nil {
		panic(err)
	}
}

func InitLogger(conf Config) (err error) {
	logr := logrus.New()

	// 添加app默认字段
	appId := fmt.Sprintf("%s-%s", conf.ProjectName, conf.AppName)
	ent := logr.WithFields(logrus.Fields{appKey: appId, "type": "applog"}) // *logrus.Entry

	// 解析level
	lv, err := logrus.ParseLevel(conf.Level)
	if err != nil {
		return err
	}
	logr.SetLevel(lv)
	logr.SetFormatter(&logrus.JSONFormatter{})
	logr.SetNoLock() // 关闭logrus互斥锁

	// 多种日志输出
	if len(conf.Outputs) > 0 {
		logr.SetOutput(multiWriter(conf.Outputs, conf.FileName))
	}

	mylogger := &MyLogger{logr, ent}
	mylogger.AddHooks(conf, appId) // 添加多个Hook
	Logger = mylogger

	return nil
}

// 生成ID
func getUUID() string {
	uuidv4 := uuid.New().String()
	newuuid := strings.Replace(uuidv4, "-", "", -1)
	return newuuid
}

// 打印traceIdKey
func (l *MyLogger) WithTraceId(c context.Context) *logrus.Entry {
	if c == nil {
		c = context.Background()
	}
	fields := logrus.Fields{}
	if id := c.Value(traceIdKey); id != nil {
		fields[traceIdKey] = id
	} else {
		fields[traceIdKey] = getUUID()
	}
	return l.WithFields(fields)
}

// 出错时打印堆栈信息
func (l *MyLogger) WithStack(c context.Context) *logrus.Entry {
	return l.WithTraceId(c).WithFields(logrus.Fields{stackKey: string(debug.Stack())})
}

// 添加Hook
func (l *MyLogger) AddHooks(conf Config, indexName string) {
	for _, t := range conf.Hooks {
		switch t {
		case "es":
			if hook := NewEsHook(conf.EsServer, indexName, l.log.GetLevel()); hook != nil {
				l.Logger.AddHook(hook)
			}
		case "stash":
			if hook := NewLogstashHook(conf.StashServer); hook != nil {
				l.Logger.AddHook(hook)
			}
		}
	}
}

// 多种写入日志方法
func multiWriter(outputs []string, fn string) io.Writer {
	var (
		writers = make([]io.Writer, 0)
	)
	for _, wr := range outputs {

		switch strings.ToLower(wr) {
		case "stdout":
			writers = append(writers, os.Stdout)

		case "file":
			f, err := os.OpenFile(fn, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
			if err != nil {
				fmt.Println("open file:", err)
			}
			writers = append(writers, f)
		}
	}
	return io.MultiWriter(writers...)
}

// 封装输出日志的方法
func Debug(c context.Context, v ...interface{}) {
	Logger.WithTraceId(c).Debug(v...)
}

func Debugf(c context.Context, format string, v ...interface{}) {
	Logger.WithTraceId(c).Debugf(format, v...)
}

func Info(c context.Context, v ...interface{}) {
	Logger.WithTraceId(c).Info(v...)
}

func Infof(c context.Context, format string, v ...interface{}) {
	Logger.WithTraceId(c).Infof(format, v...)
}

func Warning(c context.Context, v ...interface{}) {
	Logger.WithTraceId(c).Warning(v...)
}

func Warningf(c context.Context, format string, v ...interface{}) {
	Logger.WithTraceId(c).Warningf(format, v...)
}

func Error(c context.Context, v ...interface{}) {
	Logger.WithStack(c).Error(v...)
}

func Errorf(c context.Context, format string, v ...interface{}) {
	Logger.WithStack(c).Errorf(format, v...)
}
