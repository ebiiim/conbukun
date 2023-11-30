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
	"github.com/ebiiim/conbukun/pkg/presence"
)

type Config struct {
	Verbose                       int
	Token                         string
	GuildID                       string
	ROANavHandlerSaveFile         string
	ROANavHandlerSuggestionsLimit int
}

func run(cfg Config) error {
	s, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		return err
	}
	lg.Info().Msgf("bot created")

	if cfg.GuildID != "" {
		lg.Info().Str("guild_id", cfg.GuildID).Msgf("guild id specified")
	}

	lg.Info().Msgf("adding handlers...")
	s.AddHandler(handlers.OnReady)
	s.AddHandler(handlers.OnMessageCreate)
	s.AddHandler(handlers.OnInteractionCreate)
	s.AddHandler(handlers.OnMessageReactionAdd)

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
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, cfg.GuildID, v)
		if err != nil {
			lg.Error().Err(err).Msgf("could not create command: %s", v.Name)
		}
		lg.Debug().Msgf("command created: %s(%s)", cmd.ID, cmd.Name)
	}

	lg.Info().Msgf("bot started (CTRL+C to stop)")
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	stopped := make(chan struct{})
	go func() {
		lg.Info().Msg("PresenceUpdateLoop started")
		presence.PresenceUpdateLoop(ctx, s)
		lg.Info().Msg("PresenceUpdateLoop stopped")
		stopped <- struct{}{}
	}()

	<-ctx.Done() // wait for SIGINT
	lg.Info().Msg("SIGINT received")

	lg.Info().Msg("removing commands...")
	registeredCommands, err := s.ApplicationCommands(s.State.User.ID, cfg.GuildID)
	if err != nil {
		lg.Error().Err(err).Msg("could not fetch registered commands")
	} else {
		for _, cmd := range registeredCommands {
			lg.Debug().Msgf("removing command: %s(%s)", cmd.ID, cmd.Name)
			if err := s.ApplicationCommandDelete(s.State.User.ID, cfg.GuildID, cmd.ID); err != nil {
				lg.Error().Err(err).Msgf("could not delete command: %s", cmd.Name)
			}
		}
	}

	// wait for goroutines to exit
	<-stopped
	lg.Info().Msg("all components are stopped")

	lg.Info().Msg("bot stopped")
	return nil
}

var version = "dev"

// some variables shoud be initialized in init()
func init() {
	handlers.Version = version

	// NOTE: this is a workaroud; currenltly ROANavHandler is initalized in init()
	roanavSaveFile := os.Getenv("CONBUKUN_ROANAV_SAVEFILE")
	if roanavSaveFile == "" {
		handlers.ROANavHandlerSaveFile = roanavSaveFile
	}
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
	flag.Parse()

	envToken := os.Getenv("CONBUKUN_AUTH_TOKEN")
	if token == "" {
		token = envToken
	}
	envGID := os.Getenv("CONBUKUN_GUILD_ID")
	if gid == "" {
		gid = envGID
	}

	cfg := Config{
		Verbose:                       logLevel,
		Token:                         token,
		GuildID:                       gid,
		ROANavHandlerSaveFile:         "NOT_USED",
		ROANavHandlerSuggestionsLimit: 5, // not used
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
	if err := run(cfg); err != nil {
		lg.Fatal().Err(err).Msg("stopped")
	}
}
