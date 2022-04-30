package unary

import (
	"context"
	"google.golang.org/grpc"
)

/*
 * @desc : 每一次请求均会触发的拦截器.
 * @return : grpc.UnaryServerInterceptor是一个func类型的type
 */
func RequestAudit() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		// 具体逻辑...
		//log.Println("你好")
		resp, err := handler(ctx, req)
		//log.Println("再见")
		return resp, err

	}
}
