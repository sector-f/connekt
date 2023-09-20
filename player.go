package main

import (
	mpv "lncn.dev/x/libmpv"
)

type player struct {
	m *mpv.Handle
}

func newPlayer() *player {
	h := mpv.Create()
	mpv.Initialize(h)
	return &player{h}
}

func (p *player) play(url string) error {
	mpv.Command(p.m, "loadfile "+url)

	return nil
}

func (p *player) stop() {
	mpv.Command(p.m, "quit")
}
