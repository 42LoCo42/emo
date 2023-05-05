package main

import (
	"encoding/base64"
	"net/http"
	"time"

	"github.com/42LoCo42/emo/api"
	"github.com/asdine/storm/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jamesruan/sodium"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type Server struct {
	db     *storm.DB
	jwtKey []byte
}

// GetLoginUser implements api.ServerInterface
func (s *Server) GetLoginUser(ctx echo.Context, name string) error {
	var user api.User
	if err := s.db.One("Name", name, &user); err != nil {
		return errors.Wrap(err, "user lookup failed")
	}

	pubkey := sodium.BoxPublicKey{
		Bytes: user.PublicKey,
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS512,
		jwt.RegisteredClaims{
			Issuer:   "emo",
			Subject:  name,
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	)

	signed, err := token.SignedString(s.jwtKey)
	if err != nil {
		return errors.Wrap(err, "could not sign token")
	}

	return ctx.String(http.StatusOK, base64.StdEncoding.EncodeToString(
		sodium.Bytes(signed).SealedBox(pubkey),
	))
}

// DeleteSongsName implements api.ServerInterface
func (s *Server) DeleteSongsName(ctx echo.Context, name string) error {
	panic("unimplemented")
}

// DeleteStatsId implements api.ServerInterface
func (s *Server) DeleteStatsId(ctx echo.Context, id uint64) error {
	panic("unimplemented")
}

// DeleteUsersName implements api.ServerInterface
func (s *Server) DeleteUsersName(ctx echo.Context, name string) error {
	panic("unimplemented")
}

// GetSongs implements api.ServerInterface
func (s *Server) GetSongs(ctx echo.Context) error {
	panic("unimplemented")
}

// GetSongsName implements api.ServerInterface
func (s *Server) GetSongsName(ctx echo.Context, name string) error {
	panic("unimplemented")
}

// GetSongsNameFile implements api.ServerInterface
func (s *Server) GetSongsNameFile(ctx echo.Context, name string) error {
	panic("unimplemented")
}

// GetStats implements api.ServerInterface
func (s *Server) GetStats(ctx echo.Context) error {
	panic("unimplemented")
}

// GetStatsId implements api.ServerInterface
func (s *Server) GetStatsId(ctx echo.Context, id uint64) error {
	panic("unimplemented")
}

// GetStatsSongSong implements api.ServerInterface
func (s *Server) GetStatsSongSong(ctx echo.Context, song string) error {
	panic("unimplemented")
}

// GetStatsUserUser implements api.ServerInterface
func (s *Server) GetStatsUserUser(ctx echo.Context, user string) error {
	panic("unimplemented")
}

// GetUsers implements api.ServerInterface
func (s *Server) GetUsers(ctx echo.Context) error {
	panic("unimplemented")
}

// GetUsersName implements api.ServerInterface
func (s *Server) GetUsersName(ctx echo.Context, name string) error {
	panic("unimplemented")
}

// PostSongs implements api.ServerInterface
func (s *Server) PostSongs(ctx echo.Context) error {
	panic("unimplemented")
}

// PostStats implements api.ServerInterface
func (s *Server) PostStats(ctx echo.Context) error {
	panic("unimplemented")
}

// PostUsers implements api.ServerInterface
func (s *Server) PostUsers(ctx echo.Context) error {
	panic("unimplemented")
}
