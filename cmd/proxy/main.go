package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/jmhobbs/dayz-crimson-zamboni-deathmatch-webhook-proxy/pkg/discord"
	web "github.com/jmhobbs/dayz-crimson-zamboni-deathmatch-webhook-proxy/pkg/http"
	"github.com/jmhobbs/dayz-crimson-zamboni-deathmatch-webhook-proxy/pkg/scoreboard"

	"github.com/peterbourgon/ff/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

	fs := flag.NewFlagSet("dayz-crimson-zamboni-deathmatch-webhook-proxy", flag.ContinueOnError)
	var (
		listenAddr = fs.String("listen", ":8080", "listen address")
		webhookUrl = fs.String("webhook-url", "", "Discord webhook URL")
		debug      = fs.Bool("debug", false, "log debug information")
		_          = fs.String("config", "", "config file (optional)")
	)

	err := ff.Parse(fs, os.Args[1:],
		ff.WithEnvVarNoPrefix(),
		ff.WithConfigFileFlag("config"),
		ff.WithConfigFileParser(ff.PlainParser),
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse flags or config")
		os.Exit(1)
	}

	if *webhookUrl == "" {
		log.Error().Msg("webhook-url is required")
		os.Exit(1)
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	sb := scoreboard.New()
	notifier := discord.New(*webhookUrl)

	http.HandleFunc("/webhook", web.NewWebhookHandler(sb, notifier))

	log.Info().Str("listen", *listenAddr).Msg("starting server")
	if err := http.ListenAndServe(*listenAddr, nil); err != nil {
		log.Error().Err(err).Msg("server error")
	}
}
