package handlers

import (
	"context"
	"fmt"
	"hashService/pkg/hashService"
)

type Server struct {
	ListOfStrings []*hashService.ListOfStrings
	hashService.UnimplementedHashServiceServer
}

func (s *Server) GetHash(ctx context.Context, in *hashService.ListOfStrings) (*hashService.ListOfHashes, error) {
	var hashes hashService.ListOfHashes

	for _, str := range in.Strings {
		fmt.Println(str)
	}

	return &hashes, nil
}
