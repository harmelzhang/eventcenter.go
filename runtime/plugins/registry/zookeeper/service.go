package zookeeper

import (
	"encoding/json"
	"errors"
	"eventcenter-go/runtime/registry"
	"fmt"
	"github.com/go-zookeeper/zk"
	"github.com/google/uuid"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"sync"
)

type service struct {
	conn      *zk.Conn
	path      string
	mutex     sync.RWMutex
	instances map[string][]registry.Instance // 服务清单
}

var registryService = new(service)

// Spring Cloud 规范最小子集
type serverInfo struct {
	Name    string `json:"name"`
	Id      string `json:"id"`
	Address string `json:"address"`
	Port    int    `json:"port"`
}

// Register 服务注册
func (s *service) Register(serviceName, address, protocol string) (err error) {
	registryService.instances = make(map[string][]registry.Instance)

	ipAndPort := strings.Split(address, ":")
	if len(ipAndPort) != 2 {
		err = errors.New(fmt.Sprintf("address format err, must IP:PORT"))
		return
	}
	port, err := strconv.Atoi(ipAndPort[1])
	if err != nil {
		return
	}

	// 顶级节点
	path := s.path
	err = s.createNode(path)
	if err != nil {
		return
	}
	//err = s.queryServers()
	//if err != nil {
	//	return
	//}
	go s.watch(path, false)

	// 服务节点
	path = s.path + "/" + serviceName
	err = s.createNode(path)
	if err != nil {
		return
	}
	//if _, isOk := s.instances[serviceName]; !isOk {
	//	s.instances[serviceName] = make([]registry.Instance, 0)
	//}
	//go s.watch(path, true)

	// 实例节点
	info := serverInfo{
		Name:    serviceName,
		Id:      uuid.NewString(),
		Address: ipAndPort[0],
		Port:    port,
	}
	data, err := json.Marshal(info)
	if err != nil {
		return
	}
	_, err = s.conn.Create(s.path+"/"+serviceName+"/"+info.Id, data, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err != nil {
		return err
	}
	//instance := registry.Instance{Address: info.Address, Port: info.Port}
	//s.instances[serviceName] = append(s.instances[serviceName], instance)

	err = s.queryServers()
	if err != nil {
		return
	}

	return
}

// FindService 查找服务
func (s *service) FindService(serviceName string) (ins *registry.Instance, err error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if _, isOK := s.instances[serviceName]; !isOK {
		return
	}
	cnt := len(s.instances[serviceName])
	return &s.instances[serviceName][rand.Intn(cnt)], nil
}

// 监听节点变动
func (s *service) watch(path string, isInstance bool) {
	for {
		if s.conn.State() == zk.StateConnected || s.conn.State() == zk.StateHasSession {
			isExists, _, err := s.conn.Exists(path)
			if err != nil {
				log.Printf("query node exists err: %v", err)
				continue
			}
			if !isExists {
				continue
			}

			_, _, ch, err := s.conn.ChildrenW(path)
			if err != nil {
				log.Printf("new watch err: %v", err)
				continue
			}
			<-ch

			if isInstance {
				paths := strings.Split(path, "/")
				serverName := paths[len(paths)-1]
				err = s.queryInstances(serverName)
				if err != nil {
					log.Printf("query instances err [%s]: %v", path, err)
					s.instances[serverName] = make([]registry.Instance, 0)
				}
			} else {
				err = s.queryServers()
				if err != nil {
					log.Printf("query servers err [%s]: %v", path, err)
					s.instances = make(map[string][]registry.Instance, 0)
				}
			}
			log.Println("server list:", s.instances)
		}
	}
}

// 创建节点
func (s *service) createNode(path string) error {
	isExists, _, err := s.conn.Exists(path)
	if err != nil {
		return err
	}
	if !isExists {
		_, err = s.conn.Create(path, []byte{}, 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			return err
		}
	}
	return nil
}

// 查询子节点
func (s *service) findChildNode(path string) (nodePaths []string, err error) {
	nodePaths, _, err = s.conn.Children(path)
	return
}

// 查询服务
func (s *service) queryServers() (err error) {
	servicePaths, err := s.findChildNode(s.path)
	if err != nil {
		return
	}
	for _, serviceName := range servicePaths {
		if _, isOK := s.instances[serviceName]; !isOK {
			go s.watch(s.path+"/"+serviceName, true)
		}
		err = s.queryInstances(serviceName)
		if err != nil {
			return err
		}
	}
	return
}

// 查询实例
func (s *service) queryInstances(serverName string) (err error) {
	instancePaths, err := s.findChildNode(s.path + "/" + serverName)
	if err != nil {
		return
	}
	instances := make([]registry.Instance, 0)
	for _, instancePath := range instancePaths {
		data, _, err := s.conn.Get(s.path + "/" + serverName + "/" + instancePath)
		if err != nil {
			return err
		}
		instance := registry.Instance{}
		err = json.Unmarshal(data, &instance)
		if err != nil {
			return err
		}
		instances = append(instances, instance)
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.instances[serverName] = instances
	return nil
}
