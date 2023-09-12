package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/ebiiim/conbukun/pkg/handlers"
)

func run(gid string, token string) error {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		return err
	}
	lg.Info().Msgf("bot created")

	lg.Info().Msgf("adding handlers...")
	s.AddHandler(handlers.OnReady)
	s.AddHandler(handlers.OnMessageCreate)
	s.AddHandler(handlers.OnInteractionCreate)

	if err := s.Open(); err != nil {
		return err
	}
	defer func() {
		if err := s.Close(); err != nil {
			lg.Error().Err(err).Msgf("could not close the bot session")
		}
	}()
	lg.Info().Msgf("bot session opened")

	lg.Info().Msg("creating commands...")
	for _, v := range handlers.Commands {
		lg.Debug().Msgf("creating command: %s", v.Name)
		if _, err := s.ApplicationCommandCreate(s.State.User.ID, gid, v); err != nil {
			lg.Error().Err(err).Msgf("could not create command: %s", v.Name)
		}
	}

	lg.Info().Msgf("bot started (CTRL+C to stop)")
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	<-ctx.Done() // wait for SIGINT
	lg.Info().Msg("SIGINT received")

	lg.Info().Msg("removing commands...")
	registeredCommands, err := s.ApplicationCommands(s.State.User.ID, gid)
	if err != nil {
		lg.Error().Err(err).Msg("could not fetch registered commands")
	} else {
		for _, cmd := range registeredCommands {
			lg.Debug().Msgf("removing command: %s(%s)", cmd.ID, cmd.Name)
			if err := s.ApplicationCommandDelete(s.State.User.ID, gid, cmd.ID); err != nil {
				lg.Error().Err(err).Msgf("could not delete command: %s", cmd.Name)
			}
		}
	}

	lg.Info().Msg("bot stopped")
	return nil
}

var version = "dev"

var lg zerolog.Logger = log.With().Str("component", "Conbukun Bot (main)").Logger()

func main() {
	flag.CommandLine.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "Usage: %s [OPTION]... \n\n", flag.CommandLine.Name())
		flag.PrintDefaults()
		fmt.Fprintf(w, "\nVersion: %s\n", version)
	}

	var logLevel int
	flag.IntVar(&logLevel, "v", 3, "log level")
	var gid string
	flag.StringVar(&gid, "gid", "", "Guild ID or registers commands globally")
	var token string
	flag.StringVar(&token, "token", "", "Bot authentication token")
	flag.Parse()

	envGID := os.Getenv("CONBUKUN_GUILD_ID")
	envToken := os.Getenv("CONBUKUN_AUTH_TOKEN")
	if gid == "" {
		gid = envGID
	}
	if token == "" {
		token = envToken
	}

	switch logLevel {
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

	if err := run(gid, token); err != nil {
		lg.Fatal().Err(err).Msg("stopped")
	}
}
