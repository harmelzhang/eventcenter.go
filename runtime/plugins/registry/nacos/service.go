package nacos

import (
	"errors"
	"eventcenter-go/runtime/registry"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"strconv"
	"strings"
)

type service struct {
	client naming_client.INamingClient
	weight float64
}

var registryService = new(service)

// Register 服务注册
func (s *service) Register(serviceName, address, protocol string) (err error) {
	ipAndPort := strings.Split(address, ":")
	if len(ipAndPort) != 2 {
		err = errors.New(fmt.Sprintf("address format err, must IP:PORT"))
		return
	}
	port, err := strconv.ParseUint(ipAndPort[1], 10, 64)
	if err != nil {
		return
	}

	isSuccess, err := s.client.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          ipAndPort[0],
		Port:        port,
		ServiceName: serviceName,
		Healthy:     true,
		Enable:      true,
		Weight:      s.weight,
		Metadata:    map[string]string{"protocol": protocol},
		Ephemeral:   true,
	})
	if err != nil {
		return
	}
	if !isSuccess {
		err = errors.New(fmt.Sprintf("register service instance failure"))
		return
	}

	return
}

// FindService 查找服务
func (s *service) FindService(serviceName string) (ins *registry.Instance, err error) {
	instance, err := s.client.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: serviceName,
	})
	if instance != nil {
		ins = new(registry.Instance)
		ins.Address = instance.Ip
		ins.Port = int(instance.Port)
	}
	return
}
