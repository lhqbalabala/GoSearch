package grpcExpand

import (
	"User_Api/pkg/models"
	"context"
	"fmt"
	"github.com/lhqbalabala/etcd/discovery"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/resolver"
)

var client models.SearchClient

const App = "Search_RPC"

func InitEtcdClient(opts ...grpc.DialOption) error {
	options := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, "consistent_hash_x")),
		//grpc.WithBalancerName(roundrobin.Name),
	}

	addrs := []string{"http://43.142.141.48:2379"}
	r := discovery.NewResolver(addrs, zap.NewNop())
	resolver.Register(r)

	conn, err := grpc.Dial("etcd:///"+App, options...)
	if err != nil {
		return err
	}
	client = models.NewSearchClient(conn)
	_, err = Client().GetSearchResult(context.WithValue(context.Background(), roundrobin.Name, "123"), nil)
	fmt.Println(err)
	return nil
}
func Client() models.SearchClient {
	return client
}
