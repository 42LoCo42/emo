package main

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/42LoCo42/emo/api"
	"github.com/42LoCo42/emo/shared"
	"github.com/asdine/storm/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jamesruan/sodium"
)

type Server struct {
	db     *storm.DB
	jwtKey []byte
}

// NewError implements api.Handler
func (s *Server) NewError(ctx context.Context, err error) *api.ErrorStatusCode {
	shared.Trace(err)

	code := http.StatusBadRequest
	msg := err.Error()

	cause := shared.RCause(err)
	if cause == storm.ErrNotFound {
		code = http.StatusNotFound
	}

	return &api.ErrorStatusCode{
		StatusCode: code,
		Response: api.Error{
			Msg: msg,
		},
	}
}

// LoginUserGet implements api.Handler
func (s *Server) LoginUserGet(ctx context.Context, params api.LoginUserGetParams) (r api.LoginUserGetOK, err error) {
	name := params.User
	var user api.User

	if err := s.db.One("ID", name, &user); err != nil {
		return r, shared.Wrap(err, "user %v not found", name)
	}

	pubkey := sodium.BoxPublicKey{
		Bytes: user.PublicKey,
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS512,
		jwt.RegisteredClaims{
			Issuer:   "emo",
			Subject:  string(name),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	)

	signed, err := token.SignedString(s.jwtKey)
	if err != nil {
		return r, shared.Wrap(err, "could not sign token")
	}

	r.Data = strings.NewReader(base64.StdEncoding.EncodeToString(
		sodium.Bytes(signed).SealedBox(pubkey),
	))
	return r, nil
}

// SongsGet implements api.Handler
func (s *Server) SongsGet(ctx context.Context) ([]api.Song, error) {
	var songs []api.Song
	return songs, shared.WrapP(s.db.All(&songs), "could not get songs")
}

// SongsNameDelete implements api.Handler
func (s *Server) SongsNameDelete(ctx context.Context, params api.SongsNameDeleteParams) error {
	song, err := s.SongsNameGet(ctx, api.SongsNameGetParams{
		Name: params.Name,
	})
	if err != nil {
		return shared.Wrap(err, "song not found")
	}

	if err := s.db.DeleteStruct(song); err != nil {
		return shared.Wrap(err, "could not delete song from DB")
	}

	if err := os.Remove(string(song.ID)); err != nil {
		return shared.Wrap(err, "could not delete song file")
	}

	return nil
}

// SongsNameFileGet implements api.Handler
func (s *Server) SongsNameFileGet(ctx context.Context, params api.SongsNameFileGetParams) (r api.SongsNameFileGetOK, err error) {
	song, err := s.SongsNameGet(ctx, api.SongsNameGetParams{
		Name: params.Name,
	})
	if err != nil {
		return r, shared.Wrap(err, "song not found")
	}

	file, err := os.Open(string(song.ID))
	if err != nil {
		return r, shared.Wrap(err, "could not open song file")
	}

	r.Data = file
	return r, nil
}

// SongsNameGet implements api.Handler
func (s *Server) SongsNameGet(ctx context.Context, params api.SongsNameGetParams) (*api.Song, error) {
	var song api.Song
	return &song, shared.WrapP(s.db.One("Name", params.Name, &song), "could not get song")
}

// SongsPost implements api.Handler
func (s *Server) SongsPost(ctx context.Context, req api.OptSongsPostReq) error {
	song := req.Value.Song
	file := req.Value.File

	if song.ID == "" {
		data := make([]byte, 32)
		if _, err := rand.Read(data); err != nil {
			return shared.Wrap(err, "could not generate random ID")
		}

		song.ID = api.SongID(hex.EncodeToString(data))
	}

	if err := s.db.Save(&song); err != nil {
		return shared.Wrap(err, "could not save song in DB")
	}

	out, err := os.Create(string(song.ID))
	if err != nil {
		return shared.Wrap(err, "could not create song file")
	}
	defer out.Close()

	if _, err := io.Copy(out, file.File); err != nil {
		return shared.Wrap(err, "could not write to song file")
	}

	return nil
}

// StatsBulkaddPost implements api.Handler
func (s *Server) StatsBulkaddPost(ctx context.Context, req []api.Stat) error {
	tx, err := s.db.Begin(true)
	if err != nil {
		return shared.Wrap(err, "could not begin transaction")
	}
	defer tx.Rollback()

	for _, delta := range req {
		var stat api.Stat
		if err := tx.One("ID", delta.ID, &stat); err != nil {
			return shared.Wrap(err, "stat not found")
		}

		stat.Count += delta.Count
		stat.Boost += delta.Boost
		stat.Time += delta.Time

		if err := tx.Save(&stat); err != nil {
			return shared.Wrap(err, "could not save stat")
		}
	}

	return shared.WrapP(tx.Commit(), "could not commit transaction")
}

// StatsGet implements api.Handler
func (s *Server) StatsGet(ctx context.Context) ([]api.Stat, error) {
	var stats []api.Stat
	return stats, shared.WrapP(s.db.All(&stats), "could not get stats")
}

// StatsIDDelete implements api.Handler
func (s *Server) StatsIDDelete(ctx context.Context, params api.StatsIDDeleteParams) error {
	stat, err := s.StatsIDGet(ctx, api.StatsIDGetParams{
		ID: params.ID,
	})
	if err != nil {
		return shared.Wrap(err, "stat not found")
	}

	if err := s.db.DeleteStruct(stat); err != nil {
		return shared.Wrap(err, "could not delete stat from DB")
	}

	return nil
}

// StatsIDGet implements api.Handler
func (s *Server) StatsIDGet(ctx context.Context, params api.StatsIDGetParams) (*api.Stat, error) {
	var stat api.Stat
	return &stat, shared.WrapP(s.db.One("ID", params.ID, &stat), "could not get stat")
}

// StatsPost implements api.Handler
func (s *Server) StatsPost(ctx context.Context, req api.OptStat) (*api.Stat, error) {
	return &req.Value, shared.WrapP(s.db.Save(&req.Value), "could not save stat")
}

// StatsSongSongGet implements api.Handler
func (s *Server) StatsSongSongGet(ctx context.Context, params api.StatsSongSongGetParams) ([]api.Stat, error) {
	var stats []api.Stat
	return stats, shared.WrapP(s.db.Find("Song", params.Song, &stats), "could not get stats")
}

// StatsUserGet implements api.Handler
func (s *Server) StatsUserGet(ctx context.Context) ([]api.Stat, error) {
	panic("TODO")
}

// StatsUserUserGet implements api.Handler
func (s *Server) StatsUserUserGet(ctx context.Context, params api.StatsUserUserGetParams) ([]api.Stat, error) {
	var stats []api.Stat
	return stats, shared.WrapP(s.db.Find("User", params.User, &stats), "could not get stats")
}

// UsersGet implements api.Handler
func (s *Server) UsersGet(ctx context.Context) ([]api.User, error) {
	var users []api.User
	return users, shared.WrapP(s.db.All(&users), "could not get users")
}

// UsersNameDelete implements api.Handler
func (s *Server) UsersNameDelete(ctx context.Context, params api.UsersNameDeleteParams) error {
	song, err := s.UsersNameGet(ctx, api.UsersNameGetParams{
		Name: params.Name,
	})
	if err != nil {
		return shared.Wrap(err, "user not found")
	}

	if err := s.db.DeleteStruct(song); err != nil {
		return shared.Wrap(err, "could not delete user from DB")
	}

	return nil
}

// UsersNameGet implements api.Handler
func (s *Server) UsersNameGet(ctx context.Context, params api.UsersNameGetParams) (*api.User, error) {
	var user api.User
	return &user, shared.WrapP(s.db.One("ID", params.Name, &user), "could not get user")
}

// UsersPost implements api.Handler
func (s *Server) UsersPost(ctx context.Context, req api.OptUser) error {
	return shared.WrapP(s.db.Save(&req.Value), "could not save user")
}
