package main

import (
	"embed"
	"os"
	"text/template"

	"github.com/jmhobbs/dayz-crimson-zamboni-deathmatch-webhook-proxy/pkg/scoreboard"

	"github.com/rs/zerolog/log"
)

//go:embed normal-match.json
var f embed.FS

func main() {
	summary, err := template.ParseFiles(os.Args[1])
	if err != nil {
		log.Fatal().Err(err).Msg("unable to load summary template")
	}

	r, err := f.Open("normal-match.json")
	if err != nil {
		log.Fatal().Err(err).Msg("unable to open sample JSON")
	}

	sbd, err := scoreboard.NewFromJSON(r)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to parse sample JSON")
	}

	err = summary.Execute(os.Stdout, sbd)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to execute summary template")
	}
}
