package main

import (
	"context"
	"github.com/iia-micro-service/go-grpc/config"
	"github.com/iia-micro-service/go-grpc/core"
	pb "github.com/iia-micro-service/go-grpc/example/passport"
	"log"
)

type PassportService struct{}

func (s *PassportService) Login(ctx context.Context, r *pb.LoginRequest) (*pb.LoginReply, error) {
	return &pb.LoginReply{Message: "hello.world"}, nil
}

func main() {
	// 初始化配置
	trpcConfig := config.Config{}
	trpcConfig.Ip = "0.0.0.0"
	trpcConfig.GrpcPort = "9999"

	// New方法获取一个svr实例
	svr := core.New(&trpcConfig)

	// 注入pb服务
	pb.RegisterPassportServer(svr.GetGrpcServer().GetRawGrpcServer(), &PassportService{})

	// 执行svr实例的Run()方法，Run()方法中用新的协程去运行go grpc标准库服务
	err := svr.Run()
	if err != nil {
		log.Fatalf("tRPC : Run error %v\n", err)
	}

	// "主协程"阻塞挂起，保证不退出，并同时等待系统结束信号即sigint或者sigterm信号
	svr.WaitTermination(nil)
}