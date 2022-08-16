package kafka

import (
	"RelatedSearch_RPC/pkg/logger"
	"RelatedSearch_RPC/pkg/trie"
	"github.com/Shopify/sarama"
)

func Consumer() {
	consumer, err := sarama.NewConsumer([]string{"120.48.100.204:9092", "43.142.141.48:9092", "121.196.207.80:9092"}, nil)
	if err != nil {
		logger.Logger.Infoln("fail to start consumer,err:%v\n", err)
		return
	}
	partitionList, err := consumer.Partitions("SearchRequest") //根据topic取到所有的分区
	if err != nil {
		logger.Logger.Infoln("fail to get list of partition:", err)
		return
	}
	for partition := range partitionList {
		pc, err := consumer.ConsumePartition("SearchRequest", int32(partition), sarama.OffsetNewest)
		if err != nil {
			logger.Logger.Infoln("failed to start consumer for partition %d,err:%v", partition, err)
			return
		}
		defer pc.AsyncClose()
		//异步从每个分区消费消息
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				trie.Tree.InsertString(string(msg.Value))
			}
		}(pc)
	}
	select {}
}
