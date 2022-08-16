package go_micro

import (
	"github.com/gin-gonic/gin"
	"github.com/go-micro/plugins/v4/registry/etcd"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/selector"
	"go-micro.dev/v4/web"
)

var consulReg registry.Registry
var sel selector.Selector

func Init() {
	consulReg = etcd.NewRegistry(func(options *registry.Options) {
		// consul注册中心地址
		options.Addrs = []string{":2379"}
	})
	sel = selector.NewSelector(
		selector.Registry(consulReg),
		selector.SetStrategy(selector.RoundRobin),
	)

}
func Run(engine *gin.Engine) {
	http := web.NewService(
		web.Name("Search_Api"),
		web.Address(":8040"),
		web.Registry(consulReg),
		web.Handler(engine),
		web.Advertise(":8040"),
	)
	http.Init()
	http.Run()
}
func Getsel() selector.Selector {
	return sel
}
