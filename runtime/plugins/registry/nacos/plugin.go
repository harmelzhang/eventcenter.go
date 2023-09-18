package nacos

import (
	"errors"
	"eventcenter-go/runtime/plugins"
	"eventcenter-go/runtime/registry"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"net"
	"strings"
	"time"
)

var (
	defaultConnectTimeout         = 5 * time.Second
	defaultWeight         float64 = 10
)

type plugin struct{}

func init() {
	plugins.Register(plugins.NameRegistryNacos, &plugin{})
}

// Type 插件类型
func (p *plugin) Type() string {
	return plugins.TypeRegistry
}

// Init 初始化
func (p *plugin) Init(config map[string]*gvar.Var) (err error) {
	param := vo.NacosClientParam{
		ClientConfig: &constant.ClientConfig{
			TimeoutMs: uint64(defaultConnectTimeout.Milliseconds()),
		},
	}

	addressConfig := config["address"].String()
	if addressConfig == "" {
		err = errors.New("registry plugin [nacos] config must set [address]")
		return
	}
	for _, address := range strings.Split(addressConfig, ",") {
		ip, port, err := net.SplitHostPort(address)
		if err != nil {
			return err
		}
		param.ServerConfigs = append(param.ServerConfigs, constant.ServerConfig{IpAddr: ip, Port: gconv.Uint64(port)})
	}

	registryService.client, err = clients.NewNamingClient(param)
	if err != nil {
		return err
	}

	serviceConfig := config["service"].MapStrVar()

	weight := serviceConfig["weight"].Float64()
	if weight == 0 {
		weight = defaultWeight
	}
	registryService.weight = weight

	return
}

// Service 注册服务
func (p *plugin) Service() registry.Service {
	return registryService
}
