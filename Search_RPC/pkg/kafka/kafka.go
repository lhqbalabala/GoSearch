package kafka

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Shopify/sarama"
)

const (
	kafkaTimeOut = time.Second * 5

	// kafka生产者发送信息方式
	Sync  = "Sync"
	Async = "Async"
)

var (
	kafkaAddressError = errors.New("kafka address is error")
)

// NewProducer 新建kafka生产者并选择同步方式
// 一般使用异步方式
func NewProducer(address, topic string, duration time.Duration, syncOrAsync string) AbsKafkaProducer {
	var producer AbsKafkaProducer
	switch syncOrAsync {
	case Sync:
		producer = &SyncKafkaProducer{}
	case Async:
		producer = &AsyncKafkaProducer{}
	default:
		fmt.Println("kafka mode error please choose sync or async")
		return nil
	}
	err := producer.NewKafkaProducer(address, topic, duration)
	if err != nil {
		fmt.Println("new kafka producer error")
		return nil
	}
	return producer
}

// AbsKafkaProducer 生产者接口
type AbsKafkaProducer interface {
	NewKafkaProducer(address, topic string, duration time.Duration) error
	Send(value []byte)
}

// KafkaConfig kafka生产者
type KafkaConfig struct {
	addressList []string       //地址列表
	topic       string         //kafka topic
	config      *sarama.Config //kafka配置信息
}

// NewProducerByMessage 创建kafka基本生产者
func NewProducerByMessage(address, topic string, duration time.Duration) *KafkaConfig {

	//根据字符串解析地址列表
	addressList := strings.Split(address, ",")
	if len(addressList) < 1 || addressList[0] == "" {
		fmt.Println("kafka addr error")
		return nil
	}
	//配置producer参数
	sendConfig := sarama.NewConfig()
	sendConfig.Producer.Return.Successes = true
	sendConfig.Producer.Timeout = kafkaTimeOut
	if duration != 0 {
		sendConfig.Producer.Timeout = duration
	}
	return &KafkaConfig{
		addressList: addressList,
		topic:       topic,
		config:      sendConfig,
	}
}

// SyncKafkaProducer 同步kafka生产者
type SyncKafkaProducer struct {
	KafkaConfig *KafkaConfig
}

func (k *SyncKafkaProducer) NewKafkaProducer(address, topic string, duration time.Duration) error {
	if len(address) == 0 {
		return kafkaAddressError
	}
	k.KafkaConfig = NewProducerByMessage(address, topic, duration)
	return nil
}

func (k *SyncKafkaProducer) Send(value []byte) {
	if k == nil || k.KafkaConfig == nil || value == nil {
		return
	}
	p, err := sarama.NewSyncProducer(k.KafkaConfig.addressList, k.KafkaConfig.config)
	if err != nil {
		fmt.Println("sarama.NewSyncProducer err")
		return
	}

	defer p.Close()
	msg := &sarama.ProducerMessage{
		Topic: k.KafkaConfig.topic,
		Value: sarama.ByteEncoder(value),
	}
	_, _, err = p.SendMessage(msg)
	if err != nil {
		fmt.Println("send kafka message err")
	}
}

// AsyncKafkaProducer kafka异步生产者
type AsyncKafkaProducer struct {
	KafkaConfig *KafkaConfig         //共有的kafka生产者配置在这个里面
	producer    sarama.AsyncProducer //异步生产者
	isClose     chan struct{}        // 监听producer是否可以关闭
}

// NewKafkaProducer 创建kafka异步生产者实例，并初始化参数
func (k *AsyncKafkaProducer) NewKafkaProducer(address, topic string, duration time.Duration) error {
	if len(address) == 0 {
		return kafkaAddressError
	}
	//配置异步producer启动参数
	k.KafkaConfig = NewProducerByMessage(address, topic, duration)
	k.isClose = make(chan struct{}, 2)
	//启动kafka异步producer
	k.Run()
	return nil
}

// Send kafka异步发送消息
func (k *AsyncKafkaProducer) Send(value []byte) {
	//如果实例或者配置为空，直接返回。如果发送数据为空也直接返回
	if k == nil || k.KafkaConfig == nil || value == nil {
		return
	}
	// 封装消息实例
	msg := &sarama.ProducerMessage{
		Topic: k.KafkaConfig.topic,
		Value: sarama.ByteEncoder(value),
	}
	//这里一般不会出现，producer实例为空时，表示创建异步producer失败
	if k.producer == nil {
		k.Run()
	}
	select {
	//producer 出现error需要重新启动
	case <-k.isClose:
		if k.producer != nil {
			// 收到可以关闭producer的消息isClose,关闭producer并重启
			k.producer.Close()
			k.Run()
		}
		//直接返回，此条消息浪费掉了，如果后期需要收集未发送成功的消息可以在此收集，输出到日志或者等待producer重启成功后再发送
		return
	default:
		// 正常情况发送消息
		k.producer.Input() <- msg
	}
}

// Run kafka异步生产者
func (k *AsyncKafkaProducer) Run() {
	if k == nil || k.KafkaConfig == nil {
		return
	}
	//创建异步producer
	producer, err := sarama.NewAsyncProducer(k.KafkaConfig.addressList, k.KafkaConfig.config)
	//如果创建失败主动置空k.producer，否则producer不为空，在重启的时候k.producer是会有值的
	if err != nil {
		k.isClose <- struct{}{}
		k.producer = nil
		fmt.Println("sarama.NewAsyncProducer err")
		return
	}
	if producer == nil {
		k.isClose <- struct{}{}
		k.producer = nil
		fmt.Println("sarama.NewSyncProducer is null")
		return
	}
	//如果创建成功为实例k的prodcer赋值
	k.producer = producer
	go func(p sarama.AsyncProducer) {
		errors := p.Errors()
		success := p.Successes()
		for {
			select {
			//出现了错误
			case rc := <-errors:
				if rc != nil {
					//标记producer出现error，在send时会监听到这个标记
					k.isClose <- struct{}{}
					fmt.Println("send kafka data error")
				}
				return
			case res := <-success:
				data, _ := res.Value.Encode()
				fmt.Printf("发送成功，value=%s \n", string(data))
			}
		}
	}(producer)

}
