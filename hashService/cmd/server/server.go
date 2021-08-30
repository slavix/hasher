package main

import (
	"context"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"hashService/internal/handlers"
	"hashService/internal/interceptors"
	"hashService/pkg/hashService"
	"hashService/pkg/logger"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logrusLogger := logger.Init("grpc-hash-service", 5)

	if err := initConfig(); err != nil {
		logger.Panic("main", "main", err, "error initializing configs")
		panic(err)
	}

	netListener, err := net.Listen("tcp", fmt.Sprintf(":%s", viper.GetString("port")))
	if err != nil {
		logger.Panic("main", "main", err, "net listener fail")
		panic(err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		grpc_logrus.UnaryServerInterceptor(logrus.NewEntry(logrusLogger)), interceptors.ContextRequestInterceptor())))
	server := &handlers.Server{}

	hashService.RegisterHashServiceServer(s, server)
	if err := s.Serve(netListener); err != nil {
		logger.Panic("main", "main", err, "grpc server fail")
		panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	<-c

	if err := shutdown(s, netListener); err != nil {
		logger.Panic("main", "main", err, "shutdown failed")
		panic(err)
	}
}

func shutdown(grpcServer *grpc.Server, netListener net.Listener) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	grpcServer.Stop()

	if err := netListener.Close(); err != nil {
		logger.Panic("main", "shutdown", err, "grps netlistener doesn't close connection")
	}

	logger.Info(ctx, "main", "shutdown", "shutdown success", "")

	return nil
}

func initConfig() error {
	viper.AddConfigPath("/app/configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
