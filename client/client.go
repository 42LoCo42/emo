package main

import (
	"bufio"
	"bytes"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/jamesruan/sodium"
	"golang.org/x/crypto/argon2"
	"golang.org/x/term"
)

func main() {
	fmt.Fprint(os.Stderr, "Username: ")
	username, err := bufio.NewReader(os.Stdin).ReadBytes('\n')
	if err != nil {
		log.Fatal(err)
	}
	username = bytes.TrimSpace(username)

	fmt.Fprint(os.Stderr, "Password: ")
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintln(os.Stderr)

	var seed sodium.SignSeed
	seed.Bytes = argon2.IDKey(
		password,
		sha512.New().Sum(append(append(username, 0), password...)),
		uint32(sodium.CryptoPWHashOpsLimitInteractive),
		uint32(sodium.CryptoPWHashMemLimitInteractive >> 10),
		uint8(runtime.NumCPU()),
		uint32(seed.Size()),
	)

	kp := sodium.SeedSignKP(seed)
	signature := base64.URLEncoding.EncodeToString(
		sodium.Bytes(username).SignDetached(kp.SecretKey).Bytes,
	)

	log.Printf("pubkey: %#v", kp.PublicKey.Bytes)
	log.Print("signature: ", signature)
}
