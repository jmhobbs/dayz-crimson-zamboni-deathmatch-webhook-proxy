package scoreboard_test

import (
	"testing"

	"github.com/jmhobbs/dayz-crimson-zamboni-deathmatch-webhook-proxy/pkg/scoreboard"
	"github.com/jmhobbs/dayz-crimson-zamboni-deathmatch-webhook-proxy/pkg/types"

	"github.com/stretchr/testify/assert"
)

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
