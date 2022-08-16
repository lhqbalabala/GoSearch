package balancer

import (
	"fmt"
	"google.golang.org/grpc/resolver"
	"sync"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/grpclog"
)

// Name is the name of round_robin balancer.
const Name = "consistent_hash_x"
const ConsistentHash = "consistent_hash_x"

var logger = grpclog.Component("consistent_hash_x")

// newBuilder creates a new roundrobin balancer builder.
func newBuilder() balancer.Builder {
	return base.NewBalancerBuilder(Name, &rrPickerBuilder{Name}, base.Config{HealthCheck: true})
}
func InitConsistentHashBuilder(consistanceHashKey string) {
	balancer.Register(newBuilder())
}

func init() {
	balancer.Register(newBuilder())
}

type rrPickerBuilder struct {
	consistentHashKey string
}

func (b *rrPickerBuilder) Build(buildInfo base.PickerBuildInfo) balancer.Picker {

	grpclog.Infof("consistentHashPicker: newPicker called with buildInfo: %v", buildInfo)
	if len(buildInfo.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}

	picker := &consistentHashPicker{
		subConns:          make(map[string]balancer.SubConn),
		hash:              NewKetama(10, nil),
		consistentHashKey: b.consistentHashKey,
	}

	for sc, conInfo := range buildInfo.ReadySCs {
		weight := GetWeight(conInfo.Address)
		for i := 0; i < weight; i++ {
			node := wrapAddr(conInfo.Address.Addr, i)
			picker.hash.Add(node)
			picker.subConns[node] = sc
		}
	}
	return picker
}
func GetWeight(addr resolver.Address) int {
	if addr.Metadata == nil {
		return 1
	}
	weight := addr.Metadata.(int64)
	return int(weight)
}

func (p *consistentHashPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	keys, ok := info.Ctx.Value(ConsistentHash).(string)
	var ret balancer.PickResult

	targetAddr, ok := p.hash.Get(keys)
	if ok {
		ret.SubConn = p.subConns[targetAddr]
	}

	return ret, nil
}

type consistentHashPickerBuilder struct {
	consistentHashKey string
}

func wrapAddr(addr string, idx int) string {
	return fmt.Sprintf("%s-%d", addr, idx)
}

type consistentHashPicker struct {
	subConns          map[string]balancer.SubConn
	hash              *Ketama
	mu                sync.Mutex
	consistentHashKey string
}
