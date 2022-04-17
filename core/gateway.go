package core

import (
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

func NewGateway() *runtime.ServeMux {
	gatewayMux := runtime.NewServeMux()
	return gatewayMux
}
