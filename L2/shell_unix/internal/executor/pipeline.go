package executor

import (
	"io"
	"sync"

	"github.com/merkulovlad/wildberries-L2/shell_unix/internal/parser"
	"github.com/merkulovlad/wildberries-L2/shell_unix/internal/process"
)

func (e *Executor) executePipeline(pipe *parser.Pipeline) error {
	if pipe == nil || len(pipe.Commands) == 0 {
		return nil
	}
	if len(pipe.Commands) == 1 {
		return e.Execute(pipe.Commands[0])
	}

	readers := make([]*io.PipeReader, len(pipe.Commands)-1)
	writers := make([]*io.PipeWriter, len(pipe.Commands)-1)
	for i := range readers {
		readers[i], writers[i] = io.Pipe()
	}

	group := process.NewGroup()
	e.setForeground(group)
	defer e.clearForeground()

	errCh := make(chan error, len(pipe.Commands))
	var wg sync.WaitGroup

	for i, node := range pipe.Commands {
		stdin := e.stdin
		if i > 0 {
			stdin = readers[i-1]
		}

		stdout := e.stdout
		if i < len(pipe.Commands)-1 {
			stdout = writers[i]
		}

		child := e.withIO(stdin, stdout, e.stderr)

		var stdinCloser io.Closer
		if i > 0 {
			stdinCloser = readers[i-1]
		}

		var stdoutCloser io.Closer
		if i < len(pipe.Commands)-1 {
			stdoutCloser = writers[i]
		}

		wg.Add(1)
		go func(node parser.Node, child *Executor, stdinCloser io.Closer, stdoutCloser io.Closer) {
			defer wg.Done()

			err := child.Execute(node)

			if stdoutCloser != nil {
				_ = stdoutCloser.Close()
			}
			if stdinCloser != nil {
				_ = stdinCloser.Close()
			}

			errCh <- err
		}(node, child, stdinCloser, stdoutCloser)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			return err
		}
	}

	return nil
}
