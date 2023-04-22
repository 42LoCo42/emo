package util

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"

	"golang.org/x/term"
)

func InputCreds() ([]byte, []byte) {
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

	return userID, password
}
