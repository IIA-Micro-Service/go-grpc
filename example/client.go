package main

import (
	"context"
	"fmt"
	pb "github.com/iia-micro-service/go-grpc/example/passport"
	"google.golang.org/grpc"
	"io"
	"log"
)

func main() {
	ctx := context.Background()
	var opt []grpc.DialOption
	opt = append(opt, grpc.WithInsecure())
	conn, _ := grpc.DialContext(ctx, "localhost:9998", opt...)
	defer conn.Close()

	client := pb.NewPassportClient(conn)
	stream, _ := client.Login(ctx)
	for i := 1; i <= 6; i++ {
		stream.Send(&pb.LoginRequest{
			Name: "clientName",
		})
		resp, err := stream.Recv()
		if io.EOF == err {
			break
		}
		if err != nil {
			fmt.Println("err in client recv:", err)
		}
		log.Println("resp:", resp)
	}
	stream.CloseSend()
}

func GetClientConn(ctx context.Context, target string, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(opts, grpc.WithInsecure())
	return grpc.DialContext(ctx, target, opts...)
}
