package main

import (
	"net/http"
	"strconv"

	"github.com/42LoCo42/emo/shared"
	"github.com/aerogo/aero"
	"gorm.io/gorm"
)

func getStats(ctx aero.Context) error {
	name := ctx.Request().Internal().FormValue(shared.PARAM_NAME)
	data := db.
		Model(shared.Stat{}).
		Joins("join songs on stats.song_id = songs.id").
		Select("songs.id, songs.name, stats.count, stats.boost").
		Where(shared.Stat{UserID: ctx.Data().(*ContextData).UserID})

	if name == "" {
		// return all stats
		var stats []shared.StatQuery

		if err := data.Find(&stats).Error; err != nil {
			return onDBError(ctx, err, "No stats found!")
		}

		return ctx.JSON(stats)
	} else {
		// return selected stat
		var stat shared.StatQuery

		res := data.Find(&stat, "songs.name = ?", name)
		if err := res.Error; err != nil {
			return onDBError(ctx, err, "get stat failed!")
		}
		if res.RowsAffected == 0 {
			return onDBError(ctx, gorm.ErrRecordNotFound, "Stat not found!")
		}

		return ctx.JSON(stat)
	}
}

func setStat(ctx aero.Context) error {
	// get name, required
	name := ctx.Request().Internal().FormValue(shared.PARAM_NAME)
	if name == "" {
		return ctx.Error(http.StatusBadRequest, "song name needed!")
	}

	// get song, we need the ID later
	song := shared.Song{Name: name}
	if db.Find(&song, song).RowsAffected == 0 {
		return onDBError(ctx, gorm.ErrRecordNotFound, "song not found!")
	}

	// stat with conditions
	stat := shared.Stat{
		UserID: ctx.Data().(*ContextData).UserID,
		SongID: song.ID,
	}
	if err := db.
		// stat attributes for when we create a new one
		Attrs(shared.Stat{
			Count: 0,
			Boost: 1,
		}).
		FirstOrCreate(&stat, stat).Error; err != nil {
		return onDBError(ctx, err, "could not create or load stat!")
	}

	// perform changes

	if count, err := strconv.ParseUint(
		ctx.Request().Internal().FormValue(shared.PARAM_COUNT),
		10, 64,
	); err == nil {
		stat.Count = count
	}

	if boost, err := strconv.ParseUint(
		ctx.Request().Internal().FormValue(shared.PARAM_BOOST),
		10, 64,
	); err == nil {
		stat.Boost = boost
	}

	// give new stat to DB and user
	db.Save(stat)
	return getStats(ctx)
}

func deleteStat(ctx aero.Context) error {
	// get name, required
	name := ctx.Request().Internal().FormValue(shared.PARAM_NAME)
	if name == "" {
		return ctx.Error(http.StatusBadRequest, "song name needed!")
	}

	// get stat of current user and with song ID
	res := db.
		Where("user_id = ?", ctx.Data().(*ContextData).UserID).
		Where("song_id = (?)", db.
			Model(shared.Song{}).
			Select("id").
			Where("name = ?", name),
		).
		Delete(shared.Stat{})
	if err := res.Error; err != nil {
		return onDBError(ctx, err, "delete stat failed!")
	}
	if res.RowsAffected == 0 {
		return onDBError(ctx, gorm.ErrRecordNotFound, "stat not found!")
	}

	return nil
}
