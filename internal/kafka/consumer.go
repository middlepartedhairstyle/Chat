package Kafka

import (
	"github.com/middlepartedhairstyle/HiWe/internal/utils"
	"github.com/segmentio/kafka-go"
)

type SetConsumerGroupConfig func(config *kafka.ReaderConfig)

// SetConsumerTopic 设置消费者的topic
func SetConsumerTopic(str string) SetConsumerGroupConfig {
	return func(config *kafka.ReaderConfig) {
		config.Topic = str
	}
}

// SetConsumerBrokers 设置消费者broker
func SetConsumerBrokers(broker []string) SetConsumerGroupConfig {
	return func(config *kafka.ReaderConfig) {
		config.Brokers = broker
	}
}

// SetConsumerGroupID 设置消费者groupID
func SetConsumerGroupID(str string) SetConsumerGroupConfig {
	return func(config *kafka.ReaderConfig) {
		config.GroupID = str
	}
}

// SetConsumerMaxBytes 设置消费者maxBytes
func SetConsumerMaxBytes(max int) SetConsumerGroupConfig {
	return func(config *kafka.ReaderConfig) {
		config.MaxBytes = max
	}
}

// NewConsumer 新建消费者
func NewConsumer(cfg ...SetConsumerGroupConfig) *kafka.Reader {
	// make a new reader that consumes from topic-A
	var config = kafka.ReaderConfig{
		Brokers:  []string{utils.Cfg.Kafka.Addr},
		GroupID:  "0",
		Topic:    "tp1",
		MaxBytes: 10e6, // 10MB
	}
	for _, c := range cfg {
		c(&config)
	}

	r := kafka.NewReader(config)
	return r
}

//// CreateTopicWithRetention 创建话题并设置过期时间,args[0]为过期时间,args[1]为addr(kafka服务地址)
//func (producer *Producer) CreateTopicWithRetention(topic string, args ...string) bool {
//	retention := "86400000"
//	addr := "23.95.15.178:9092"
//	if len(args) > 0 {
//		retention = args[0]
//	}
//	if len(args) > 1 {
//		addr = args[1]
//	}
//
//	conn, err := kafka.Dial("tcp", addr)
//	if err != nil {
//		panic(err.Error())
//	}
//	defer conn.Close()
//
//	controller, err := conn.Controller()
//	if err != nil {
//		panic(err.Error())
//	}
//	var controllerConn *kafka.Conn
//	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
//	if err != nil {
//		panic(err.Error())
//	}
//	defer controllerConn.Close()
//
//	configEntries := make([]kafka.ConfigEntry, 1)
//	configEntries[0] = kafka.ConfigEntry{
//		ConfigName:  "retention.ms",
//		ConfigValue: retention,
//	}
//
//	topicConfigs := []kafka.TopicConfig{
//		{
//			Topic:             topic,
//			NumPartitions:     1,
//			ReplicationFactor: 1,
//			ConfigEntries:     configEntries,
//		},
//	}
//
//	err = controllerConn.CreateTopics(topicConfigs...)
//	if err != nil {
//		return false
//	}
//	return true
//}
//
//// GetTopic 获取kafka，topic是否存在
//func (producer *Producer) GetTopic(topic string, args ...string) bool {
//	addr := "23.95.15.178:9092"
//	for _, arg := range args {
//		addr = arg
//	}
//	conn, err := kafka.Dial("tcp", addr)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer conn.Close()
//
//	// 获取所有 topic
//	topics, err := conn.ReadPartitions(topic)
//
//	if err != nil {
//		return false
//	}
//
//	// 打印所有 topic 的名称
//	for _, t := range topics {
//		fmt.Println("Topic:", t.Topic)
//	}
//	return true
//}
