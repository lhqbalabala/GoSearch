package main

import (
	"User_Api/pkg/db"
	"User_Api/pkg/grpcExpand"
	"User_Api/pkg/logger"
	"User_Api/router"
	"github.com/go-micro/plugins/v4/registry/etcd"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/web"

	"os"
	"os/signal"
)

// @title           lgdSearch API
// @version         1.0
// @description     This is a simple search engine.

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:9090

// @securityDefinitions.apikey  Token
// @in                          header
// @name                        Authorization
// @description					should be set with extra string "Bearer " before it, sample: "Authorization:Bearer XXXXXXXXXXX(token)"
func main() {
	etcdReg := etcd.NewRegistry(func(options *registry.Options) {
		// consul注册中心地址
		options.Addrs = []string{"43.142.141.48:2379"}
	})
	//初始化日志
	if err := logger.InitLog(logrus.DebugLevel); err != nil {
		panic("log initialization failed")
	}
	//初始化etcd
	grpcExpand.InitEtcdClient()
	//加载配置文件
	//环境变量>配置文件>默认值
	viper.SetConfigFile(".env")
	viper.AutomaticEnv() //自动匹配环境
	if err := viper.ReadInConfig(); err != nil {
		logger.Logger.Errorf("Fatal error config file:%s", err)
		return
	}
	engine := router.Init()
	db.Init()
	server := web.NewService(
		web.Name("User_Api"),
		web.Address("10.0.12.9:8060"),
		web.Handler(engine),
		web.Registry(etcdReg),
		web.Metadata(map[string]string{"protocol": "http"}),
		web.Advertise("43.142.141.48:8060"),
	)
	go server.Run()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c
	logger.Logger.Infoln("程序即将结束 -----------------------------")

	logger.Logger.Fatal("program interrupted.\n\n\n")
}
