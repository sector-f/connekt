package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle()

type model struct {
	table          table.Model
	streams        []stream
	streamNames    map[string]int // Maps station names to indices
	currentStation int
	lastUpdated    time.Time
	updateCount    int
	player         *player
	winHeight      int
}

func newModel(streams []stream, streamMap map[string]int) model {
	columns := []table.Column{
		{Title: "Station", Width: 15},
		{Title: "Song", Width: 30},
		{Title: "Artist", Width: 30},
	}

	rows := []table.Row{}
	for _, s := range streams {
		rows = append(rows, table.Row{s.name, "", ""})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(streams)+1),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	t.SetStyles(s)

	return model{
		table:          t,
		streams:        streams,
		streamNames:    streamMap,
		currentStation: -1,
		player:         newPlayer(),
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.player.stop()
			return m, tea.Quit
		case "p":
			switch m.player.paused {
			case true:
				m.player.unpause()
			case false:
				m.player.pause()
			}
		case "s":
			m.player.stop()
			m.currentStation = -1
		case "enter":
			idx := m.table.Cursor()

			if idx == m.currentStation {
				return m, cmd
			}

			m.currentStation = idx
			m.player.play(m.streams[idx].streamURL)
		}
	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseWheelUp:
			m.table.MoveUp(1)
		case tea.MouseWheelDown:
			m.table.MoveDown(1)
		case tea.MouseRight:
			idx := m.table.Cursor()

			if idx == m.currentStation {
				return m, cmd
			}

			m.currentStation = idx
			m.player.play(m.streams[idx].streamURL)
		}
		return m, cmd
	case eventMessage:
		stationIdx, ok := m.streamNames[msg.Station]
		if !ok {
			return m, cmd
		}

		m.lastUpdated = msg.timestamp
		m.updateCount++

		rows := m.table.Rows()
		rows[stationIdx][1] = msg.Title
		rows[stationIdx][2] = msg.Artist
		m.table.SetRows(rows)
		m.table.UpdateViewport()
	case tea.WindowSizeMsg:
		m.winHeight = msg.Height

		flexWidth := float64(msg.Width / 5.0)

		columns := []table.Column{
			{Title: "Station", Width: int(flexWidth)},
			{Title: "Song", Width: (int(flexWidth) * 2) - 2},
			{Title: "Artist", Width: (int(flexWidth) * 2) - 2},
		}

		m.table.SetColumns(columns)
		m.table.SetWidth(msg.Width)
		m.table.UpdateViewport()
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	header := "\n c o n n e k t . f m\n\n"
	headerHeight := lipgloss.Height(header)

	tableView := m.table.View()
	tableHeight := lipgloss.Height(tableView)

	var footer string
	if m.currentStation != -1 {
		footer += "\n" + "      Now Playing: " + m.streams[m.currentStation].name
		if m.player.paused {
			footer += " (Paused)"
		}
	}
	if !m.lastUpdated.IsZero() {
		footer += fmt.Sprintf("\n Playlist Updated: %s", m.lastUpdated.Format("2006-01-02 3:04:05 PM"))
	}
	footerHeight := lipgloss.Height(footer)

	var space string
	newLineCount := m.winHeight - (headerHeight + tableHeight + footerHeight) + 2
	for i := 0; i < newLineCount; i++ {
		space += "\n"
	}

	return baseStyle.Render(header + tableView + space + footer)
}
