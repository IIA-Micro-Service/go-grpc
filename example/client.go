package main

import (
	"context"
	pb "github.com/iia-micro-service/go-grpc/example/passport"
	"google.golang.org/grpc"
	"log"
)

func main() {
	ctx := context.Background()
	clientConn, _ := GetClientConn(ctx, "localhost:9999", nil)
	defer clientConn.Close()

	passportServiceClient := pb.NewPassportClient(clientConn)
	resp, _ := passportServiceClient.Login(ctx, &pb.LoginRequest{Name: "Go"})

	log.Printf("resp: %v", resp)
}

func GetClientConn(ctx context.Context, target string, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(opts, grpc.WithInsecure())
	return grpc.DialContext(ctx, target, opts...)
}
