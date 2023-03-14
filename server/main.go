package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aerogo/aero"
	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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
	jwtKeyS, ok := os.LookupEnv("EMO_KEY")
	if !ok {
		die(errors.New("Could not get JWT key from EMO_KEY environment variable"), "")
	}
	jwtKey = []byte(jwtKeyS)

	// connect to database
	db, err = gorm.Open(sqlite.Open("emo.db"))
	if err != nil {
		die(err, "Could not open database")
	}

	// setup database
	if err := db.AutoMigrate(
		&User{},
		&Song{},
	); err != nil {
		die(err, "Automatic migration failed")
	}

	// create admin:admin
	admin := User{
		Name: "admin",
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
	app.Use(logRequests, secureEndpointCheck)

	// serve static files
	app.Get("/*file", func(ctx aero.Context) error {
		return ctx.File("static/" + ctx.Get("file"))
	})

	app.Post("/login", login)

	app.Run()
}
