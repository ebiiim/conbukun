package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/ebiiim/conbukun/pkg/bot"
	"github.com/ebiiim/conbukun/pkg/handlers"
)

var version = "dev"

func init() {
	handlers.Version = version // TODO: this is a workaround
}

var lg zerolog.Logger = log.With().Str("component", "Conbukun Bot").Logger()

func main() {
	flag.CommandLine.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "Usage: %s [OPTION]... \n\n", flag.CommandLine.Name())
		flag.PrintDefaults()
		fmt.Fprintf(w, "\nVersion: %s\n", version)
	}

	var logLevel int
	flag.IntVar(&logLevel, "v", 3, "log level")
	var token string
	flag.StringVar(&token, "token", "", "Bot authentication token")
	var gid string
	flag.StringVar(&gid, "gid", "", "Guild ID or registers commands globally")
	var saveDir string
	flag.StringVar(&saveDir, "save-dir", "", "Save directory")
	flag.Parse()

	envToken := os.Getenv("CONBUKUN_AUTH_TOKEN")
	if token == "" {
		token = envToken
	}
	envGID := os.Getenv("CONBUKUN_GUILD_ID")
	if gid == "" {
		gid = envGID
	}
	envSaveDir := os.Getenv("CONBUKUN_SAVE_DIR")
	if saveDir == "" {
		saveDir = envSaveDir
	}

	cfg := bot.Config{
		Version:                       version,
		Verbose:                       logLevel,
		Token:                         token,
		GuildID:                       gid,
		ROANavHandlerSaveFile:         filepath.Join(saveDir, "roanav.json"),
		ROANavHandlerSuggestionsLimit: handlers.DefaultROANavHandlerSuggestionsLimit,
	}

	switch cfg.Verbose {
	default:
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case 1:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case 2:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case 3:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case 4:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case 5:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}

	lg.Info().Msgf("conbukun version=%v", version)

	b, err := bot.New(cfg)
	if err != nil {
		lg.Fatal().Err(err).Msg("could not create bot")
	}

	if err := b.Run(cfg); err != nil {
		lg.Fatal().Err(err).Msg("stopped")
	}
}
