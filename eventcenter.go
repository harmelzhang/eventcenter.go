package main

import (
	"eventcenter-go/runtime"
	"log"
)

func main() {
	if err := runtime.LoadPlugins(); err != nil {
		log.Fatalf("设置插件出错：%v", err)
	}
	if err := runtime.InitPliguins(); err != nil {
		log.Fatalf("初始化插件出错：%v", err)
	}
	if err := runtime.Start(); err != nil {
		log.Fatalf("启动服务出错：%v", err)
	}
}
