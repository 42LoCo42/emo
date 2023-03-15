package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/42LoCo42/emo/shared"
	"github.com/jamesruan/sodium"
	"golang.org/x/crypto/argon2"
	"golang.org/x/term"
)

func main() {
	// input user ID
	fmt.Fprint(os.Stderr, "User ID: ")
	userID, err := bufio.NewReader(os.Stdin).ReadBytes('\n')
	if err != nil {
		log.Fatal(err)
	}
	userID = bytes.TrimSpace(userID)

	// input password
	fmt.Fprint(os.Stderr, "Password: ")
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintln(os.Stderr)

	// create keypair seed via strong hashing
	var seed sodium.BoxSeed
	seed.Bytes = argon2.IDKey(
		password,
		userID,
		uint32(sodium.CryptoPWHashOpsLimitInteractive),
		uint32(sodium.CryptoPWHashMemLimitInteractive>>10),
		1,
		uint32(seed.Size()),
	)

	// create keypair from seed
	kp := sodium.SeedBoxKP(sodium.BoxSeed(seed))
	log.Printf("pubkey: %#v", kp.PublicKey.Bytes)

	// perform login request
	resp, err := http.PostForm(
		"http://localhost:37812/login",
		url.Values{
			shared.PARAM_NAME: []string{string(userID)},
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	// get response body
	buf := new(strings.Builder)
	if _, err := io.Copy(buf, resp.Body); err != nil {
		log.Fatal(err)
	}

	// get encrypted token from base64-encoded response body
	encryptedToken, err := base64.StdEncoding.DecodeString(buf.String())
	if err != nil {
		log.Fatal(encryptedToken)
	}

	// decrypt token with keypair
	token, err := sodium.Bytes(encryptedToken).SealedBoxOpen(kp)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(string(token))
}
