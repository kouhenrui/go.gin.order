package grpc

import (
	"fmt"
	"go.gin.order/src/config"
	"google.golang.org/grpc"
	"log"
	"net"
)

// 微服务文件生成：protoc --go_out=. --go-grpc_out=. auth.proto
func InitGrpc() {
	log.Println("开启监听tcp", config.Tcp)
	// 加载证书和密钥
	//creeds, err := credentials.NewServerTLSFromFile("../../config/https/server.crt", "../../config/https/server.key")
	//if err != nil {
	//	logger.Fatalf("Failed to generate credentials: %v", err)
	//
	//	panic(err)
	//}
	//grpc.Creds(creeds)
	// 创建 gRPC 服务器
	server := grpc.NewServer()

	// 将server结构体注册到grpc服务中
	//au.RegisterAuthServiceServer(server, &auth.AuthService{})
	//sm.RegisterSmtpServiceServer(server, &smtp.SMTPService{})
	ln, err := net.Listen("tcp", config.Tcp)
	if err != nil {
		fmt.Println("网络异常：", err)
		panic(err)
		//return
	}
	// 监听服务
	err = server.Serve(ln)
	if err != nil {
		fmt.Println("监听异常：", err)
		panic(err)
		//return
	}
	log.Println("tcp listen:", config.Tcp)
}
