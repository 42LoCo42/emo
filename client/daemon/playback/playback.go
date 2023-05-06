package playback

import (
	"fmt"
	"strconv"

	"github.com/42LoCo42/emo/client/daemon/state"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func Cmd(state *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use: "playback",
	}

	cmd.AddCommand(
		show(state),
		next(state),
		play(state),
		pause(state),
		toggle(state),
		forward(state),
		backward(state),
	)

	return cmd
}

func show(state *state.State) *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show the currently playing song",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(cmd.OutOrStdout(), state.CurrentSong)
		},
	}
}

func next(state *state.State) *cobra.Command {
	return &cobra.Command{
		Use:   "next",
		Short: "Move to the next song (either from queue or randomly selected)",
		Run: func(cmd *cobra.Command, args []string) {
			var song string

			if len(state.Queue) > 0 {
				song = state.Queue[0]
				state.Queue = state.Queue[1:]
			} else {
				panic("TODO randomselect song")
			}

			state.SelectSong(song)
			fmt.Fprintln(cmd.OutOrStdout(), song)
		},
	}
}

func play(state *state.State) *cobra.Command {
	return &cobra.Command{
		Use:   "play",
		Short: "Start playback",
		Run: func(cmd *cobra.Command, args []string) {
			state.SetPaused(false)
		},
	}
}

func pause(state *state.State) *cobra.Command {
	return &cobra.Command{
		Use:   "pause",
		Short: "Stop playback",
		Run: func(cmd *cobra.Command, args []string) {
			state.SetPaused(true)
		},
	}
}

func toggle(state *state.State) *cobra.Command {
	return &cobra.Command{
		Use:   "toggle",
		Short: "Toggle playback",
		Run: func(cmd *cobra.Command, args []string) {
			state.SetPaused(!state.Paused)
			fmt.Fprintln(cmd.OutOrStdout(), "paused:", state.Paused)
		},
	}
}

func forward(state *state.State) *cobra.Command {
	return &cobra.Command{
		Use:   "forward seconds",
		Short: "Move current song forwards",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			time, err := strconv.ParseFloat(args[0], 64)
			if err != nil {
				fmt.Fprintln(
					cmd.ErrOrStderr(),
					errors.Wrap(err, "could not parse time"),
				)
			}

			state.Move(time)
		},
	}
}

func backward(state *state.State) *cobra.Command {
	return &cobra.Command{
		Use:   "backward seconds",
		Short: "Move current song backwards",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			time, err := strconv.ParseFloat(args[0], 64)
			if err != nil {
				fmt.Fprintln(
					cmd.ErrOrStderr(),
					errors.Wrap(err, "could not parse time"),
				)
			}

			state.Move(-time)
		},
	}
}
