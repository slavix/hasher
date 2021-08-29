package main

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"hashService/internal/handlers"
	"hashService/pkg/hashService"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	netListener, err := net.Listen("tcp", fmt.Sprintf(":%s", viper.GetString("port")))
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	server := &handlers.Server{}

	hashService.RegisterHashServiceServer(s, server)
	if err := s.Serve(netListener); err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	<-c

	if err := shutdown(s, netListener); err != nil {
		panic(err)
	}
}

func shutdown(grpcServer *grpc.Server, netListener net.Listener) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	grpcServer.Stop()

	if err := netListener.Close(); err != nil {
		log.Println("main", "shutdown", err, "grps netlistener doesn't close connection")
	}

	log.Println(ctx, "main", "shutdown", "shutdown success", "")

	return nil
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
