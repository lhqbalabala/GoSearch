package sentinel

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/alibaba/sentinel-golang/logging"
	"go-micro.dev/v4/util/log"
)

func InitSentinel() {
	// We should initialize Sentinel first.
	conf := config.NewDefaultConfig()
	// for testing, logging output to console
	conf.Sentinel.Log.Logger = logging.NewConsoleLogger()
	err := sentinel.InitWithConfig(conf)
	if err != nil {
		log.Fatal(err)
	}

	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "GetHotSearch_flow",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Throttling, // 流控效果为匀速排队
			Threshold:              100,             // 请求的间隔控制在 1000/100=10 ms
			MaxQueueingTimeMs:      100,             // 最长排队等待时间
		},
	})
	if err != nil {
		log.Fatalf("Unexpected error: %+v", err)
		return
	}
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "InsertQuery_flow",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Throttling, // 流控效果为匀速排队
			Threshold:              500,             // 请求的间隔控制在 1000/500=2 ms
			MaxQueueingTimeMs:      100,             // 最长排队等待时间
		},
	})
	if err != nil {
		log.Fatalf("Unexpected error: %+v", err)
		return
	}
	_, err = circuitbreaker.LoadRules([]*circuitbreaker.Rule{
		// Statistic time span=5s, recoveryTimeout=3s, slowRtUpperBound=50ms, maxSlowRequestRatio=50%
		{
			Resource:                     "GetHotSearch_break",
			Strategy:                     circuitbreaker.SlowRequestRatio,
			RetryTimeoutMs:               3000, //熔断触发后持续的时间
			MinRequestAmount:             10,   //触发熔断的最小请求数目
			StatIntervalMs:               5000, //统计的时间窗口长度
			StatSlidingWindowBucketCount: 10,
			MaxAllowedRtMs:               50,  //慢调用时间
			Threshold:                    0.5, //慢调用比例50%
		},
	})
	if err != nil {
		log.Fatalf("Unexpected error: %+v", err)
		return
	}
}
