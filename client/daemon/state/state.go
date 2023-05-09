package state

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/42LoCo42/emo/api"
	"github.com/42LoCo42/emo/client/util"
	"github.com/42LoCo42/emo/shared"
	"github.com/gen2brain/go-mpv"
)

type State struct {
	Client *api.Client
	Stats  []api.Stat
	Deltas []api.Stat

	Queue []string

	CurrentSong string
	Time        float64
	Percentage  int64
	Paused      bool

	Mpv *util.Mpv
}

func New() (state *State, err error) {
	// initial state
	state = &State{
		Client: util.Client(),
		Paused: true,
		Mpv: &util.Mpv{
			Mpv: mpv.Create(),
		},
	}

	// get stats from server
	resp, err := state.Client.GetStatsUser(context.Background())
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, shared.Wrap(err, "get stats for user request failed")
	}

	// decode response
	data, err := api.ParseGetStatsUserUserResponse(resp)
	if err != nil {
		return nil, shared.Wrap(err, "could not parse get stats for user response")
	}

	// assign stats
	state.Stats = *data.JSON200

	// create empty deltas - we just need the ID for association
	state.Deltas = make([]api.Stat, len(state.Stats))
	for i, stat := range state.Stats {
		state.Deltas[i] = api.Stat{ID: stat.ID}
	}

	// init MPV
	if err := state.Mpv.Initialize(); err != nil {
		return nil, shared.Wrap(err, "could not initialize MPV")
	}

	// observe some properties
	state.Mpv.Observe("path", mpv.FORMAT_STRING, func(a any) {
		state.CurrentSong = a.(string)
	})
	state.Mpv.Observe("playback-time", mpv.FORMAT_DOUBLE, func(a any) {
		state.Time = a.(float64)
	})
	state.Mpv.Observe("percent-pos", mpv.FORMAT_INT64, func(a any) {
		state.Percentage = a.(int64)
	})
	state.Mpv.Observe("pause", mpv.FORMAT_FLAG, func(a any) {
		state.Paused = a.(bool)
	})

	// stop handler
	state.Mpv.OnStop = func(reason int) {
		switch reason {
		case util.STOP_REASON_EOF:
			log.Print("normal stop, calling next")
			state.NextSong()
		case util.STOP_REASON_STOP:
			log.Print("early stop, no action")
		case util.STOP_REASON_ERROR:
			log.Print("ERROR - halting daemon!")
		}
	}

	// we start paused
	state.SetPaused(true)
	return state, nil
}

func (state *State) SyncStats() error {
	panic("TODO")
}

func (state *State) SelectSong(song string) {
	state.CurrentSong = song
	state.Mpv.Command([]string{"loadfile", song})
	state.SetPaused(false)
	log.Print("Now playing: ", song)
}

func (state *State) NextSong() string {
	var song string

	if len(state.Queue) > 0 {
		song = state.Queue[0]
		state.Queue = state.Queue[1:]
	} else {
		song = util.RandomStat(&state.Stats).Song
	}

	state.SelectSong(song)
	return song
}

func (state *State) SetPaused(paused bool) {
	state.Paused = paused
	state.Mpv.SetProperty("pause", mpv.FORMAT_FLAG, paused)
}

func (state *State) Move(time float64) {
	state.Mpv.Command([]string{"seek", fmt.Sprint(time)})
}
