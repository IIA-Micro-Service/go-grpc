package unary

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"time"
)

/*
 * @desc : 当grpc客户端取消请求时候，用这个拦截器捕捉
 */
func RequestCancel() grpc.UnaryServerInterceptor {

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		startTime := time.Now()
		// 执行逻辑...
		resp, err := handler(ctx, req)
		endTime := time.Since(startTime)
		if ctx.Err() == context.Canceled {
			log.Println("request-cancel")
			// 记录客户端取消请求的log
			fmt.Println(endTime)
		}
		return resp, err

	}

}
