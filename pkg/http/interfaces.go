package http

import "github.com/jmhobbs/dayz-crimson-zamboni-deathmatch-webhook-proxy/pkg/types"

type Scoreboard interface {
	AddKill(killer, victim, weapon string, distance int)
	Reset()
	GetLongestKill() *types.Kill
	GetKDRatios() map[string]float64
	GetKills() []types.Kill
}

type DiscordNotifier interface {
	PostMessage(message string) error
}
