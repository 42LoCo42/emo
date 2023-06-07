package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/42LoCo42/emo/shared"
	"github.com/asdine/storm/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func errorHandler(err error, c echo.Context) {
	if shared.RCause(err) == storm.ErrNotFound {
		log.Print(err)
		c.NoContent(http.StatusNotFound)
	} else {
		log.Printf(
			"Error in %s %s: %s",
			c.Request().Method,
			c.Request().URL,
			err,
		)
		shared.Trace(err)
		c.NoContent(http.StatusInternalServerError)
	}
}

func logRequest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Printf(
			"%s: %s %s",
			c.RealIP(),
			c.Request().Method,
			c.Request().URL,
		)
		return next(c)
	}
}

func authHandler(s *Server) func(echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if strings.HasPrefix(c.Request().URL.Path, "/login/") {
				log.Print("Login request, skipping auth check")
				return next(c)
			}

			tokenString := c.Request().Header.Get(shared.AUTH_HEADER)
			claims := jwt.RegisteredClaims{}
			token, err := jwt.ParseWithClaims(
				tokenString,
				&claims,
				func(t *jwt.Token) (interface{}, error) {
					return s.jwtKey, nil
				},
			)
			if err != nil {
				log.Print(shared.Wrap(err, "Could not parse JWT"))
				return c.NoContent(http.StatusForbidden)
			}

			if !token.Valid {
				log.Print("JWT is not valid!")
				return c.NoContent(http.StatusForbidden)
			}

			log.Printf("JWT is valid for %s!", claims.Subject)
			c.Set("user", claims.Subject)

			return next(c)
		}
	}
}
