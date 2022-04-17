package core

import (
	"github.com/iia-micro-service/go-grpc/config"
	"google.golang.org/grpc"
	"log"
	"net"
)

/*
 * @desc : 声明一个gRPC interface接口
 */
type GrpcServer interface {
	Run() error
	Stop()
	RegisterService()
	GetListenSocket() net.Listener
	GetRawGrpcServer() *grpc.Server
}

/*
 * @desc : grpcServer实现了GrpcServer接口
 */
type grpcServer struct {
	server       *grpc.Server
	listenSocket net.Listener
}

/*
 * @desc : 开启一个新的协程去run gRPC标准库
 */
func (grpcServer *grpcServer) Run(config *config.Config) error {
	listenSocketAddr := config.Ip + ":" + config.GrpcPort
	listenSocket, err := net.Listen("tcp", listenSocketAddr)
	if err != nil {
		log.Fatalf("tRPC - gRPC net.Listen err: %v", err)
		return err
	}
	grpcServer.listenSocket = listenSocket
	go grpcServer.serv()
	log.Printf("tRPC - Run gRPC server on %s\n", listenSocketAddr)
	return nil
}
func (grpcServer *grpcServer) serv() {
	err := grpcServer.server.Serve(grpcServer.listenSocket)
	if err != nil {
		log.Fatalf("tRPC - gRPC server err: %v", err)
	}
}

/*
 * @desc : grpcServer实现在后台阻塞等待结束服务信号
 */
func (grpcServer *grpcServer) Stop() {
	log.Println("tRPC - Stop gRPC server on ", grpcServer.listenSocket.Addr())
	grpcServer.server.GracefulStop()
	log.Println("tRPC - Close listen socket...")
	grpcServer.listenSocket.Close()
}
func (grpcServer *grpcServer) RegisterService() {

}

/*
 * @desc : 返回用于监听的socket指针
 */
func (grpcServer *grpcServer) GetListenSocket() net.Listener {
	return grpcServer.listenSocket
}

/*
 * @desc : 返回go标准库的grpc server指针
 */
func (grpcServer *grpcServer) GetRawGrpcServer() *grpc.Server {
	return grpcServer.server
}

func NewGrpc(grpcOptions []grpc.ServerOption) *grpcServer {
	grpcServer := &grpcServer{}
	rawGrpc := grpc.NewServer(grpcOptions...)
	grpcServer.server = rawGrpc
	return grpcServer
}
