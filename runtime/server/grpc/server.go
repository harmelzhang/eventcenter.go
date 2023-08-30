package grpc

import (
	"eventcenter-go/runtime/server"
	"log"
)

type GRPCServer struct{}

func New() server.CoreServer {
	return &GRPCServer{}
}

func (s GRPCServer) Start() {
	log.Println("gRPC Server not implemented")
}

func (s GRPCServer) Stop() error {
	return nil
}
