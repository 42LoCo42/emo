package state

type State struct {
	CurrentSong string
	Queue       []string
	Progress    uint
	Paused      bool
}

func (state *State) SelectSong(song string) {
	state.CurrentSong = song
}

func (state *State) SetPaused(paused bool) {
	state.Paused = paused
}

func (state *State) Move(time float64) {

}
