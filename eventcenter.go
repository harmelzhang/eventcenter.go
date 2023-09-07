package main

import (
	"eventcenter-go/runtime"
	"log"
)

func main() {
	if err := runtime.LoadPlugins(); err != nil {
		log.Fatalf("load plugins err: %v", err)
	}
	if err := runtime.Start(); err != nil {
		log.Fatalf("start server err: %v", err)
	}
}
