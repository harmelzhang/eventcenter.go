package registry

// Instance 服务实例
type Instance struct {
	Address string // 地址
	Port    int    // 端口
}

type Service interface {
	// Register 服务注册
	Register(serviceName, address, protocol string) (err error)

	// FindService 查找服务
	FindService(serviceName string) (ins *Instance, err error)
}
