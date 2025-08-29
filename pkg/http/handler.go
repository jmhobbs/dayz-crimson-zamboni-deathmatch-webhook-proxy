package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/rs/zerolog/log"
)

type deathmatchPayload struct {
	Content string `json:"content"`
}

var killMatcher *regexp.Regexp = regexp.MustCompile(`^(.*) killed (.*) using (.*) from ([0-9]+)m$`)

func NewWebhookHandler(scoreboard Scoreboard, notifier DiscordNotifier, summary *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var payload deathmatchPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			log.Error().Err(err).Msg("unable to decode payload")
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		// forward the incoming message to Discord
		if err := notifier.PostMessage(payload.Content); err != nil {
			log.Error().Err(err).Str("content", payload.Content).Msg("failed to forward message to Discord")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			// intentionally continue after this error and keep feeding the scoreboard
		}

		if strings.HasPrefix(payload.Content, "**Leaderboard:**") {
			handleGameEnd(scoreboard, notifier, summary)
		} else if matches := killMatcher.FindStringSubmatch(payload.Content); len(matches) == 5 {
			handleKill(matches[1], matches[2], matches[3], matches[4], scoreboard)
		}

		w.WriteHeader(http.StatusOK)
	}
}

func handleGameEnd(scoreboard Scoreboard, notifier DiscordNotifier, summary *template.Template) {
	defer scoreboard.Reset()
	if scoreboard.GetLongestKill() == nil {
		// do nothing, no kills were recorded
		return
	}

	var buf bytes.Buffer
	if err := summary.Execute(&buf, scoreboard); err != nil {
		log.Error().Err(err).Msg("failed to execute summary template")
		return
	}

	if err := notifier.PostMessage(buf.String()); err != nil {
		log.Error().Err(err).Msg("failed to post end of game summary to Discord")
	}
}

func handleKill(killer, victim, weapon, distance string, scoreboard Scoreboard) {
	log.Debug().Msgf("%s killed %s using %s from %sm", killer, victim, weapon, distance)

	distanceInt, err := strconv.Atoi(distance)
	if err != nil {
		log.Error().Err(err).Str("distance", distance).Msg("unable to parse distance")
		return
	}

	scoreboard.AddKill(killer, victim, weapon, distanceInt)
}
