package main

import (
	"context"
	"expvar"
	"github.com/iia-micro-service/go-grpc/config"
	"github.com/iia-micro-service/go-grpc/core"
	pb "github.com/iia-micro-service/go-grpc/example/passport"
	"google.golang.org/grpc"
	"io"
	"log"
	"runtime"
)

var loginMethodCounter *expvar.Int
var GOMAXPROCSCounter *expvar.Int
var goroutineCounter *expvar.Int

func init() {
	// 初始化expvar metric计数器
	loginMethodCounter = expvar.NewInt("loginCounter")
	GOMAXPROCSCounter = expvar.NewInt("GOMAXPROCS")
	GOMAXPROCSCounter.Set(int64(runtime.NumCPU()))
	goroutineCounter = expvar.NewInt("goroutine")
	goroutineCounter.Set(int64(runtime.NumGoroutine()))
}

type PassportService struct{}

func (s *PassportService) Login(stream pb.Passport_LoginServer) error {
	//loginCounter := expvar.NewInt("login")
	loginMethodCounter.Add(1)
	n := 0
	for {
		stream.Send(&pb.LoginReply{
			Message: "serverMessage",
		})
		resp, err := stream.Recv()
		if io.EOF == err {
			return nil
		}
		if err != nil {
			return err
		}
		n++
		log.Println("resp:", resp)
		//time.Sleep(time.Duration(1) * time.Second)
	}
	return nil
}

func main() {
	// 初始化配置
	svrConfig := config.Config{}
	svrConfig.Ip = "0.0.0.0"
	svrConfig.GrpcPort = "9998"
	svrConfig.HttpPort = "9999"
	svrConfig.RunHTTP = true
	svrConfig.PortReuse = false

	// New方法获取一个svr实例
	svr := core.New(&svrConfig)

	// 注入pb服务
	pb.RegisterPassportServer(svr.GetGrpcServer().GetRawGrpcServer(), &PassportService{})
	pb.RegisterPassportHandlerFromEndpoint(
		context.Background(),
		svr.GetHttpServer().GetGatewayMux(),
		svrConfig.Ip+":"+svrConfig.GrpcPort,
		[]grpc.DialOption{grpc.WithInsecure()})

	// 执行svr实例的Run()方法，Run()方法中用新的协程去运行go grpc标准库服务
	err := svr.Run()
	if err != nil {
		log.Fatalf("tRPC : Run error %v\n", err)
	}

	// "主协程"阻塞挂起，保证不退出，并同时等待系统结束信号即sigint或者sigterm信号
	svr.WaitTermination(nil)

}
