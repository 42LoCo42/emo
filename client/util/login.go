package util

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/42LoCo42/emo/shared"
	"github.com/jamesruan/sodium"
	"golang.org/x/crypto/argon2"
)

var token string = ""

func Token() string {
	if token == "" {
		tokenPath, err := TokenFilePath()
		if err != nil {
			log.Fatal(err)
		}

		tokenBytes, err := os.ReadFile(tokenPath)
		if err != nil {
			log.Fatal(err)
		}

		token = string(tokenBytes)
	}

	return token
}

func TokenFilePath() (string, error) {
	confDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return path.Join(confDir, "emo-token"), nil
}

func Login(userID, password []byte) string {
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

	// perform login request
	encodedToken, err := Request(
		"",
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
	encryptedToken := make([]byte, base64.StdEncoding.DecodedLen(len(encodedToken)))
	if _, err := base64.StdEncoding.Decode(encryptedToken, encodedToken); err != nil {
		log.Fatal(err)
	}

	// decrypt token with keypair
	tokenBytes, err := sodium.Bytes(encryptedToken).SealedBoxOpen(kp)
	if err != nil {
		log.Fatal(err)
	}

	return string(tokenBytes)
}
