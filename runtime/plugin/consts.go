package plugin

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
	NameConnectorRabbitMQ   = "rabbitmq"   // RabbitMQ
	NameConnectorKafka      = "kafka"      // Kafka
	NameConnectorPulsar     = "pulsar"     // Pulsar
)

// 监控指标插件

const (
	TypeMetrics = "metrics"

	NameMetricsPrometheus = "prometheus" // Prometheus
)
