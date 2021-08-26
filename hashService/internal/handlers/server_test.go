package handlers

import (
	"context"
	"github.com/stretchr/testify/require"
	"hashService/pkg/hashService"
	"testing"
)

var testStrings = &hashService.ListOfStrings{
	Strings: []string{
		"sdfsdf",
		"wewewe",
		"gdfgfdgdfg",
		"ddd",
		"",
	},
}

func TestFalseOnExistKey(t *testing.T) {
	req := require.New(t)

	tests := map[string]struct {
		testStrings []string
		wantHashes  []string
	}{
		"simple": {
			testStrings: []string{"dddd", "sfsdfsf", "werwerwer", ""},
			wantHashes: []string{
				"cc2ba14b88e3b96aa05347427bf7d3ab92d5d1ff2ec7a60bb8be9355b198ba75",
				"9afb658a8ef03a88fcc6f6416dc39126fcaeb3d685d2ed19b60bef42a8044143",
				"4c3992427320e8d4dd01805a1686dd85f771a6303d99f72e71c7a11b26a0ff8d",
				"a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a",
			},
		},
		"empty string also has hash": {
			testStrings: []string{""},
			wantHashes: []string{
				"a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a",
			},
		},
	}

	ctx := context.Background()
	server := &Server{}

	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			resHashes, err := server.GetHash(ctx, &hashService.ListOfStrings{
				Strings: testCase.testStrings,
			})

			req.NoError(err)

			for i, hash := range resHashes.Data {
				req.Equal(testCase.testStrings[i], hash.Str)
				req.Equal(testCase.wantHashes[i], hash.Hash)
			}
		})
	}
}
