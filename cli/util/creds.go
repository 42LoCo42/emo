package util

import (
	"fmt"
	"os"

	"github.com/jamesruan/sodium"
	"golang.org/x/crypto/argon2"
	"golang.org/x/term"
)

func AskPassword() ([]byte, error) {
	fmt.Fprint(os.Stderr, "Password: ")
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return nil, err
	}
	fmt.Fprintln(os.Stderr)

	return password, nil
}

func MakeKey(username, password []byte) sodium.BoxKP {
	var seed sodium.BoxSeed
	seed.Bytes = argon2.IDKey(
		password,
		username,
		uint32(sodium.CryptoPWHashOpsLimitInteractive),
		uint32(sodium.CryptoPWHashMemLimitInteractive>>10),
		1,
		uint32(seed.Size()),
	)

	return sodium.SeedBoxKP(seed)
}
