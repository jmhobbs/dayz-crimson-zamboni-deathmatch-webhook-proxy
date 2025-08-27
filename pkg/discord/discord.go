package discord

import (
	"bytes"
	"encoding/json"
	"net/http"

	web "github.com/jmhobbs/dayz-crimson-zamboni-deathmatch-webhook-proxy/pkg/http"
)

type notifier struct {
	url string
}

var _ web.DiscordNotifier = (*notifier)(nil)

type discordWebhookPayload struct {
	Content string `json:"content"`
}

func New(url string) *notifier {
	return &notifier{url: url}
}

func (n *notifier) PostMessage(message string) error {
	payload := discordWebhookPayload{Content: message}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.Encode(payload)

	req, err := http.NewRequest(http.MethodPost, n.url, &buf)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	http.DefaultClient.Do(req)

	return nil
}
