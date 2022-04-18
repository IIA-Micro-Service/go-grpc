package core

import (
	"crypto/tls"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/iia-micro-service/go-grpc/config"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"strings"
)

type HttpServer interface {
	GetRawHttpServer()
	GetGatewayMux()
}

type httpServer struct {
	gatewayMux *runtime.ServeMux
	httpServer *http.Server
}

func (hSvr *httpServer) GetRawHttpServer() *http.Server {
	return hSvr.httpServer
}

func (hSvr *httpServer) GetGatewayMux() *runtime.ServeMux {
	return hSvr.gatewayMux
}

func (hSvr *httpServer) Run() {
	go hSvr.serv()
	log.Printf("tRPC - Run Http server on %s\n", hSvr.httpServer.Addr)
}
func (hSvr *httpServer) serv() {
	err := hSvr.httpServer.ListenAndServe()
	if err != nil {
		log.Fatalf("tRPC - Http server err: %v", err)
	}
}

func (hSvr *httpServer) Stop() {
	hSvr.httpServer.Close()
	log.Printf("tRPC - Stop Http server on %s\n", hSvr.httpServer.Addr)
}

func NewHttp(config *config.Config, grpc *grpc.Server) *httpServer {
	var rawHttpServer *http.Server
	var gatewayMux *runtime.ServeMux
	httpServer := &httpServer{}

	// 如果开启了端口复用
	if true == config.PortReuse {
		Addr := config.Ip + ":" + config.GrpcPort
		gatewayMux = NewGateway()
		rawHttpServer = &http.Server{
			//TLSConfig: getTLSConfig(),
			TLSConfig: &tls.Config{InsecureSkipVerify: true},
			Addr:      Addr,
			Handler: h2c.NewHandler(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.ProtoMajor == 2 && strings.Contains(r.Header.Get(`Content-Type`), `application/grpc`) {
						fmt.Println("gRPC request")
						grpc.ServeHTTP(w, r)
					} else {
						fmt.Println("Http request")
						gatewayMux.ServeHTTP(w, r)
					}
				}),
				&http2.Server{}),
		}
	} else {
		Addr := config.Ip + ":" + config.HttpPort
		gatewayMux = NewGateway()
		rawHttpServer = &http.Server{
			//TLSConfig: getTLSConfig(),
			TLSConfig: &tls.Config{InsecureSkipVerify: true},
			Addr:      Addr,
			Handler:   gatewayMux,
		}
	}

	httpServer.gatewayMux = gatewayMux
	httpServer.httpServer = rawHttpServer
	return httpServer
}
