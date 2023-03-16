package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/42LoCo42/emo/shared"
	"github.com/aerogo/aero"
	"github.com/akyoto/uuid"
	"gorm.io/gorm"
)

func songPath(song *Song) string {
	return path.Join(shared.SONG_DIR, song.ID)
}

func getSongs(ctx aero.Context) error {
	name := ctx.Request().Internal().FormValue(shared.PARAM_NAME)

	if name == "" {
		// return all songs
		var songs []Song

		if err := db.Find(&songs).Error; err != nil {
			return onDBError(ctx, err, "No songs found!")
		}

		return ctx.JSON(songs)
	} else {
		// return selected song
		song := Song{Name: name}

		if err := db.First(&song, song).Error; err != nil {
			return onDBError(ctx, err, "Song not found")
		}

		return ctx.JSON(song)
	}
}

func getSongFile(ctx aero.Context) error {
	return ctx.File(path.Join(shared.SONG_DIR, ctx.Get(shared.PARAM_NAME)))
}

func uploadSong(ctx aero.Context) error {
	name := ctx.Request().Internal().FormValue(shared.PARAM_NAME)
	song := Song{Name: name}

	if err := db.Transaction(func(tx *gorm.DB) error {
		// create DB entry for song
		if err := tx.
			Where(song).
			Attrs(Song{ID: uuid.New().String()}).
			FirstOrCreate(&song).
			Error; err != nil {
			return err
		}

		// create song folder
		if err := os.MkdirAll(shared.SONG_DIR, 0755); err != nil {
			log.Printf(
				"Could not create song directory %s: %s",
				shared.SONG_DIR, err,
			)
			return err
		}

		// create song file
		path := songPath(&song)
		file, err := os.Create(path)
		if err != nil {
			log.Printf("Could not create song file %s: %s", path, err)
			return err
		}
		defer file.Close()

		// write request body into song file
		if _, err := io.Copy(file, ctx.Request().Internal().Body); err != nil {
			log.Printf("Could not write song file %s: %s", path, err)
			return err
		}

		log.Printf("Successfully uploaded song %s (%s)", name, path)
		return nil
	}); err != nil {
		return ctx.Error(http.StatusInternalServerError, err)
	}

	return ctx.JSON(song)
}

func deleteSong(ctx aero.Context) error {
	name := ctx.Request().Internal().FormValue(shared.PARAM_NAME)
	song := Song{Name: name}

	// find song
	if err := db.First(&song, song).Error; err != nil {
		return onDBError(ctx, err, "Song not found")
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		// delete DB entry for song
		if err := tx.Delete(song).Error; err != nil {
			log.Printf("Could not delete song %s from DB: %s", name, err)
			return err
		}

		// delete song file
		path := songPath(&song)
		if err := os.Remove(path); err != nil {
			log.Printf(
				"Could not delete song's %s file %s: %s",
				name, path, err,
			)
			return err
		}

		return nil
	}); err != nil {
		return ctx.Error(http.StatusInternalServerError, err)
	}

	return nil
}
