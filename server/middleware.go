package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/42LoCo42/emo/shared"
	"github.com/aerogo/aero"
	"github.com/golang-jwt/jwt/v5"
)

func logRequests(h aero.Handler) aero.Handler {
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

type ContextData struct {
	Authed bool
	UserID string
}

func parseAuthToken(h aero.Handler) aero.Handler {
	return func(ctx aero.Context) error {
		// deny on true returned
		if func() bool {
			// get token string - empty/no token isn't an error
			tokenS := ctx.Request().Header("authorization")
			if tokenS == "" {
				return false
			}

			// parse token
			token, err := jwt.Parse(tokenS, func(t *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})
			if err != nil {
				log.Print("Parsing token failed: ", err)
				return true
			}

			// check token validity
			if !token.Valid {
				log.Print("Invalid token!")
				return true
			}

			// get user in token
			subject, err := token.Claims.GetSubject()
			if err != nil {
				log.Print("Couldn't get username/subject! ", err)
				return true
			}
			log.Print("User ID: ", subject)

			// authenticate user in current context
			ctx.SetData(ContextData{
				Authed: true,
				UserID: subject,
			})

			return false
		}() {
			return ctx.Error(http.StatusUnauthorized)
		}

		return h(ctx)
	}
}

func authCheck(
	h aero.Handler,
	ctx aero.Context,
	endpoint string,
	adminOly bool,
) error {
	// only trigger on selected endpoint - deny if true returned
	if strings.HasPrefix(ctx.Path(), endpoint) && func() bool {
		log.Print("Access to secure endpoint ", endpoint)

		// deny if not authed
		data, ok := ctx.Data().(ContextData)
		if !ok || !data.Authed {
			log.Printf("Not authenticated!")
			return true
		}

		// if this endpoint requires admin privileges:
		if adminOly {
			// get user from DB
			user := User{ID: data.UserID}

			// deny if user isn't an admin
			if err := db.First(&user).Error; err != nil || !user.Admin {
				log.Printf("User is NOT an administrator!")
				return true
			}

			log.Printf("User is an administrator!")
		}

		return false
	}() {
		return ctx.Error(http.StatusUnauthorized)
	}

	return h(ctx)
}

func secureEndpointCheck(h aero.Handler) aero.Handler {
	return func(ctx aero.Context) error {
		return authCheck(h, ctx, shared.ENDPOINT_SECURE, false)
	}
}

func adminEndpointCheck(h aero.Handler) aero.Handler {
	return func(ctx aero.Context) error {
		return authCheck(h, ctx, shared.ENDPOINT_ADMIN, true)
	}
}
