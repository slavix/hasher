package main

import (
	"context"
	"github.com/go-openapi/loads"
	"github.com/jessevdk/go-flags"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose"
	"github.com/spf13/viper"
	"hashServer/internal/generated/restapi"
	"hashServer/internal/generated/restapi/operations"
	"hashServer/internal/handler"
	"hashServer/internal/repository"
	"hashServer/pkg/logger"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger.Init("swagger-hash-server", 5)

	if err := initConfig(); err != nil {
		logger.Panic("main", "main", err, "error initializing configs")
	}

	// database
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     viper.GetString("db.port"),
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logger.Panic("main", "main", err, "failed to initialize db")
	}

	err = goose.Up(db.DB, "/app/internal/migrations")
	if err != nil {
		logger.Panic("main", "main", err, "migrations failed")
	}

	//configure services
	repos := repository.NewRepository(db)
	handler.InitHandler(repos)

	//swagger server
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		logger.Panic("main", "main", err, "swagger loader failed")
	}

	api := operations.NewHashServerAPI(swaggerSpec)
	server := restapi.NewServer(api)
	server.Host = os.Getenv("APP_HOST")
	server.Port = viper.GetInt("port")

	parser := flags.NewParser(server, flags.Default)
	parser.ShortDescription = "hash_server"
	parser.LongDescription = "Данный сервис должен, взаимодействуя с сервисом считающим хэши (по выбранному вами протоколу), получать из входящих строк их хэши, сохранять их в свою БД (выбор так же за вами) с присвоем id, по которым далее можно будет запрашивать хэши."
	server.ConfigureFlags()
	for _, optsGroup := range api.CommandLineOptionsGroups {
		_, err := parser.AddGroup(optsGroup.ShortDescription, optsGroup.LongDescription, optsGroup.Options)
		if err != nil {
			log.Fatalln(err)
		}
	}

	if _, err := parser.Parse(); err != nil {
		code := 1
		if fe, ok := err.(*flags.Error); ok {
			if fe.Type == flags.ErrHelp {
				code = 0
			}
		}
		os.Exit(code)
	}

	server.ConfigureAPI()

	if err := server.Serve(); err != nil {
		logger.Panic("main", "main", err, "swagger server failed")
	}

	//shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	<-c

	if err := shutdown(db, server); err != nil {
		logger.Panic("main", "main", err, "shutdown failed")
	}
}

func shutdown(db *sqlx.DB, server *restapi.Server) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.Close(); err != nil {
		logger.Panic("main", "shutdown", err, "db doesn't close connection")
	}

	if err := server.Shutdown(); err != nil {
		logger.Panic("main", "shutdown", err, "swagger server doesn't close connection")
	}

	logger.Info(ctx, "main", "shutdown", "shutdown success", "")

	return nil
}

func initConfig() error {
	viper.AddConfigPath("/app/configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
