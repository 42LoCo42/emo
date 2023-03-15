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

	KEY_IS_AUTHED = "is_authed"
	KEY_IS_ADMIN  = "is_admin"
	KEY_USERID    = "userid"

	PARAM_NAME = "name"
)
