package main

import (
	mpv "lncn.dev/x/libmpv"
)

type player struct {
	m      *mpv.Handle
	paused bool
}

func newPlayer() *player {
	h := mpv.Create()
	mpv.Initialize(h)
	return &player{h, false}
}

func (p *player) play(url string) error {
	mpv.Command(p.m, "loadfile "+url)

	return nil
}

func (p *player) pause() {
	mpv.SetProperty(p.m, "pause", "yes")
	p.paused = true
}

func (p *player) unpause() {
	mpv.SetProperty(p.m, "pause", "no")
	p.paused = false
}

func (p *player) stop() {
	mpv.Command(p.m, "quit")
}
