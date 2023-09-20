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
	h := mpv.Create()
	mpv.Initialize(h)

	return &player{
		m:         h,
		paused:    false,
		isRunning: false,
	}
}

func (p *player) play(url string) {
	mpv.Command(p.m, "loadfile "+url)
	p.isRunning = true
	p.unpause()
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
