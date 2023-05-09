package state

import (
	"github.com/42LoCo42/emo/client/util"
	"github.com/42LoCo42/emo/shared"
	"github.com/gen2brain/go-mpv"
)

type State struct {
	CurrentSong string
	Queue       []string
	Time        float64
	Percentage  int64
	Paused      bool
	Mpv         *util.Mpv
}

func New() (state *State, err error) {
	state = &State{
		Paused: true,
		Mpv: &util.Mpv{
			Mpv: mpv.Create(),
		},
	}

	if err := state.Mpv.Initialize(); err != nil {
		return nil, shared.Wrap(err, "could not initialize MPV")
	}

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

	state.Mpv.SetProperty("pause", mpv.FORMAT_FLAG, state.Paused)

	return state, nil
}

func (state *State) SelectSong(song string) {
	state.CurrentSong = song
	state.Mpv.Command([]string{"loadfile", song})
}

func (state *State) SetPaused(paused bool) {
	state.Paused = paused
	state.Mpv.SetProperty("pause", mpv.FORMAT_FLAG, paused)
}

func (state *State) Move(time float64) {

}
