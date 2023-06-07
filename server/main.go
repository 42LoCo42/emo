package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/42LoCo42/emo/api"
	"github.com/42LoCo42/emo/shared"
	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/codec/msgpack"
	"github.com/labstack/echo/v4"
)

func main() {
	address := flag.String("a", "", "Address to listen on")
	port := flag.Int("p", 37812, "Port to listen on")
	flag.Parse()

	var err error
	server := new(Server)

	server.jwtKey = []byte(os.Getenv(shared.EMO_KEY_VAR))
	if len(server.jwtKey) == 0 {
		shared.Die(nil, shared.EMO_KEY_VAR+" is not set!")
	}

	server.db, err = storm.Open(
		"emo.db",
		storm.Codec(msgpack.Codec),
	)
	if err != nil {
		shared.Die(err, "Could not open database")
	}
	defer server.db.Close()

	admin := api.User{
		CanUploadSongs: true,
		IsAdmin:        true,
		Name:           "admin",
		PublicKey: []byte{
			0x74, 0x20, 0x61, 0x1b, 0xcc, 0xad, 0x16, 0xcd,
			0x41, 0x5b, 0x74, 0x3b, 0x29, 0xf2, 0x94, 0xf3,
			0xcc, 0xf4, 0xde, 0x4b, 0xe5, 0xb8, 0x18, 0xc1,
			0x53, 0x0a, 0x8a, 0x12, 0x0f, 0x90, 0x4a, 0x59,
		},
	}
	server.db.Save(&admin)

	if err := os.MkdirAll("songs", 0755); err != nil {
		shared.Die(err, "could not create songs directory")
	}
	if err := os.Chdir("songs"); err != nil {
		shared.Die(err, "could not chdir to songs directory")
	}

	e := echo.New()
	e.HideBanner = true
	api.RegisterHandlers(e, server)

	e.HTTPErrorHandler = errorHandler
	e.Use(
		logRequest,
		authHandler(server),
	)

	if err := e.Start(fmt.Sprintf("%s:%d", *address, *port)); err != nil {
		log.Fatal(err)
	}
}
