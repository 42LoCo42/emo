package state

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/42LoCo42/emo/api"
	"github.com/42LoCo42/emo/daemon/util"
	"github.com/42LoCo42/emo/shared"
	"github.com/gen2brain/go-mpv"
)

const (
	BOOST_SUB_LO = 2
	BOOST_SUB_HI = 20
	BOOST_ADD_LO = 80
)

type State struct {
	Client *api.Client
	Stats  []api.Stat
	Deltas map[api.StatID]api.Stat

	Queue []string

	CurrentFile string
	CurrentStat *api.Stat
	Time        float64
	Percentage  int64
	Paused      bool

	Mpv *util.Mpv
}

func NewState() (state *State, err error) {
	// initial state
	state = &State{
		Client: shared.Client(),
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
	state.Deltas = map[api.StatID]api.Stat{}

	// init MPV
	if err := state.Mpv.Initialize(); err != nil {
		return nil, shared.Wrap(err, "could not initialize MPV")
	}

	// observe some properties
	state.Mpv.Observe("playback-time", mpv.FORMAT_DOUBLE, func(a any) {
		state.Time = a.(float64)
	})
	state.Mpv.Observe("percent-pos", mpv.FORMAT_INT64, func(a any) {
		state.Percentage = a.(int64)
	})

	// not really neccessary
	state.Mpv.Observe("pause", mpv.FORMAT_FLAG, func(a any) {
		state.Paused = a.(bool)
	})

	// stop handler
	state.Mpv.OnStop = func(reason int) {
		switch reason {
		case util.STOP_REASON_EOF:
			log.Print("normal stop, calling next")
			if _, err := state.NextSong(); err != nil {
				log.Print(err)
			}

		case util.STOP_REASON_STOP:
			// do nothing

		case util.STOP_REASON_ERROR:
			log.Print("ERROR - halting daemon!")
		}
	}

	// we start paused
	state.SetPaused(true)
	return state, nil
}

func (state *State) SyncStats() error {
	values := make([]api.Stat, 0, len(state.Deltas))
	for _, v := range state.Deltas {
		values = append(values, v)
	}

	resp, err := state.Client.PostStatsBulkadd(context.Background(), values)
	if err != nil || resp.StatusCode != http.StatusOK {
		return shared.Wrap(err, "stat bulk update failed")
	}

	return nil
}

func (state *State) WithCurrentDelta(f func(*api.Stat)) {
	if state.CurrentStat != nil {
		delta, ok := state.Deltas[state.CurrentStat.ID]
		if !ok {
			delta = api.Stat{
				ID:   state.CurrentStat.ID,
				Song: state.CurrentStat.Song,
				User: state.CurrentStat.User,
			}
		}

		f(&delta)
		state.Deltas[state.CurrentStat.ID] = delta
		log.Printf("%#v", state.Deltas)
	}
}

func (state *State) NextSong() (string, error) {
	var newStat api.Stat

	// pop next song from queue if present, else select random
	if len(state.Queue) > 0 {
		name := state.Queue[0]
		state.Queue = state.Queue[1:]

		for _, stat := range state.Stats {
			if stat.Song == name {
				newStat = stat
				break
			}
		}

		return "", shared.Wrap(nil, fmt.Sprintf("song %s not found!", name))
	} else {
		newStat = util.RandomStat(&state.Stats)
	}

	state.WithCurrentDelta(func(delta *api.Stat) {
		p := state.Percentage
		if p >= BOOST_SUB_LO && p <= BOOST_SUB_HI {
			delta.Boost--
		} else if p >= BOOST_ADD_LO {
			delta.Boost++
			delta.Count++
		}
	})

	state.CurrentStat = &newStat
	if err := state.PlaySong(); err != nil {
		return "", shared.Wrap(err, "could not play song")
	}

	return state.CurrentStat.Song, nil
}

func (state *State) PlaySong() error {
	// download song to temporary file
	tmpFile, err := os.CreateTemp(os.TempDir(), "emo")
	if err != nil {
		return shared.Wrap(err, "could not create temp song file")
	}

	resp, err := state.Client.GetSongsNameFile(
		context.Background(),
		state.CurrentStat.Song,
	)
	if err != nil {
		return shared.Wrap(err, "could not download song")
	}
	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		return shared.Wrap(err, "could not save song")
	}

	// remove old temp file, update state
	os.Remove(state.CurrentFile)
	state.CurrentFile = tmpFile.Name()

	// start playback
	state.Mpv.Command([]string{"loadfile", tmpFile.Name()})
	state.SetPaused(false)
	log.Print("Now playing: ", state.CurrentStat.Song)

	return nil
}

func (state *State) SetPaused(paused bool) {
	state.Paused = paused
	state.Mpv.SetProperty("pause", mpv.FORMAT_FLAG, paused)
}

func (state *State) Move(time float64) {
	// TODO doesn't work, time is not absolute
	// if time == 0 && state.Percentage >= BOOST_ADD_LO {
	// 	state.WithCurrentDelta(func(delta *api.Stat) {
	// 		delta.Count++
	// 	})
	// }

	state.Mpv.Command([]string{"seek", fmt.Sprint(time)})
}