package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aerogo/aero"
	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/42LoCo42/emo/shared"
)

var (
	jwtKey []byte
	db     *gorm.DB
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func die(err error, msg string) {
	if msg != "" {
		err = errors.Wrap(err, msg)
	}

	log.Print(err)

	if trace, ok := err.(stackTracer); ok {
		for _, frame := range trace.StackTrace() {
			fmt.Fprintf(os.Stderr, "%+v\n", frame)
		}
	}

	os.Exit(1)
}

func main() {
	var err error

	// get JWT key from environment
	jwtKeyS, ok := os.LookupEnv(shared.EMO_KEY_VAR)
	if !ok {
		die(errors.Errorf(
			"Could not get JWT key from %s environment variable",
			shared.EMO_KEY_VAR,
		), "")
	}
	jwtKey = []byte(jwtKeyS)

	// connect to database
	db, err = gorm.Open(sqlite.Open(shared.DATABASE))
	if err != nil {
		die(err, "Could not open database")
	}

	// setup custom associations
	if err := db.SetupJoinTable(shared.User{}, "Songs", shared.Stat{}); err != nil {
		die(err, "Could not setup stat table")
	}

	// setup database
	if err := db.AutoMigrate(
		shared.User{},
		shared.Song{},
	); err != nil {
		die(err, "Automatic database migration failed")
	}

	// create admin:admin
	admin := shared.User{
		ID:   "admin",
		Name: "Administrator",
		PublicKey: []byte{
			0x74, 0x20, 0x61, 0x1b, 0xcc, 0xad, 0x16, 0xcd,
			0x41, 0x5b, 0x74, 0x3b, 0x29, 0xf2, 0x94, 0xf3,
			0xcc, 0xf4, 0xde, 0x4b, 0xe5, 0xb8, 0x18, 0xc1,
			0x53, 0x0a, 0x8a, 0x12, 0x0f, 0x90, 0x4a, 0x59,
		},
		Admin: true,
	}
	db.FirstOrCreate(&admin, admin)

	// create aero application
	app := aero.New()
	app.Config.Ports.HTTP = 37812
	app.Use(
		logRequests,
		parseAuthToken,
		secureEndpointCheck,
		adminEndpointCheck,
	)

	// serve static files
	// app.Get("/*file", func(ctx aero.Context) error {
	// 	return ctx.File(path.Join(shared.STATIC_DIR, ctx.Get("file")))
	// })

	app.Post(shared.ENDPOINT_LOGIN, login)

	// songs endpoint
	app.Get(shared.ENDPOINT_SONGS, getSongs)
	app.Get(shared.ENDPOINT_SONGS+"/*"+shared.PARAM_NAME, getSongFile)
	app.Post(shared.ENDPOINT_SONGS, uploadSong)
	app.Delete(shared.ENDPOINT_SONGS, deleteSong)

	// stats endpoint
	app.Get(shared.ENDPOINT_STATS, getStats)
	app.Post(shared.ENDPOINT_STATS, setStat)
	app.Delete(shared.ENDPOINT_STATS, deleteStat)

	app.Run()
}

func onDBError(
	ctx aero.Context,
	err error,
	extra ...any,
) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ctx.Error(http.StatusNotFound, extra...)
	} else {
		return ctx.Error(http.StatusInternalServerError, extra...)
	}
}
