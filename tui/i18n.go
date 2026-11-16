package tui

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Nesoriel/neurolink/playerapi"
	"github.com/Nesoriel/neurolink/statusapi"
)

type Language string

const (
	LanguageEnglish           Language = "en"
	
)

func ParseLanguage(value string) (Language, error) {
	value = strings.TrimSpace(value)
	switch value {
	case "", string(LanguageEnglish):
		return LanguageEnglish, nil
		default:
		return "", fmt.Errorf("unsupported language %q (supported: en)", value)
	}
}

type lexicon struct {
	language Language

	stopped          string
	loading          string
	headerSubtitle   string
	idleMode         string
	battleMode       string
	lastUpdatePrefix string

	statusSummaryTitle string
	refreshAge         string
	historyTrend       string
	refreshQueued      string
	impactedServices   string
	downDegraded       string
	unknownChecks      string
	apiMode            string
	demoDataNotLive    string

	communityPulseTitle string
	refreshIssuePrefix  string
	noReportsLine1      string
	noReportsLine2      string
	attributionPrefix   string
	reportReceived      string

	noRegionalDetail string
	updatedPrefix    string

	footerSetup string
	footerLive  string

	playerLookupTitle    string
	playerLookupSubtitle string
	playerInputLabel     string
	playerPlaceholder    string
	playerPlatformLabel  string
	pcOriginNote         string
	playerLoading        string
	playerResultTitle    string
	playerErrorTitle     string
	playerIdleLine1      string
	playerIdleLine2      string
	playerSource         string
	playerUID            string
	playerPlatform       string
	playerLevel          string
	playerRank           string
	playerLegend         string
	playerTrackers       string

	configTitle    string
	configSource   string
	configLanguage string
	configAPIKey   string
	configEnv      string
	configDemo     string

	helpTitle     string
	helpDashboard string
	helpPlayer    string
	helpConfig    string
	helpHelp      string
	helpSearch    string
	helpPlatform  string
	helpRefresh   string
	helpQuit      string
	helpPrivacy   string

	running  string
	degraded string
	down     string
	unknown  string

	apiKeyRejected  string
	refreshCanceled string
}

func lexiconFor(language Language) lexicon {
	return lexicon{
		language: LanguageEnglish,

		stopped:          "neurolink stopped\n",
		loading:          "NEUROLINK // acquiring service feed...",
		headerSubtitle:   "Crypto surveillance drone / Apex service health",
		idleMode:         "IDLE",
		battleMode:       "BATTLE",
		lastUpdatePrefix: "Last update ",

		statusSummaryTitle: "STATUS SUMMARY",
		refreshAge:         "Refresh age",
		historyTrend:       "Trend",
		refreshQueued:      "Refresh requested; waiting for the status feed.",
		impactedServices:   "Impacted services",
		downDegraded:       "Down / degraded",
		unknownChecks:      "Unknown checks",
		apiMode:            "API mode",
		demoDataNotLive:    "Demo data is not live Apex service status.",

		communityPulseTitle: "COMMUNITY PULSE",
		refreshIssuePrefix:  "Refresh issue: ",
		noReportsLine1:      "/servers may not include a report feed.",
		noReportsLine2:      "No player profile or player-status lookup is performed.",
		attributionPrefix:   "Source: ",
		reportReceived:      "Report received",

		noRegionalDetail: "No regional detail",
		updatedPrefix:    "Updated ",

		footerSetup: "tab views  1 dashboard  2 player  3 config  ? help  / player input  enter lookup  r refresh  q quit",
		footerLive:  "tab views  1 dashboard  2 player  3 config  ? help  / player input  enter lookup  p platform  r refresh  q quit",

		playerLookupTitle:    "PLAYER LOOKUP",
		playerLookupSubtitle: "Explicit /bridge lookup by player name and platform; no background player tracking.",
		playerInputLabel:     "Player name",
		playerPlaceholder:    "Origin / PSN / Xbox name",
		playerPlatformLabel:  "Platform",
		pcOriginNote:         "PC lookup generally uses the Origin account name, even for Steam-linked accounts.",
		playerLoading:        "Looking up player...",
		playerResultTitle:    "PLAYER RESULT",
		playerErrorTitle:     "LOOKUP ERROR",
		playerIdleLine1:      "Type a player name, choose a platform, then press Enter.",
		playerIdleLine2:      "Supports PC, PS4, and X1; lookup only runs after your visible action.",
		playerSource:         "Source",
		playerUID:            "UID",
		playerPlatform:       "Platform",
		playerLevel:          "Level",
		playerRank:           "Rank",
		playerLegend:         "Selected legend",
		playerTrackers:       "Trackers",

		configTitle:    "CONFIG",
		configSource:   "Current status source",
		configLanguage: "Language",
		configAPIKey:   "API key: use --api-key, NEUROLINK_APEX_API_KEY, or neurolink config set api-key.",
		configEnv:      "Language: --lang or NEUROLINK_LANG; supported value: en.",
		configDemo:     "Without an API key, clearly labeled demo data is used; player lookup also returns demo results.",

		helpTitle:     "HELP",
		helpDashboard: "1 / tab: dashboard",
		helpPlayer:    "2 / /: player lookup view and focused input",
		helpConfig:    "3: config view",
		helpHelp:      "?: help view",
		helpSearch:    "enter: search from the player view",
		helpPlatform:  "p: cycle PC / PS4 / X1 when the player input is not focused",
		helpRefresh:   "r / F5: request an immediate service refresh",
		helpQuit:      "Esc: dashboard from other views; q / ctrl+c: quit",
		helpPrivacy:   "Player data is requested only after explicit lookup; neurolink does not continuously track players or hide telemetry.",

		running:  "RUNNING",
		degraded: "DEGRADED",
		down:     "DOWN",
		unknown:  "UNKNOWN",

		apiKeyRejected:  "API key rejected by status API",
		refreshCanceled: "refresh canceled",
	}
}

