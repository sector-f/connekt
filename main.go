package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/r3labs/sse/v2"
)

type stream struct {
	name      string
	streamURL string
}

func main() {
	streams := []stream{
		{
			name:      "REKT",
			streamURL: "https://stream.nightride.fm:8443/rekt/aac_hifi.m3u8",
		},
		{
			name:      "REKTORY",
			streamURL: "https://stream.nightride.fm:8443/rektory/aac_hifi.m3u8",
		},
		{
			name:      "NIGHTRIDE FM",
			streamURL: "https://stream.nightride.fm:8443/nightride/aac_hifi.m3u8",
		},
		{
			name:      "CHILLSYNTH FM",
			streamURL: "https://stream.nightride.fm:8443/chillsynth/aac_hifi.m3u8",
		},
		{
			name:      "DATAWAVE FM",
			streamURL: "https://stream.nightride.fm:8443/datawave/aac_hifi.m3u8",
		},
		{
			name:      "SPACESYNTH FM",
			streamURL: "https://stream.nightride.fm:8443/spacesynth/aac_hifi.m3u8",
		},
		{
			name:      "DARKSYNTH",
			streamURL: "https://stream.nightride.fm:8443/darksynth/aac_hifi.m3u8",
		},
		{
			name:      "HORRORSYNTH",
			streamURL: "https://stream.nightride.fm:8443/horrorsynth/aac_hifi.m3u8",
		},
		{
			name:      "EBSM",
			streamURL: "https://stream.nightride.fm:8443/ebsm/aac_hifi.m3u8",
		},
	}

	streamMap := map[string]int{
		"rekt":        0,
		"rektory":     1,
		"nightride":   2,
		"chillsynth":  3,
		"datawave":    4,
		"spacesynth":  5,
		"darksynth":   6,
		"horrorsynth": 7,
		"ebsm":        8,
	}

	m := newModel(streams, streamMap)
	program := tea.NewProgram(m, tea.WithAltScreen())

	go func() {
		events := make(chan *sse.Event)
		eventClient := sse.NewClient("https://rekt.network/meta")
		eventClient.SubscribeChanRaw(events)

		refreshTime := 1 * time.Minute // The server should send a keepalive every minute

		ticker := time.NewTicker(refreshTime)

		for {
			select {
			case <-ticker.C:
				// The events channel occasionally stops receiving events for some reason.
				// I can't figure out why.
				// So, just restart it if we go 5 minutes without receiving anything.
				eventClient.Unsubscribe(events)
				eventClient.SubscribeChanRaw(events)
			case event := <-events:
				ticker.Reset(refreshTime)

				rcvdAt := time.Now()

				if event.Data == nil {
					continue
				}

				var es []eventMessage
				err := json.Unmarshal(event.Data, &es)
				if err != nil {
					continue
				}

				for _, e := range es {
					e.timestamp = rcvdAt
					program.Send(e)
				}
			}
		}
	}()

	if _, err := program.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
