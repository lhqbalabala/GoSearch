package main

import (
	"Search_Api/pkg/db"
	go_micro "Search_Api/pkg/go-micro"
	"Search_Api/pkg/grpcExpand"
	balancer "Search_Api/pkg/grpcExpand/balance"
	"Search_Api/pkg/logger"
	"Search_Api/router"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	//初始化日志
	if err := logger.InitLog(logrus.DebugLevel); err != nil {
		panic("log initialization failed")
	}
	//加载配置文件
	//环境变量>配置文件>默认值
	viper.SetConfigFile(".env")
	viper.AutomaticEnv() //自动匹配环境
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Fatal error config file:%s", err)
		return
	}
	db.Init()
	balancer.InitConsistentHashBuilder("consistent_hash_x")
	grpcExpand.InitEtcdClient()
	engine := router.InitRouter()
	go_micro.Init()
	go_micro.Run(engine)

}
