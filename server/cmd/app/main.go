package main

import (
	"fmt"
	"github.com/spf13/viper"
	"hashServer/internal/repository"
	"log"
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

	repos := repository.NewRepository(db)

	//cwt, _ := context.WithTimeout(context.Background(), time.Second*5)
	//conn, err := grpc.DialContext(cwt, fmt.Sprintf(":%s", configs.Get("HASH_SERVICE_PORT")),
	//	grpc.WithInsecure(), grpc.WithBlock())
	//
	//if err != nil {
	//	panic(err)
	//}
	//defer conn.Close()
	//
	//hash := hashService.NewHashServiceClient(conn)
	//
	//stringData := &hashService.ListOfStrings{
	//	Strings: []string{"ddd", "ddd", "sss", "sdsfsd", "wewewe"},
	//}
	//
	//hashResult, err := hash.GetHash(cwt, stringData)
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println(hashResult)
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
