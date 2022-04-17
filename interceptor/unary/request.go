package unary

import (
	"context"
	"google.golang.org/grpc"
	"log"
)

/*
 * @desc : 每一次请求均会触发的拦截器.
 * @return : grpc.UnaryServerInterceptor是一个func类型的type
 */
func Request() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		log.Println("你好")
		resp, err := handler(ctx, req)
		log.Println("再见")
		return resp, err
	}
}
