package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/42LoCo42/emo/api"
	"github.com/42LoCo42/emo/shared"
	"github.com/asdine/storm/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jamesruan/sodium"
	"github.com/labstack/echo/v4"
)

type Server struct {
	db     *storm.DB
	jwtKey []byte
}

func (s *Server) getSongByName(name string) (*api.Song, error) {
	var song api.Song
	if err := s.db.One("Name", name, &song); err != nil {
		return nil, shared.Wrap(err, "could not find song")
	}

	return &song, nil
}

// GetLoginUser implements api.ServerInterface
func (s *Server) GetLoginUser(ctx echo.Context, name string) error {
	var user api.User
	if err := s.db.One("Name", name, &user); err != nil {
		return err
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
		return shared.Wrap(err, "could not sign token")
	}

	return ctx.String(http.StatusOK, base64.StdEncoding.EncodeToString(
		sodium.Bytes(signed).SealedBox(pubkey),
	))
}

// DeleteSongsName implements api.ServerInterface
func (s *Server) DeleteSongsName(ctx echo.Context, name string) error {
	song, err := s.getSongByName(name)
	if err != nil {
		return err
	}

	if err := s.db.DeleteStruct(song); err != nil {
		return shared.Wrap(err, "could not delete song from DB")
	}

	if err := os.Remove(song.ID); err != nil {
		return shared.Wrap(err, "could not delete song file")
	}

	return ctx.NoContent(http.StatusOK)
}

// DeleteStatsId implements api.ServerInterface
func (s *Server) DeleteStatsId(ctx echo.Context, id uint64) error {
	return s.db.DeleteStruct(&api.Stat{ID: id})
}

// DeleteUsersName implements api.ServerInterface
func (s *Server) DeleteUsersName(ctx echo.Context, name string) error {
	return s.db.DeleteStruct(&api.User{Name: name})
}

// GetSongs implements api.ServerInterface
func (s *Server) GetSongs(ctx echo.Context) error {
	var songs []api.Song
	if err := s.db.All(&songs); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, songs)
}

// GetSongsName implements api.ServerInterface
func (s *Server) GetSongsName(ctx echo.Context, name string) error {
	song, err := s.getSongByName(name)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, song)
}

// GetSongsNameFile implements api.ServerInterface
func (s *Server) GetSongsNameFile(ctx echo.Context, name string) error {
	song, err := s.getSongByName(name)
	if err != nil {
		return err
	}

	return ctx.File(song.ID)
}

// GetStats implements api.ServerInterface
func (s *Server) GetStats(ctx echo.Context) error {
	var stats []api.Stat
	if err := s.db.All(&stats); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, stats)
}

// GetStatsId implements api.ServerInterface
func (s *Server) GetStatsId(ctx echo.Context, id uint64) error {
	var stat api.Stat
	if err := s.db.One("ID", id, &stat); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, stat)
}

func (s *Server) GetStatsUser(ctx echo.Context) error {
	var stats []api.Stat
	if err := s.db.Find("User", ctx.Get("user"), &stats); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, stats)
}

// GetStatsSongSong implements api.ServerInterface
func (s *Server) GetStatsSongSong(ctx echo.Context, song string) error {
	var stats []api.Stat
	if err := s.db.Find("Song", song, &stats); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, stats)
}

// GetStatsUserUser implements api.ServerInterface
func (s *Server) GetStatsUserUser(ctx echo.Context, user string) error {
	var stats []api.Stat
	if err := s.db.Find("User", user, &stats); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, stats)
}

// GetUsers implements api.ServerInterface
func (s *Server) GetUsers(ctx echo.Context) error {
	var users []api.User
	if err := s.db.All(&users); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, users)
}

// GetUsersName implements api.ServerInterface
func (s *Server) GetUsersName(ctx echo.Context, name string) error {
	var user api.User
	if err := s.db.One("Name", name, &user); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, user)
}

// PostSongs implements api.ServerInterface
func (s *Server) PostSongs(ctx echo.Context) error {
	form, err := ctx.MultipartForm()
	if err != nil {
		return shared.Wrap(err, "could not get multipart form")
	}

	infos := form.Value["Info"]
	if len(infos) != 1 {
		return ctx.NoContent(http.StatusBadRequest)
	}

	var info api.Song
	if err := json.NewDecoder(strings.NewReader(infos[0])).Decode(&info); err != nil {
		return shared.Wrap(err, "could not decode song info")
	}

	orig, err := s.getSongByName(info.Name)
	if err != nil && shared.RCause(err) != storm.ErrNotFound {
		return shared.Wrap(err, "could not get original song")
	}

	if orig != nil {
		info.ID = orig.ID
	} else {
		data := make([]byte, 32)
		if _, err := rand.Read(data); err != nil {
			return shared.Wrap(err, "could not generate random ID")
		}

		info.ID = hex.EncodeToString(data)
	}

	if err := s.db.Save(&info); err != nil {
		return shared.Wrap(err, "could not save song info")
	}

	files := form.File["File"]
	if len(files) == 1 {
		multipartFile, err := files[0].Open()
		if err != nil {
			return shared.Wrap(err, "could not open song file")
		}
		defer multipartFile.Close()

		file, err := os.Create(info.ID)
		if err != nil {
			return shared.Wrap(err, "could not create file")
		}
		defer file.Close()

		if _, err := io.Copy(file, multipartFile); err != nil {
			return shared.Wrap(err, "could not copy multipart body to file")
		}
	}

	return ctx.NoContent(http.StatusOK)
}

// PostStats implements api.ServerInterface
func (s *Server) PostStats(ctx echo.Context) error {
	var stat api.Stat
	if err := json.NewDecoder(ctx.Request().Body).Decode(&stat); err != nil {
		return shared.Wrap(err, "could not decode stat")
	}

	if err := s.db.Save(&stat); err != nil {
		return shared.Wrap(err, "could not save stat")
	}

	return ctx.NoContent(http.StatusOK)
}

// PostStatsBulkadd implements api.ServerInterface
func (s *Server) PostStatsBulkadd(ctx echo.Context) error {
	tx, err := s.db.Begin(true)
	if err != nil {
		return shared.Wrap(err, "could not begin transaction")
	}
	defer tx.Rollback()

	var stats []api.Stat
	if err := json.NewDecoder(ctx.Request().Body).Decode(&stats); err != nil {
		return shared.Wrap(err, "could not decode stats")
	}

	for _, stat := range stats {
		var realStat api.Stat
		if err := tx.One("ID", stat.ID, &realStat); err != nil {
			if err == storm.ErrNotFound {
				log.Printf("Stat with ID %d not found!", stat.ID)
			} else {
				return shared.Wrap(err, "could not get stat")
			}
		}

		realStat.Count += stat.Count
		realStat.Boost += stat.Boost

		if err := tx.Save(&realStat); err != nil {
			return shared.Wrap(err, "could not save stat")
		}
	}

	return tx.Commit()
}

// PostUsers implements api.ServerInterface
func (s *Server) PostUsers(ctx echo.Context) error {
	var user api.User
	if err := json.NewDecoder(ctx.Request().Body).Decode(&user); err != nil {
		return shared.Wrap(err, "could not decode body")
	}

	if err := s.db.Save(&user); err != nil {
		return shared.Wrap(err, "could not save user")
	}

	log.Printf("User %s updated", user.Name)
	return nil
}
