package main

import (
	"go.gin.order/src/config"
	"go.gin.order/src/internal/api/rest"
)

func main() {

	config.InitConfig()
	//go grpc.InitGrpc()
	rest.InitHttp()

}
