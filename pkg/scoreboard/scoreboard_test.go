package scoreboard_test

import (
	"bytes"
	"testing"

	"github.com/jmhobbs/dayz-crimson-zamboni-deathmatch-webhook-proxy/pkg/scoreboard"
	"github.com/jmhobbs/dayz-crimson-zamboni-deathmatch-webhook-proxy/pkg/types"

	"github.com/stretchr/testify/assert"
)

func Test_NewFromJSON(t *testing.T) {
	t.Run("returns an error if not provided valid JSON", func(t *testing.T) {
		_, err := scoreboard.NewFromJSON(bytes.NewBufferString("not json"))
		assert.Error(t, err)
	})

	t.Run("decodes kills out of JSON", func(t *testing.T) {
		sb, err := scoreboard.NewFromJSON(bytes.NewBufferString(`{"kills":[{"killer":"Player1","victim":"Player2","weapon":"M70 Tundra","distance":50}]}`))
		assert.Nil(t, err)
		assert.Equal(t, 1, len(sb.GetKills()))
	})
}

func Test_Scoreboard_GetLongestKill(t *testing.T) {
	t.Run("returns the longest kill", func(t *testing.T) {
		expectedLongest := types.Kill{Killer: "Player3", Victim: "Player4", Weapon: "Mlock", Distance: 100}

		s := scoreboard.New()
		s.AddKill("Player1", "Player2", "M70 Tundra", 50)
		s.AddKill(expectedLongest.Killer, expectedLongest.Victim, expectedLongest.Weapon, expectedLongest.Distance)
		s.AddKill("Player5", "Player6", "M79", 30)

		longest := s.GetLongestKill()
		assert.NotNil(t, longest)

		assert.Equal(t, expectedLongest, *longest)
	})

	t.Run("returns the first kill when there is a tie", func(t *testing.T) {
		expectedLongest := types.Kill{Killer: "Player3", Victim: "Player4", Weapon: "Mlock", Distance: 100}

		s := scoreboard.New()
		s.AddKill("Player1", "Player2", "M70 Tundra", 50)
		s.AddKill(expectedLongest.Killer, expectedLongest.Victim, expectedLongest.Weapon, expectedLongest.Distance)
		s.AddKill("Player5", "Player6", "M79", 30)
		s.AddKill("Player5", "Player6", "M79", 100)

		longest := s.GetLongestKill()
		assert.NotNil(t, longest)

		assert.Equal(t, expectedLongest, *longest)
	})
}

func Test_Scoreboard_GetKDRatios(t *testing.T) {
	s := scoreboard.New()
	s.AddKill("Player 1", "Player 2", "M70 Tundra", 50)
	s.AddKill("Player 1", "Player 2", "M70 Tundra", 50)
	s.AddKill("Player 2", "Player 1", "M79", 30)
	s.AddKill("Player 3", "Player 4", "Screwdriver", 30)
	s.AddKill("Player 3", "Player 1", "Screwdriver", 300)
	s.AddKill("Player 3", "Player 4", "Screwdriver", 150)

	ratios := s.GetKDRatios()

	assert.Equal(t, map[string]float64{
		"Player 1": 1.0, // 2:1
		"Player 2": 0.5, // 1:2
		"Player 3": 3.0, // 3:0
		"Player 4": 0.0, // 0:2
	}, ratios)
}

func Test_Scoreboard_Reset(t *testing.T) {
	s := scoreboard.New()
	s.AddKill("Player1", "Player2", "M70 Tundra", 50)
	assert.NotNil(t, s.GetLongestKill())
	s.Reset()
	assert.Nil(t, s.GetLongestKill())
}
