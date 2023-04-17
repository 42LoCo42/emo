package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"github.com/42LoCo42/emo/shared"
	"github.com/jamesruan/sodium"
	"golang.org/x/crypto/argon2"
)

func Login(userID, password []byte) {
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
	log.Printf("Pubkey: %#v", kp.PublicKey.Bytes)

	// perform login request
	encodedToken, err := Request(
		http.MethodPost,
		fmt.Sprintf(
			"%s?%s=%s",
			shared.ENDPOINT_LOGIN,
			shared.PARAM_NAME,
			string(userID),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	// get encrypted token from base64-encoded response body
	encryptedToken, err := base64.StdEncoding.DecodeString(encodedToken)
	if err != nil {
		log.Fatal(encryptedToken)
	}

	// decrypt token with keypair
	tokenBytes, err := sodium.Bytes(encryptedToken).SealedBoxOpen(kp)
	if err != nil {
		log.Fatal(err)
	}

	token = string(tokenBytes)
	log.Print("Token: ", token)
}
