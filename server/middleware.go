package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/aerogo/aero"
	"github.com/golang-jwt/jwt/v5"
)

var (
	logRequests = func(h aero.Handler) aero.Handler {
		return func(ctx aero.Context) error {
			log.Printf(
				"%s: %s %s",
				ctx.Request().Internal().RemoteAddr,
				ctx.Request().Method(),
				ctx.Path(),
			)

			return h(ctx)
		}
	}

	secureEndpointCheck = func(h aero.Handler) aero.Handler {
		return func(ctx aero.Context) error {
			if strings.HasPrefix(ctx.Path(), "/secure/") {
				log.Print("Access to secure endpoint")

				tokenS := ctx.Request().Header("authorization")
				token, err := jwt.Parse(tokenS, func(t *jwt.Token) (interface{}, error) {
					return jwtKey, nil
				})
				if err != nil {
					log.Print("Parsing token failed: ", err)
					return ctx.Error(http.StatusUnauthorized)
				}

				subject, err := token.Claims.GetSubject()
				if err != nil {
					log.Print("Couldn't get username/subject! ", err)
					return ctx.Error(http.StatusUnauthorized)
				}
				log.Print("Username: ", subject)

				if !token.Valid {
					log.Print("Invalid token!")
					return ctx.Error(http.StatusUnauthorized)
				}
			}

			return h(ctx)
		}
	}
)
