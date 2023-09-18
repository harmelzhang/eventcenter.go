package zookeeper

import (
	"errors"
	"eventcenter-go/runtime/plugins"
	"eventcenter-go/runtime/registry"
	"fmt"
	"github.com/go-zookeeper/zk"
	"github.com/gogf/gf/v2/container/gvar"
	"strings"
	"time"
)

var (
	defaultConnectTimeout = 10 * time.Second
	defaultPath           = "/services"
)

type plugin struct{}

func init() {
	plugins.Register(plugins.NameRegistryZK, &plugin{})
}

// Type 插件类型
func (p *plugin) Type() string {
	return plugins.TypeRegistry
}

// Init 初始化
func (p *plugin) Init(config map[string]*gvar.Var) (err error) {
	addresses := make([]string, 0)
	addressConfig := config["address"].String()
	if addressConfig == "" {
		err = errors.New("registry plugin [zookeeper] config must set [address]")
		return
	}
	for _, address := range strings.Split(addressConfig, ",") {
		addresses = append(addresses, address)
	}

	timeout := config["timeout"].Int()
	if timeout != 0 {
		defaultConnectTimeout = time.Duration(timeout) * time.Second
	}

	conn, _, err := zk.Connect(addresses, defaultConnectTimeout)
	if err != nil {
		err = errors.New(fmt.Sprintf("conn zookeeper err: %v", err))
		return
	}

	path := config["path"].String()
	if path == "" {
		path = defaultPath
	}

	registryService.conn = conn
	registryService.path = path

	return
}

// Service 注册服务
func (p *plugin) Service() registry.Service {
	return registryService
}
