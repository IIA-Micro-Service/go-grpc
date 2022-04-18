package config

type Config struct {
	Ip            string // 服务器IP
	GrpcPort      string // gRPC服务端口号
	PortReuse     bool   // http服务与gRPC服务是否端口复用，仅在开启RunHTTP为true时候有效
	RunHTTP       bool   // 是否开启http服务
	HttpPort      string // http服务所在端口号，仅在开启RunHTTP为true时且PortReuse为false时有效
	TLSKey        string // 证书
	TLSPerm       string // 证书
	RunReflection bool   // 是否开启gRPC服务反射
}
