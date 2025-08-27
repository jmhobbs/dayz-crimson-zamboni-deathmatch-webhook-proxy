package http_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	web "github.com/jmhobbs/dayz-crimson-zamboni-deathmatch-webhook-proxy/pkg/http"
	"github.com/jmhobbs/dayz-crimson-zamboni-deathmatch-webhook-proxy/pkg/types"
	"github.com/stretchr/testify/assert"
)

func Test_WebhookHandler(t *testing.T) {
	t.Run("rejects non-POST methods", func(t *testing.T) {
		scoreboard := NewMockScoreboard(t)
		notifier := NewMockDiscordNotifier(t)
		handler := web.NewWebhookHandler(scoreboard, notifier)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/webhook", nil)
		assert.NoError(t, err)

		handler(w, req)
		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})

	t.Run("requires valid JSON", func(t *testing.T) {
		scoreboard := NewMockScoreboard(t)
		notifier := NewMockDiscordNotifier(t)
		handler := web.NewWebhookHandler(scoreboard, notifier)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/webhook", bytes.NewBufferString("not json"))
		assert.NoError(t, err)

		handler(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("forwards the request to Discord", func(t *testing.T) {
		scoreboard := NewMockScoreboard(t)
		notifier := NewMockDiscordNotifier(t)
		handler := web.NewWebhookHandler(scoreboard, notifier)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/webhook", bytes.NewBufferString(`{"content":"some message here"}`))
		assert.NoError(t, err)

		notifier.EXPECT().PostMessage("some message here").Return(nil)
		handler(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("posts no summary if no kills were made", func(t *testing.T) {
		scoreboard := NewMockScoreboard(t)
		notifier := NewMockDiscordNotifier(t)
		handler := web.NewWebhookHandler(scoreboard, notifier)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/webhook", bytes.NewBufferString(`{"content":"**Leaderboard:**..."}`))
		assert.NoError(t, err)

		scoreboard.EXPECT().GetLongestKill().Return(nil)
		scoreboard.EXPECT().GetKDRatios().Return(map[string]float64{})
		scoreboard.EXPECT().Reset()
		notifier.EXPECT().PostMessage("**Leaderboard:**...").Return(nil)
		handler(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("posts a summary with the longest kill and K/D ratios to Discord", func(t *testing.T) {
		scoreboard := NewMockScoreboard(t)
		notifier := NewMockDiscordNotifier(t)
		handler := web.NewWebhookHandler(scoreboard, notifier)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/webhook", bytes.NewBufferString(`{"content":"**Leaderboard:**..."}`))
		assert.NoError(t, err)

		scoreboard.EXPECT().GetLongestKill().Return(&types.Kill{Killer: "jmhobbs", Victim: "Reader", Weapon: "Screwdriver", Distance: 420})
		scoreboard.EXPECT().GetKDRatios().Return(map[string]float64{"jmhobbs": 1.0, "Reader": 0.0})
		scoreboard.EXPECT().Reset()
		notifier.EXPECT().PostMessage("**Leaderboard:**...").Return(nil)
		notifier.EXPECT().PostMessage("Longest kill: jmhobbs killed Reader with Screwdriver at 420m").Return(nil)
		notifier.EXPECT().PostMessage("K/D Ratios:\n```jmhobbs: 1.0\nReader: 0.0\n```").Return(nil)
		handler(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("records kills sent in the payload", func(t *testing.T) {
		scoreboard := NewMockScoreboard(t)
		notifier := NewMockDiscordNotifier(t)
		handler := web.NewWebhookHandler(scoreboard, notifier)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/webhook", bytes.NewBufferString(`{"content":"jmhobbs killed Reader using Screwdriver from 420m"}`))
		assert.NoError(t, err)

		scoreboard.EXPECT().AddKill("jmhobbs", "Reader", "Screwdriver", 420)
		notifier.EXPECT().PostMessage("jmhobbs killed Reader using Screwdriver from 420m").Return(nil)
		handler(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
