package playback

import (
	"fmt"
	"strconv"

	"github.com/42LoCo42/emo/daemon/state"
	"github.com/42LoCo42/emo/daemon/util"
	"github.com/42LoCo42/emo/shared"
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
		seek(state, "forward", 1),
		seek(state, "backward", -1),
		rewind(state),
		random10k(state),
	)

	return cmd
}

func show(state *state.State) *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show the currently playing song",
		Run: func(cmd *cobra.Command, args []string) {
			if state.CurrentStat == nil {
				fmt.Fprintln(cmd.ErrOrStderr(), "error: no song is playing!")
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), "song:   ", state.CurrentStat.Song)
				fmt.Fprintln(cmd.OutOrStdout(), "time:   ", state.Time)
				fmt.Fprintln(cmd.OutOrStdout(), "percent:", state.Percentage)
				fmt.Fprintln(cmd.OutOrStdout(), "paused: ", state.Paused)
			}
		},
	}
}

func next(state *state.State) *cobra.Command {
	return &cobra.Command{
		Use:   "next",
		Short: "Move to the next song (either from queue or randomly selected)",
		Run: func(cmd *cobra.Command, args []string) {
			song, err := state.NextSong()
			if err != nil {
				fmt.Fprintln(cmd.ErrOrStderr(), "error:", err)
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), "song:", song)
			}
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

func seek(
	state *state.State,
	name string,
	factor float64,
) *cobra.Command {
	return &cobra.Command{
		Use:   name + " seconds",
		Short: "Seek " + name + " in current song",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			time, err := strconv.ParseFloat(args[0], 64)
			if err != nil {
				fmt.Fprintln(
					cmd.ErrOrStderr(),
					shared.Wrap(err, "could not parse time"),
				)
			}

			state.Move(time * factor)
			fmt.Fprintln(cmd.OutOrStdout(), "time:", state.Time)
			fmt.Fprintln(cmd.OutOrStdout(), "percent:", state.Percentage)
		},
	}
}

func rewind(state *state.State) *cobra.Command {
	return &cobra.Command{
		Use:   "rewind",
		Short: "Rewinds the current song",
		Run: func(cmd *cobra.Command, args []string) {
			if err := state.Rewind(); err != nil {
				fmt.Fprintln(cmd.ErrOrStderr(), shared.Wrap(err, "rewind failed"))
			}
		},
	}
}

func random10k(state *state.State) *cobra.Command {
	return &cobra.Command{
		Use:   "random10k",
		Short: "Debug command. Prints 10000 randomly selected songs.",
		Run: func(cmd *cobra.Command, args []string) {
			for i := 0; i < 10000; i++ {
				fmt.Fprintln(
					cmd.OutOrStdout(),
					util.RandomStat(&state.Stats, &state.Deltas).Song,
				)
			}
		},
	}
}
