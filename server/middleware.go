package main

import (
	"log"
	"net/http"

	"github.com/42LoCo42/emo/shared"
	"github.com/asdine/storm/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func errorHandler(err error, c echo.Context) {
	if errors.Unwrap(errors.Unwrap(err)) == storm.ErrNotFound {
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
			tokenString := c.Request().Header.Get(shared.AUTH_HEADER)
			if len(tokenString) == 0 {
				return next(c)
			}

			claims := jwt.RegisteredClaims{}
			token, err := jwt.ParseWithClaims(
				tokenString,
				&claims,
				func(t *jwt.Token) (interface{}, error) {
					return s.jwtKey, nil
				},
			)
			if err != nil {
				return errors.Wrap(err, "Could not parse JWT")
			}

			if !token.Valid {
				return errors.New("JWT is not valid!")
			}

			log.Print("JWT is valid!")

			return next(c)
		}
	}
}
