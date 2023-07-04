package util

import (
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/42LoCo42/emo/shared"
)

func RunPipeline(
	i io.Reader,
	o io.Writer,
	cmds ...[]string,
) error {
	rPipes := make([]io.Reader, len(cmds))
	wPipes := make([]io.Writer, len(cmds))
	processes := make([]*exec.Cmd, len(cmds))

	for i := 1; i < len(rPipes); i++ {
		rPipes[i], wPipes[i-1] = io.Pipe()
	}

	rPipes[0] = i
	wPipes[len(wPipes)-1] = o

	wg := &sync.WaitGroup{}

	run := func(i int, cmd []string) {
		proc := exec.Command(cmd[0], cmd[1:]...)
		proc.Stdin = rPipes[i]
		proc.Stdout = wPipes[i]
		proc.ExtraFiles = append(proc.ExtraFiles, os.Stdout)
		processes[i] = proc

		if err := proc.Start(); err != nil {
			shared.Wrap(err, "could not start command %v", cmd)
		}

		wg.Add(1)
		go func() {
			proc.Wait()
			rPipes[i].(io.ReadCloser).Close()
			wPipes[i].(io.WriteCloser).Close()
			wg.Done()
		}()
	}

	for i, cmd := range cmds {
		run(i, cmd)
	}

	wg.Wait()
	return nil
}
