package scoreboard

import (
	"github.com/jmhobbs/dayz-crimson-zamboni-deathmatch-webhook-proxy/pkg/http"
	"github.com/jmhobbs/dayz-crimson-zamboni-deathmatch-webhook-proxy/pkg/types"
)

var _ http.Scoreboard = (*scoreboard)(nil)

type scoreboard struct {
	kills []types.Kill
}

func New() *scoreboard {
	return &scoreboard{kills: []types.Kill{}}
}

func (s *scoreboard) Reset() {
	s.kills = []types.Kill{}
}

func (s *scoreboard) AddKill(killer, victim, weapon string, distance int) {
	s.kills = append(s.kills, types.Kill{
		Killer:   killer,
		Victim:   victim,
		Weapon:   weapon,
		Distance: distance,
	})
}

func (s *scoreboard) GetLongestKill() *types.Kill {
	var longest *types.Kill
	for _, kill := range s.kills {
		if longest == nil || kill.Distance > longest.Distance {
			longest = &kill
		}
	}
	return longest
}

func (s *scoreboard) GetKDRatios() map[string]float64 {
	kills := make(map[string]float64)
	deaths := make(map[string]float64)
	for _, kill := range s.kills {
		kills[kill.Killer]++
		deaths[kill.Victim]++
	}
	ratios := make(map[string]float64)
	for name, count := range kills {
		deathCount := deaths[name]
		if deathCount == 0 {
			ratios[name] = count
		} else {
			ratios[name] = count / deathCount
		}
	}
	for name := range deaths {
		if _, ok := kills[name]; !ok {
			ratios[name] = 0.0
		}
	}
	return ratios
}
