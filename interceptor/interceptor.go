package interceptor

import (
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/iia-micro-service/go-grpc/interceptor/unary"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func requestErrorHandler(p interface{}) (err error) {
	log.Printf("tRPC : error in request +%v", err)
	return status.Errorf(codes.Internal, "Something went wrong :( ")
}

/*
 * @desc : 获取所有一元拦截器（服务器侧）
 */
func GetServerUnrayInterceptors() []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{
		unary.Request(),
		grpcRecovery.UnaryServerInterceptor(grpcRecovery.WithRecoveryHandler(requestErrorHandler)),
	}
}

/*
 * @desc : 获取所有stream拦截器（服务器侧）
 */
func GetServerStreamInterceptors() []grpc.StreamServerInterceptor {
	return []grpc.StreamServerInterceptor{
		grpcRecovery.StreamServerInterceptor(grpcRecovery.WithRecoveryHandler(requestErrorHandler)),
	}
}
