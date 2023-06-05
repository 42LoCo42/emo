package queue

import (
	"fmt"
	"strconv"

	"github.com/42LoCo42/emo/daemon/state"
	"github.com/42LoCo42/emo/shared"
	"github.com/spf13/cobra"
)

func Cmd(state *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use: "queue",
	}

	cmd.AddCommand(
		show(state),
		push(state),
		del(state),
	)

	return cmd
}

func show(state *state.State) *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show the current queue",
		Run: func(cmd *cobra.Command, args []string) {
			for _, song := range state.Queue {
				fmt.Fprintln(cmd.OutOrStdout(), song)
			}
		},
	}
}

func push(state *state.State) *cobra.Command {
	return &cobra.Command{
		Use:   "push song...",
		Short: "Push one or multiple song(s) onto the queue",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			for _, song := range args {
				state.Queue = append(state.Queue, song)
			}
		},
	}
}

func del(state *state.State) *cobra.Command {
	return &cobra.Command{
		Use:   "delete index",
		Short: "Delete a queue entry at an index",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			index, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				fmt.Fprintln(
					cmd.ErrOrStderr(),
					shared.Wrap(err, "could not parse index"),
				)
				return
			}

			if int(index) >= len(state.Queue) {
				fmt.Fprintf(
					cmd.ErrOrStderr(),
					"Index out of bounds (max. %d)\n",
					len(state.Queue)-1,
				)
				return
			}

			state.Queue = append(
				state.Queue[:index],
				state.Queue[index+1:]...,
			)
		},
	}
}
