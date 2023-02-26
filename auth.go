package main

import (
	"encoding/base64"
	"log"
	"net/http"

	"github.com/aerogo/aero"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jamesruan/sodium"
)

func auth(ctx aero.Context) error {
	var err error

	username := ctx.Request().Internal().PostFormValue("username")
	log.Print("Authentication attempt from ", username)

	var signature sodium.Signature
	signature.Bytes, err = base64.URLEncoding.DecodeString(
		ctx.Request().Internal().PostFormValue("signature"),
	)
	if err != nil {
		log.Print("Could not decode signature: ", err)
		return ctx.Error(http.StatusBadRequest)
	}
	if len(signature.Bytes) != signature.Size() {
		log.Print("Signature length invalid")
		return ctx.Error(http.StatusBadRequest)
	}

	var user User
	if err := db.First(&user, "name = ?", username).Error; err != nil {
		log.Print("No such user: ", err)
		return ctx.Error(http.StatusUnauthorized)
	}

	pubkey := sodium.SignPublicKey{
		Bytes: user.PublicKey,
	}
	if len(pubkey.Bytes) != pubkey.Size() {
		log.Print("Public key length invalid")
		return ctx.Error(http.StatusBadRequest)
	}

	if err := sodium.Bytes(username).SignVerifyDetached(signature, pubkey); err != nil {
		log.Print("Could not verify signature: ", err)
		return ctx.Error(http.StatusUnauthorized)
	}

	log.Print("Success")

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.RegisteredClaims{
		Issuer:    "emo",
		Subject:   username,
	})
	ss, err := token.SignedString(jwtKey)
	if err != nil {
		log.Print("Could not create JWT: ", err)
		return ctx.Error(http.StatusInternalServerError)
	}

	return ctx.Text(ss)
}
