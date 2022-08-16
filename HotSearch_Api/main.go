package main

import (
	"HotSearch_Api/models"
	"HotSearch_Api/pkg/logger"
	"HotSearch_Api/pkg/sentinel"
	"HotSearch_Api/router"
	"context"
	"fmt"
	"github.com/go-micro/plugins/v4/registry/etcd"
	"github.com/sirupsen/logrus"
	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/metadata"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/web"
)

type logWrapper struct {
	client.Client
}

func (l *logWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	md, _ := metadata.FromContext(ctx)
	fmt.Printf("[Log Wrapper] ctx: %v service: %s method: %s\n", md, req.Service(), req.Endpoint())
	return l.Client.Call(ctx, req, rsp)
}

func newLogWrapper(c client.Client) client.Client {
	return &logWrapper{c}
}

func main() {
	//初始化日志
	if err := logger.InitLog(logrus.DebugLevel); err != nil {
		panic("log initialization failed")
	}
	sentinel.InitSentinel()
	etcdReg := etcd.NewRegistry(func(options *registry.Options) {
		// consul注册中心地址
		options.Addrs = []string{":2379"}
	})
	app := micro.NewService(
		micro.Registry(etcdReg),
		micro.WrapClient(newLogWrapper),
	)
	Service := models.NewHotSearchService("HotSearch_RPC", app.Client())
	http := web.NewService(
		web.Name("HotSearch_Client"),
		web.Address(":8000"),
		web.Handler(router.InitRouter(Service)),
		web.Metadata(map[string]string{"protocol": "http"}),
		web.Advertise(":8000"),
	)
	http.Init()
	http.Run()
}
