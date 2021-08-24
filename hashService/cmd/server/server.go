package main

import (
	"fmt"
	"google.golang.org/grpc"
	"hashService/internal/handlers"
	"hashService/pkg/configs"
	"hashService/pkg/hashService"
	"net"
)

func main() {
	configs.InitConfig()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.Get("HASH_SERVICE_PORT")))
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	server := &handlers.Server{}

	hashService.RegisterHashServiceServer(s, server)
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}
