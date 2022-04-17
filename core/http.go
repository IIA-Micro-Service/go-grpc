package core

import (
	"crypto/tls"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/iia-micro-service/go-grpc/config"
	"google.golang.org/grpc"
	"log"
	"net/http"
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
	httpServer := &httpServer{}
	Addr := config.Ip + ":" + config.HttpPort
	//httpMux    := http.NewServeMux()
	gatewayMux := NewGateway()
	//httpMux.Handle("/", gatewayMux)
	rawHttpServer := &http.Server{
		//TLSConfig: getTLSConfig(),
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
		Addr:      Addr,
		Handler:   gatewayMux,
		/*
			Handler: h2c.NewHandler(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.ProtoMajor == 2 &&
						strings.Contains(r.Header.Get(`Content-Type`), `application/grpc`) {
						fmt.Println("gRPC request")
						grpc.ServeHTTP(w, r)
					} else {
						fmt.Println("Http request")
						httpMux.ServeHTTP(w, r)
					}
				}),
				&http2.Server{}),
		*/
	}
	httpServer.gatewayMux = gatewayMux
	httpServer.httpServer = rawHttpServer
	return httpServer
}
