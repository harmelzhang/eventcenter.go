package server

// CoreServer 核心服务
type CoreServer interface {
	// Start 启动服务
	Start()

	// Stop 停止服务
	Stop() error
}
