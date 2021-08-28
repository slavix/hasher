package main

import (
	"github.com/go-openapi/loads"
	"github.com/jessevdk/go-flags"
	"github.com/pressly/goose"
	"github.com/spf13/viper"
	"hashServer/internal/generated/restapi"
	"hashServer/internal/generated/restapi/operations"
	"hashServer/internal/handler"
	"hashServer/internal/repository"
	"log"
	"os"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host: viper.GetString("db.host"),
		Port: viper.GetString("db.port"),
		//Username: os.Getenv("POSTGRES_USER"),
		Username: "db_user",
		//Password: os.Getenv("POSTGRES_PASSWORD"),
		Password: "pwd123",
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}

	err = goose.Up(db.DB, "./internal/migrations")
	if err != nil {
		log.Fatalf("migrations failed: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	handler.InitHandler(repos)

	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewHashServerAPI(swaggerSpec)
	server := restapi.NewServer(api)
	server.Port = viper.GetInt("port")
	defer server.Shutdown()

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
		log.Fatalln(err)
	}

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
