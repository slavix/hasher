package main

import (
	"fmt"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"hashService/internal/handlers"
	"hashService/pkg/hashService"
	"log"
	"net"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", viper.GetString("port")))
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

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
