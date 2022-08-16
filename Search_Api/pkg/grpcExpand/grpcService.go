package grpcExpand

import (
	"Search_Api/models"
	"fmt"
	"github.com/lhqbalabala/etcd/discovery"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

//var consulClient *api.Client

//func Initconsul() {
//	// 服务发现
//	// 1.初始化 grpcConsul 配置
//	consulConfig := api.DefaultConfig()
//	consulConfig.Address = "localhost:8500"
//	// 2.创建 grpcConsul 对象
//	consulClient, _ = api.NewClient(consulConfig)
//	// 3.服务发现. 从 grpcConsul 上, 获取健康的服务
//	//  参数：
//	//  service: 服务名。 -- 注册服务时，指定该string
//	//  tag：外名/别名。 如果有多个， 任选一个
//	//  passingOnly：是否通过健康检查。 true
//	//  q：查询参数。 通常传 nil
//	//  返回值：
//	//  ServiceEntry： 存储服务的切片。
//	//  QueryMeta：额外查询返回值。 nil
//	//  error： 错误信息
//	consulClient.Health().Connect("queryService", "grpc", true, nil)
//}
//func GetHealthServices() ([]*api.ServiceEntry, error) {
//	services, _, err := consulClient.Health().Service("queryService", "grpc", true, nil)
//	if err != nil {
//		return nil, err
//	}
//	if services == nil {
//		return nil, errors.NewError(500, "没有可用的服务!")
//	}
//	return services, nil
//}
var client models.SearchClient

const App = "Search_RPC"

func InitEtcdClient(opts ...grpc.DialOption) error {
	options := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, "consistent_hash_x")),
		//grpc.WithBalancerName(roundrobin.Name),
	}

	addrs := []string{":2379"}
	r := discovery.NewResolver(addrs, zap.NewNop())
	resolver.Register(r)

	conn, err := grpc.Dial("etcd:///"+App, options...)
	if err != nil {
		return err
	}
	client = models.NewSearchClient(conn)
	return nil
}
func Client() models.SearchClient {
	return client
}
