package main

import (
	"google.golang.org/grpc"
	"hashService/internal/handlers"
	"hashService/pkg/hashService"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", ":3001")
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
