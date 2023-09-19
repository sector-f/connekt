package main

import (
	"context"
	"io"
	"os/exec"
)

type player struct {
	stdin   io.WriteCloser
	stopper func()
}

func (p *player) play(url string) error {
	ctx, stopper := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, "mpv", url)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		stopper()
		return err
	}

	p.stdin = stdin
	p.stopper = stopper

	go cmd.Run()

	return nil
}

func (p *player) stop() {
	if p.stopper != nil {
		p.stopper()
	}
}
