package core

import (
	"github.com/iia-micro-service/go-grpc/config"
	"github.com/iia-micro-service/go-grpc/interceptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"os"
	"os/signal"
	"syscall"
)

/*
 * @desc : trpc核心struct结构体
 */
type Core struct {
	config           *config.Config
	grpcServer       *grpcServer
	httpServer       *httpServer
	metricServer     *MetricServer
	grpcOptions      []grpc.ServerOption
	reflectionStatus bool
	shutDownHook     func()
}

/*
 * @desc : 追加grpc options属性
 */
func (core *Core) AddOption(options grpc.ServerOption) {
	core.grpcOptions = append(core.grpcOptions, options)
}

/*
 * @desc : 设置grpc的reflection状态
 */
func (core *Core) SetReflectionStatus(status bool) {
	core.reflectionStatus = status
}

/*
 * @desc : 添加Unray拦截器，这里使用grpc官方标准ChainUnaryInterceptor
 */
func (core *Core) SetUnaryInterceptors(interceptors []grpc.UnaryServerInterceptor) {
	unrayInterceptorsOpts := grpc.ChainUnaryInterceptor(interceptors...)
	core.grpcOptions = append(core.grpcOptions, unrayInterceptorsOpts)
}

/*
 * @desc : 添加Stream拦截器，这里使用grpc官方标准ChainStreamInterceptor
 */
func (core *Core) SetStreamInterceptors(interceptors []grpc.StreamServerInterceptor) {
	unrayInterceptorsOpts := grpc.ChainStreamInterceptor(interceptors...)
	core.grpcOptions = append(core.grpcOptions, unrayInterceptorsOpts)
}
func (core *Core) GetGrpcServer() *grpcServer {
	return core.grpcServer
}
func (core *Core) GetHttpServer() *httpServer {
	return core.httpServer
}

/*
 * @desc : core.grpcServer.Run()方法中开启了一个新的协程，运行gRPC服务
 */
func (core *Core) Run() error {
	var err error
	// 让gRPC服务跑起来.
	if false == core.config.PortReuse {
		err = core.grpcServer.Run(core.config)
		if err != nil {
			log.Fatalf("tRPC - err on core.Run:%v", err)
		}
	}
	// 让http服务跑起来
	core.httpServer.Run()
	// 运行基于expvar标准库的轻量级metric监控服务
	core.metricServer.Run()
	return err
}

/*
 * @desc : 主协程在阻塞等待保证不退出，等待系统信号来终止服务
 */
func (core *Core) WaitTermination(stopHook func()) {
	waitSignal := make(chan os.Signal, 1)
	signal.Notify(waitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-waitSignal
	// 如果开启了http服务
	if true == core.config.RunHTTP {
		// 结束http服务
		core.httpServer.Stop()
		if false == core.config.PortReuse {
			// 结束grpc服务
			core.grpcServer.Stop()
		}
	} else {
		// 结束grpc服务
		core.grpcServer.Stop()
	}
	core.metricServer.Stop()
	log.Println("tRPC - END")
	if stopHook != nil {
		stopHook()
	}
}

/*
 * @desc : 返回一个trpc结构体指针
 */
func New(config *config.Config) *Core {
	core := Core{}

	// 设置grpc服务的配置项
	// 框架强制开启一元拦截器 与 stream拦截器
	core.SetUnaryInterceptors(interceptor.GetServerUnrayInterceptors())
	core.SetStreamInterceptors(interceptor.GetServerStreamInterceptors())
	// 初始化一个grpc服务
	grpcSvr := NewGrpc(core.grpcOptions)
	if core.reflectionStatus {
		reflection.Register(grpcSvr.GetRawGrpcServer())
	}
	// 初始化一个http服务，通过gateway方式同时实现http服务
	if true == config.RunHTTP {
		httpSvr := NewHttp(config, grpcSvr.GetRawGrpcServer())
		core.httpServer = httpSvr
	}
	// 初始化一个metric服务
	metricServer := NewMetric()
	core.metricServer = metricServer

	// 给core结构体赋值
	core.config = config
	core.grpcServer = grpcSvr
	return &core
}
