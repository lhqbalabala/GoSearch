package main

import (
	"RelatedSearch_RPC/handler"
	"RelatedSearch_RPC/models"
	"RelatedSearch_RPC/pkg/kafka"
	"RelatedSearch_RPC/pkg/logger"
	"RelatedSearch_RPC/pkg/trie"
	"github.com/go-micro/plugins/v4/registry/etcd"
	"github.com/sirupsen/logrus"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/server"
	"os"
)

const IndexPath = "./pkg/data"

func main() {
	//初始化日志
	if err := logger.InitLog(logrus.DebugLevel); err != nil {
		panic("log initialization failed")
	}
	go kafka.Consumer()
	etcdReg := etcd.NewRegistry(func(options *registry.Options) {
		// consul注册中心地址
		options.Addrs = []string{":2379"}
	})
	app := micro.NewService(
		micro.Name("RelatedSearch_RPC"),
		micro.Address(":8030"),
		micro.Registry(etcdReg),
		func(o *micro.Options) {
			o.Server.Init(server.Advertise(":8030"))
		},
	)
	trie.InitTrie(IndexPath + string(os.PathSeparator) + trie.TrieDataPath)
	err := models.RegisterRelatedSearchHandler(app.Server(), new(handler.RelatedSearch))
	if err != nil {
		logger.Logger.Infoln(err)
	}
	app.Init()

	app.Run()

}
