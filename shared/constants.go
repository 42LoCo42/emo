package shared

import "path"

var (
	DATABASE = "emo.db"

	EMO_KEY_VAR = "EMO_KEY"

	STATIC_DIR = "static"
	SONG_DIR   = path.Join(STATIC_DIR, "songs")

	ENDPOINT_LOGIN  = "/login"
	ENDPOINT_ADMIN  = "/admin"
	ENDPOINT_SECURE = "/secure"
	ENDPOINT_SONGS  = ENDPOINT_SECURE + "/songs"
	ENDPOINT_STATS  = ENDPOINT_SECURE + "/stats"

	PARAM_NAME  = "name"
	PARAM_COUNT = "count"
	PARAM_BOOST = "boost"
)
