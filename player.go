package main

import (
	mpv "lncn.dev/x/libmpv"
)

type player struct {
	m         *mpv.Handle
	paused    bool
	isRunning bool
}

func newPlayer() *player {
	return &player{
		paused:    false,
		isRunning: false,
	}
}

func (p *player) play(url string) error {
	if !p.isRunning {
		h := mpv.Create()
		mpv.Initialize(h)
		p.m = h
	}

	mpv.Command(p.m, "loadfile "+url)
	p.unpause()
	p.isRunning = true

	return nil
}

func (p *player) pause() {
	if !p.isRunning {
		return
	}

	mpv.SetProperty(p.m, "pause", "yes")
	p.paused = true
}

func (p *player) unpause() {
	if !p.isRunning {
		return
	}

	mpv.SetProperty(p.m, "pause", "no")
	p.paused = false
}

func (p *player) stop() {
	mpv.Command(p.m, "stop")
	p.isRunning = false
	p.paused = false
}
