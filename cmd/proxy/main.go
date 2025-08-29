package main

import (
	_ "embed"
	"flag"
	"net/http"
	"os"
	"text/template"

	"github.com/jmhobbs/dayz-crimson-zamboni-deathmatch-webhook-proxy/pkg/discord"
	web "github.com/jmhobbs/dayz-crimson-zamboni-deathmatch-webhook-proxy/pkg/http"
	"github.com/jmhobbs/dayz-crimson-zamboni-deathmatch-webhook-proxy/pkg/scoreboard"

	"github.com/peterbourgon/ff/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//go:embed summary.tmpl
var defaultSummary string

func main() {
	fs := flag.NewFlagSet("dayz-crimson-zamboni-deathmatch-webhook-proxy", flag.ExitOnError)
	var (
		listenAddr      = fs.String("listen", ":8080", "listen address")
		webhookUrl      = fs.String("webhook-url", "", "Discord webhook URL")
		summaryTemplate = fs.String("summary-template", "", "Template for the summary message sent at the end of a game (optional)")
		debug           = fs.Bool("debug", false, "log debug information")
		_               = fs.String("config", "", "config file (optional)")
	)

	err := ff.Parse(fs, os.Args[1:],
		ff.WithEnvVarNoPrefix(),
		ff.WithConfigFileFlag("config"),
		ff.WithConfigFileParser(ff.PlainParser),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse flags or config")
	}

	if *webhookUrl == "" {
		log.Fatal().Msg("webhook-url is required")
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	sb := scoreboard.New()
	notifier := discord.New(*webhookUrl)

	var summary *template.Template
	if *summaryTemplate != "" {
		summary, err = template.ParseFiles(*summaryTemplate)
		if err != nil {
			log.Fatal().Err(err).Str("filename", *summaryTemplate).Msg("failed to load summary template from file")
		}
	} else {
		summary, err = template.New("summary").Parse(defaultSummary)
		if err != nil {
			panic(err)
		}
	}

	http.HandleFunc("/webhook", web.NewWebhookHandler(sb, notifier, summary))

	log.Info().Str("listen", *listenAddr).Msg("starting server")
	if err := http.ListenAndServe(*listenAddr, nil); err != nil {
		log.Error().Err(err).Msg("server error")
	}
}
