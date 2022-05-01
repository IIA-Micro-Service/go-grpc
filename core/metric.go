package core

import (
	"crypto/tls"
	"log"
	"net/http"
)

type MetricServer struct {
	httpServer *http.Server
}

func (mSvr *MetricServer) Run() {
	go mSvr.serv()
	log.Printf("tRPC - Run Metric server on %s\n", mSvr.httpServer.Addr)
}
func (mSvr *MetricServer) serv() {
	err := mSvr.httpServer.ListenAndServe()
	if err != nil {
		log.Fatalf("tRPC - Metric server err: %v", err)
	}
}

func (mSvr *MetricServer) Stop() {
	mSvr.httpServer.Close()
	log.Printf("tRPC - Stop Metric server on %s\n", mSvr.httpServer.Addr)
}

func NewMetric() *MetricServer {
	metricServer := &MetricServer{}
	rawHttpServer := &http.Server{
		//TLSConfig: getTLSConfig(),
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
		Addr:      "0.0.0.0:9997",
	}
	metricServer.httpServer = rawHttpServer
	return metricServer
}
