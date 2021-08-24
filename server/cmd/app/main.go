package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"hashServer/pkg/configs"
	"hashServer/pkg/hashService"
	"time"
)

func main() {
	configs.InitConfig()

	cwt, _ := context.WithTimeout(context.Background(), time.Second*5)
	conn, err := grpc.DialContext(cwt, fmt.Sprintf(":%s", configs.Get("HASH_SERVICE_PORT")),
		grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		panic(err)
	}
	defer conn.Close()

	hash := hashService.NewHashServiceClient(conn)

	stringData := &hashService.ListOfStrings{
		Strings: []string{"ddd", "ddd", "sss", "sdsfsd", "wewewe"},
	}

	hashResult, err := hash.GetHash(cwt, stringData)
	if err != nil {
		panic(err)
	}

	fmt.Println(hashResult)
}
