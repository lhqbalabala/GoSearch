package main

import (
	"Search_RPC/handler"
	"Search_RPC/models"
	"Search_RPC/pkg"
	"Search_RPC/pkg/logger"
	"Search_RPC/pkg/redis"
	"flag"
	"fmt"
	"github.com/lhqbalabala/etcd/discovery"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const (
	MyAddress    = ""
	MyPort       = 8050
	kafkaAddress = ":9092"
	EtcdAddress  = ":2379"
)

func main() {
	flag.Parse()
	//初始化日志
	if err := logger.InitLog(logrus.DebugLevel); err != nil {
		panic("log initialization failed")
	}

	redis.InitClient()
	redis.InitPictureRedisClient()
	PkgEngine := pkg.Engine{IndexPath: pkg.DefaultIndexPath()}
	PkgEngine.Init()
	defer PkgEngine.Close()
	pkg.Set(&PkgEngine)
	GrpcEtcdRegisterAndStart()
}

func GrpcEtcdRegisterAndStart() {
	addrs := []string{EtcdAddress}
	etcdRegister := discovery.NewRegister(addrs, zap.NewNop())
	node := discovery.Server{
		Name:   "Search_RPC",
		Addr:   "" + ":" + strconv.Itoa(MyPort),
		Weight: 3,
	}

	server, err := Start()
	if err != nil {
		panic(fmt.Sprintf("start server failed : %v", err))
	}

	if _, err := etcdRegister.Register(node, 100); err != nil {
		panic(fmt.Sprintf("server register failed: %v", err))
	}

	fmt.Println("service started listen on", MyAddress+":"+strconv.Itoa(MyPort))
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			server.Stop()
			etcdRegister.Stop()
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

func Start() (*grpc.Server, error) {
	s := grpc.NewServer()

	models.RegisterSearchServer(s, new(handler.Search))
	lis, err := net.Listen("tcp", MyAddress+":"+strconv.Itoa(MyPort))
	if err != nil {
		return nil, err
	}

	go func() {
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()

	return s, nil
}
