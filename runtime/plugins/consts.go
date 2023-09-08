package plugins

const NameActive = "active" // 激活 Key

// 存储插件

const (
	TypeStorage = "storage"

	NameStorageStandalone = "standalone" // Standalone
	NameStorageDB         = "database"   // DataBase
	NameStorageRedis      = "redis"      // Redis
	NameStorageMongodb    = "mongodb"    // MongoDB
)

// 连接器插件

const (
	TypeConnector = "connector"

	NameConnectorStandalone = "standalone" // Standalone
	NameConnectorRedis      = "redis"      // Redis
	NameConnectorRabbitMQ   = "rabbitmq"   // RabbitMQ
	NameConnectorRocketMQ   = "rocketmq"   // RocketMQ
	NameConnectorKafka      = "kafka"      // Kafka
	NameConnectorPulsar     = "pulsar"     // Pulsar
)

// 服务注册

const (
	TypeRegistry = "registry"

	NameRegistryNacos  = "nacos"     // Nocas
	NameRegistryZK     = "zookeeper" // Zookeeper
	NameRegistryConsul = "consul"    // Consul
	NameRegistryEtcd   = "etcd"      // Etcd
)

// 监控指标插件

const (
	TypeMetrics = "metrics"

	NameMetricsPrometheus = "prometheus" // Prometheus
)
