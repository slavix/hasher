package handlers

import (
	"context"
	"encoding/hex"
	"golang.org/x/crypto/sha3"
	"hashService/pkg/hashService"
	"sync"
)

type Server struct {
	ListOfStrings []*hashService.ListOfStrings
	hashService.UnimplementedHashServiceServer
}

func (s *Server) GetHash(ctx context.Context, in *hashService.ListOfStrings) (*hashService.ListOfHashes, error) {
	hashes := make([]*hashService.Hash, len(in.Strings))

	var wg sync.WaitGroup

	for i, str := range in.Strings {
		wg.Add(1)
		go func(i int, s string) {
			defer wg.Done()

			h := sha3.New256()
			h.Write([]byte(s))
			sha3Hash := hex.EncodeToString(h.Sum(nil))

			hashes[i] = &hashService.Hash{
				Str:  s,
				Hash: sha3Hash,
			}
		}(i, str)
	}

	wg.Wait()

	return &hashService.ListOfHashes{
		Data: hashes,
	}, nil
}
