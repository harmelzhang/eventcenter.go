package grpc

import (
	"eventcenter-go/runtime/server"
	"log"
)

type GRPCServer struct{}

func New() server.CoreServer {
	return &GRPCServer{}
}

// Start 启动服务
func (s GRPCServer) Start() {
	log.Println("gRPC Server not implemented")
}

// Stop 停止服务
func (s GRPCServer) Stop() error {
	return nil
}
