package gslog

import (
	"fmt"
	logrustash "github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"gopkg.in/sohlich/elogrus.v7"
	"net"
	"time"
)

func GetLocalIP() string {
	defaultIp := "127.0.0.1"
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return defaultIp
	}
	for _, addr := range addrs {
		ipAddr, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}
		if ipAddr.IP.IsLoopback() {
			continue
		}
		if !ipAddr.IP.IsGlobalUnicast() {
			continue
		}
		return ipAddr.IP.String()
	}
	return defaultIp
}

// 1. es hook
func NewEsHook(esUrls []string, name string, level logrus.Level) *elogrus.ElasticHook {
	host := GetLocalIP()
	client, err := elastic.NewClient(elastic.SetURL(esUrls...))
	if err != nil {
		fmt.Println("connect es server error:", err)
		return nil
	}

	// 根据日期，定义生成的ES索引名字
	indexFunc := func() string {
		return name + "-" + time.Now().Format("20060102")
	}

	hook, err := elogrus.NewAsyncElasticHookWithFunc(client, host, level, indexFunc)
	if err != nil {
		fmt.Println("add es hook error:", err)
		return nil
	}
	fmt.Printf("add es hook: Addr: %s, Name: %s, Level: %s\n", esUrls, name, level)
	return hook
}

// 2. logstatsh hook
func NewLogstashHook(stashUrl string) *logrustash.Hook {
	hook, err := logrustash.NewHook("tcp", stashUrl, "")
	if err != nil {
		return nil
	}
	fmt.Println("add stash hook:", stashUrl)
	return hook
}
