package main

import (
	"HotSearch_RPC/handler"
	"HotSearch_RPC/models"
	"HotSearch_RPC/pkg/logger"
	"HotSearch_RPC/pkg/redis"
	"HotSearch_RPC/pkg/trie"
	"fmt"
	"github.com/go-micro/plugins/v4/registry/etcd"
	"github.com/sirupsen/logrus"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/server"
)

func main() {
	//初始化日志
	if err := logger.InitLog(logrus.DebugLevel); err != nil {
		panic("log initialization failed")
	}
	trie.InitHotSearch("./pkg/data/HotSearch.txt")
	err := redis.InitClient()
	if err != nil {
		fmt.Println(err)
	}
	etcdReg := etcd.NewRegistry(func(options *registry.Options) {
		// consul注册中心地址
		options.Addrs = []string{":2379"}
	})
	app := micro.NewService(
		micro.Name("HotSearch_RPC"),
		micro.Address(":8010"),
		micro.Registry(etcdReg),
		func(o *micro.Options) {
			o.Server.Init(server.Advertise(":8010"))
		},
	)

	err = models.RegisterHotSearchHandler(app.Server(), new(handler.HotSearch))
	if err != nil {
		fmt.Println(err)
	}
	app.Run()
}
