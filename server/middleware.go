package main

import (
	"context"
	"log"
	"strings"

	"github.com/42LoCo42/emo/api"
	"github.com/42LoCo42/emo/shared"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ogen-go/ogen/middleware"
)

func logRequest(req middleware.Request, next middleware.Next) (middleware.Response, error) {
	if req.Raw.Method != "GET" {
		log.Printf(
			"%s: %s %s",
			req.Raw.RemoteAddr,
			req.Raw.Method,
			req.Raw.URL,
		)
	}
	return next(req)
}

func authHandler(s *Server) func(req middleware.Request, next middleware.Next) (middleware.Response, error) {
	return func(req middleware.Request, next middleware.Next) (r middleware.Response, err error) {
		if strings.HasPrefix(req.Raw.URL.Path, "/login/") {
			log.Print("Login request, skipping auth check")
			return next(req)
		}

		tokenString := req.Raw.Header.Get(shared.AUTH_HEADER)
		claims := jwt.RegisteredClaims{}
		token, err := jwt.ParseWithClaims(
			tokenString,
			&claims,
			func(t *jwt.Token) (interface{}, error) {
				return s.jwtKey, nil
			},
		)
		if err != nil {
			return r, shared.Wrap(err, "could not parse JWT")
		}

		if !token.Valid {
			return r, shared.Wrap(nil, "invalid JWT")
		}

		var user api.User
		if err := s.db.One("ID", api.UserName(claims.Subject), &user); err != nil {
			return r, shared.Wrap(err, "could not find user")
		}

		// log.Printf("JWT is valid for %s!", user.Name)
		req.Context = context.WithValue(req.Context, "user", user)
		return next(req)
	}
}
