package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

type deathmatchPayload struct {
	Content string `json:"content"`
}

var killMatcher *regexp.Regexp = regexp.MustCompile(`^(.*) killed (.*) using (.*) from ([0-9]+)m$`)

func NewWebhookHandler(scoreboard Scoreboard, notifier DiscordNotifier) http.HandlerFunc {
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
			// intentionally ignore this error and keep feeding the scoreboard
		}

		// Game end, post summary
		if strings.HasPrefix(payload.Content, "**Leaderboard:**") {
			log.Debug().Msg("End of game detected, posting summary...")
			longest := scoreboard.GetLongestKill()
			if longest != nil {
				summary := fmt.Sprintf("Longest kill: %s killed %s with %s at %dm", longest.Killer, longest.Victim, longest.Weapon, longest.Distance)
				log.Info().Msg(summary)
				if err := notifier.PostMessage(summary); err != nil {
					log.Error().Err(err).Msg("failed to post end of game summary to Discord")
				}
			} else {
				log.Info().Msg("No kills recorded.")
			}
			scoreboard.Reset()
			w.WriteHeader(http.StatusOK)
			return
		}

		// Is it a kill?
		matches := killMatcher.FindStringSubmatch(payload.Content)
		if len(matches) == 5 {
			killer := matches[1]
			victim := matches[2]
			weapon := matches[3]
			distance := matches[4]

			log.Debug().Msgf("%s killed %s using %s from %sm", killer, victim, weapon, distance)

			distanceInt, err := strconv.Atoi(distance)
			if err != nil {
				log.Error().Err(err).Str("distance", distance).Msg("unable to parse distance")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			// Add to scoreboard
			scoreboard.AddKill(killer, victim, weapon, distanceInt)
		}

		w.WriteHeader(http.StatusOK)
	}
}
