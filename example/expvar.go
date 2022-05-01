package main

import (
	_ "expvar"
	"net/http"
)

func main() {
	// 基于go官方标准库expvar的metric监控
	_ = http.ListenAndServe(":9897", http.DefaultServeMux)
}