func (c lexicon) statusLabel(status statusapi.Status) string {
	switch status {
	case statusapi.StatusHealthy:
		return "● " + c.running
	case statusapi.StatusDegraded:
		return "▲ " + c.degraded
	case statusapi.StatusDown:
		return "✕ " + c.down
	default:
		return "? " + c.unknown
	}
}

func (c lexicon) viewLabel(view appView) string {
	switch view {
	case viewPlayer:
		return "PLAYER"
	case viewConfig:
		return "CONFIG"
	case viewHelp:
		return "HELP"
	default:
		return "DASHBOARD"
	}
}

func (c lexicon) viewKeyLabel(view appView) string {
	switch view {
	case viewPlayer:
		return "2 " + c.viewLabel(view)
	case viewConfig:
		return "3 " + c.viewLabel(view)
	case viewHelp:
		return "? " + c.viewLabel(view)
	default:
		return "1 " + c.viewLabel(view)
	}
}

func (c lexicon) footerHint(view appView, source statusapi.SourceMode) string {
	if source == statusapi.SourceLive {
		return c.footerLive
	}
	return c.footerSetup
}

func (c lexicon) modeLabel(battle bool) string {
	if battle {
		return "● " + c.battleMode
	}
	return "● " + c.idleMode
}

func (c lexicon) notice(value string) string {
	return strings.TrimSpace(value)
}

func (c lexicon) attributionLine(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	return c.attributionPrefix + value
}

func (c lexicon) playerError(err error) string {
	switch {
	case errors.Is(err, playerapi.ErrAPIKeyRequired):
		return "API key is not configured. Use --api-key, NEUROLINK_APEX_API_KEY, or config set api-key."
	case errors.Is(err, playerapi.ErrInvalidAPIKey):
		return "Player lookup API rejected the API key."
	case errors.Is(err, playerapi.ErrPlayerNotFound):
		return "Player was not found. Check the player name and platform."
	case errors.Is(err, playerapi.ErrUnknownPlatform):
		return "Unknown platform. Supported platforms are PC, PS4, and X1."
	case errors.Is(err, playerapi.ErrRateLimited):
		return "Player lookup is rate limited. Try again later."
	case errors.Is(err, playerapi.ErrPlayerRequired):
		return "Enter a player name before searching."
	case errors.Is(err, playerapi.ErrInternalAPI):
		return "Player lookup API returned an internal error."
	case errors.Is(err, playerapi.ErrExternalAPI):
		return "Upstream player lookup is temporarily unavailable."
	default:
		return "Player lookup failed."
	}
}

func (c lexicon) summarizeError(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if strings.Contains(value, "HTTP 401") || strings.Contains(value, "HTTP 403") {
		return c.apiKeyRejected
	}
	if strings.Contains(value, "context canceled") {
		return c.refreshCanceled
	}
	return fit(value, 72)
}
