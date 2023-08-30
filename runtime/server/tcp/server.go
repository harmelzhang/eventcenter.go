package tcp

import (
	"eventcenter-go/runtime/server"
	"log"
)

type TCPServer struct{}

func New() server.CoreServer {
	return &TCPServer{}
}

func (s TCPServer) Start() {
	log.Println("TCP Server not implemented")
}

func (s TCPServer) Stop() error {
	return nil
}
