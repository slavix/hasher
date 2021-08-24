package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"hashServer/pkg/hashService"
	"time"
)

func main() {
	cwt, _ := context.WithTimeout(context.Background(), time.Second*5)
	conn, err := grpc.DialContext(cwt, "localhost:3001", grpc.WithInsecure(), grpc.WithBlock())

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
