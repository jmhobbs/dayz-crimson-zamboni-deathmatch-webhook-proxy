package discord_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jmhobbs/dayz-crimson-zamboni-deathmatch-webhook-proxy/pkg/discord"

	"github.com/stretchr/testify/assert"
)

func Test_DiscordNotifier_PostMessage(t *testing.T) {
	var messages [][]byte

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		message, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		messages = append(messages, message)
		defer func() {
			assert.NoError(t, r.Body.Close())
		}()
		w.WriteHeader(http.StatusOK)
	}))

	assert.NoError(t, discord.New(srv.URL).PostMessage("test message"))
	assert.Equal(t, 1, len(messages))
	assert.Equal(t, `{"content":"test message"}`+"\n", string(messages[0]))
}
