package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"time"

	"github.com/42LoCo42/emo/shared"
	"github.com/aerogo/aero"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jamesruan/sodium"
)

func login(ctx aero.Context) error {
	var err error

	// get userID from query
	userID := ctx.Request().Internal().FormValue(shared.PARAM_NAME)
	log.Print("Login attempt from ", userID)

	// get user from database
	user := shared.User{ID: userID}
	if err := db.First(&user, user).Error; err != nil {
		log.Print("No such user: ", err)
		return ctx.Error(http.StatusUnauthorized)
	}

	// get public key from user
	pubkey := sodium.BoxPublicKey{
		Bytes: user.PublicKey,
	}
	if len(pubkey.Bytes) != pubkey.Size() {
		log.Print("Stored public key length invalid")
		return ctx.Error(http.StatusInternalServerError)
	}

	// create & sign token for user
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.RegisteredClaims{
		Issuer:   "emo",
		Subject:  userID,
		IssuedAt: jwt.NewNumericDate(time.Now()),
	})
	ss, err := token.SignedString(jwtKey)
	if err != nil {
		log.Print("Could not create JWT: ", err)
		return ctx.Error(http.StatusInternalServerError)
	}

	// encrypt & send token for user
	return ctx.Text(
		base64.StdEncoding.EncodeToString(
			sodium.Bytes(ss).SealedBox(pubkey),
		),
	)
}
