package tcp

import (
	"eventcenter-go/runtime/server"
	"log"
)

type TCPServer struct{}

func New() server.CoreServer {
	return &TCPServer{}
}

// Start 启动服务
func (s TCPServer) Start() {
	log.Println("TCP Server not implemented")
}

// Stop 停止服务
func (s TCPServer) Stop() error {
	return nil
}
