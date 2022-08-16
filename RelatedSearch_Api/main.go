package main

import (
	"RelatedSearch_Api/models"
	"RelatedSearch_Api/pkg/logger"
	"RelatedSearch_Api/router"
	"context"
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
	logger.Logger.Infoln("[Log Wrapper] ctx: %v service: %s method: %s\n", md, req.Service(), req.Endpoint())
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
	etcdReg := etcd.NewRegistry(func(options *registry.Options) {
		// consul注册中心地址
		options.Addrs = []string{":2379"}
	})
	app := micro.NewService(
		micro.Registry(etcdReg),
		micro.WrapClient(newLogWrapper),
	)
	Service := models.NewRelatedSearchService("RelatedSearch_RPC", app.Client())
	http := web.NewService(
		web.Name("RelatedSearch_Client"),
		web.Address(":8020"),
		web.Handler(router.InitRouter(Service)),
		web.Metadata(map[string]string{"protocol": "http"}),
		web.Advertise(":8020"),
	)
	http.Init()
	http.Run()
}
