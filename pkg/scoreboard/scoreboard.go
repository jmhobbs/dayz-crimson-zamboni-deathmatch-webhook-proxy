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
