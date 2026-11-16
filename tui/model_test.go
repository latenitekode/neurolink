package tui

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Nesoriel/neurolink/collector"
	"github.com/Nesoriel/neurolink/playerapi"
	"github.com/Nesoriel/neurolink/statusapi"
)

func TestViewRendersPolishedDemoDashboard(t *testing.T) {
	now := time.Unix(100, 0)
	snapshot := collector.Snapshot{
		Mode:      collector.ModeIdle,
		Status:    demoSnapshotForView(now),
		UpdatedAt: now,
		LastError: "",
	}

	model := NewModel(make(chan collector.Snapshot))
	updated, _ := model.Update(tea.WindowSizeMsg{Width: 110, Height: 36})
	model = updated.(Model)
	updated, _ = model.Update(snapshotMsg(snapshot))
	model = updated.(Model)

	view := model.View()
	for _, want := range []string{
		"NEUROLINK",
		"STATUS SUMMARY",
		"COMMUNITY PULSE",
		"Refresh age",
		"Trend",
		"Lobby / Matchmaking",
		"Crossplay Auth",
		"DEMO",
		"Demo data is not live Apex service status.",
		"JP",
		"▰",
		"q quit",
	} {
		if !strings.Contains(view, want) {
			t.Fatalf("view does not contain %q\n%s", want, view)
		}
	}
}

func TestPlayerLookupViewRunsAsyncLookup(t *testing.T) {
	now := time.Unix(100, 0)
	snapshot := collector.Snapshot{
		Mode:      collector.ModeIdle,
		Status:    demoSnapshotForView(now),
		UpdatedAt: now,
	}
	provider := &fakePlayerProvider{snapshot: playerapi.PlayerSnapshot{
		Source:         playerapi.SourceLive,
		Attribution:    playerapi.Attribution,
		LookupAt:       now,
		PlayerName:     "Ace Origin",
		Platform:       playerapi.PlatformPC,
		UID:            "uid-1",
		Level:          99,
		HasLevel:       true,
		RankName:       "Gold",
		RankDivision:   "1",
		RankScore:      7600,
		HasRank:        true,
		SelectedLegend: "Crypto",
		Trackers:       []playerapi.Tracker{{Name: "Kills", Value: "123"}},
	}}

	model := NewModelWithPlayerProvider(make(chan collector.Snapshot), provider, nil)
	updated, _ := model.Update(tea.WindowSizeMsg{Width: 110, Height: 36})
	model = updated.(Model)
	updated, _ = model.Update(snapshotMsg(snapshot))
	model = updated.(Model)
	updated, _ = model.Update(keyMsg("/"))
	model = updated.(Model)
	for _, r := range []rune("Ace Origin") {
		updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		model = updated.(Model)
	}

	updated, cmd := model.Update(keyMsg("enter"))
	model = updated.(Model)
	if cmd == nil {
		t.Fatal("enter should return async lookup command")
	}
	msg := cmd()
	updated, _ = model.Update(msg)
	model = updated.(Model)

	if provider.last.Player != "Ace Origin" || provider.last.Platform != playerapi.PlatformPC {
		t.Fatalf("lookup request = %#v, want typed player on PC", provider.last)
	}
	view := model.View()
	for _, want := range []string{
		"PLAYER LOOKUP",
		"Ace Origin",
		"uid-1",
		"Gold 1 / 7600",
		"Crypto",
		"Kills",
	} {
		if !strings.Contains(view, want) {
			t.Fatalf("view does not contain %q\n%s", want, view)
		}
	}
}

func TestNavigationHelpConfigAndRefreshKey(t *testing.T) {
	now := time.Unix(100, 0)
	snapshot := collector.Snapshot{
		Mode:      collector.ModeIdle,
		Status:    demoSnapshotForView(now),
		UpdatedAt: now,
	}
	refresh := make(chan struct{}, 1)
	model := NewModelWithPlayerProvider(make(chan collector.Snapshot), &fakePlayerProvider{}, refresh)
	updated, _ := model.Update(tea.WindowSizeMsg{Width: 140, Height: 34})
	model = updated.(Model)
	updated, _ = model.Update(snapshotMsg(snapshot))
	model = updated.(Model)

	updated, cmd := model.Update(keyMsg("r"))
	model = updated.(Model)
	if cmd == nil {
		t.Fatal("r should return refresh command")
	}
	updated, _ = model.Update(cmd())
	model = updated.(Model)
	select {
	case <-refresh:
	default:
		t.Fatal("refresh channel did not receive request")
	}
	if !strings.Contains(model.View(), "Refresh requested") {
		t.Fatalf("dashboard did not show refresh request\n%s", model.View())
	}

	updated, _ = model.Update(keyMsg("3"))
	model = updated.(Model)
	if view := model.View(); !strings.Contains(view, "CONFIG") || !strings.Contains(view, "NEUROLINK_APEX_API_KEY") {
		t.Fatalf("config view missing expected content\n%s", view)
	}
	updated, _ = model.Update(keyMsg("/"))
	model = updated.(Model)
	updated, _ = model.Update(keyMsg("?"))
	model = updated.(Model)
	if view := model.View(); !strings.Contains(view, "HELP") || !strings.Contains(view, "does not continuously track players") {
		t.Fatalf("help view missing expected content\n%s", view)
	}
}

