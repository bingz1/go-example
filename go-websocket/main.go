package main

import (
	"flag"
	"fmt"
	"github.com/wantgo/go-example/go-websocket/gateway"
	"os"
	"time"
)

var (
	confFile string // 配置文件路径
)

func initArgs() {
	flag.StringVar(&confFile, "config", "./go-websocket/config.json", "where config.json is.")
	flag.Parse()
}

func main() {

	//初始化配置路径
	initArgs()
	//加载配置参数
	err := gateway.InitConfig(confFile)
	if err != nil {
		goto ERR
	}

	//开启websocket服务  用以接收客户端的socket请求
	gateway.InitWSServer()

	//开启http服务  用以接收http请求
	gateway.InitService()

	for {
		time.Sleep(1 * time.Second)
	}

	os.Exit(0)

ERR:
	fmt.Fprintln(os.Stderr, err)
	os.Exit(-1)
}
