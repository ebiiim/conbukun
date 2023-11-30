package bot

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/ebiiim/conbukun/pkg/handlers"
	"github.com/ebiiim/conbukun/pkg/presence"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var lg zerolog.Logger = log.With().Str("component", "conbukun/pkg/bot").Logger()

type Config struct {
	Version                       string
	Verbose                       int
	Token                         string
	GuildID                       string
	ROANavHandlerSaveFile         string
	ROANavHandlerSuggestionsLimit int
}

type Bot struct {
	onReadyFn              func(s *discordgo.Session, r *discordgo.Ready)
	onMessageCreateFn      func(s *discordgo.Session, m *discordgo.MessageCreate)
	onInteractionCreateFn  func(s *discordgo.Session, i *discordgo.InteractionCreate)
	onMessageReactionAddFn func(s *discordgo.Session, r *discordgo.MessageReactionAdd)

	applicationCommands []*discordgo.ApplicationCommand

	presenceUpdateLoopFn func(ctx context.Context, s *discordgo.Session)
}

func New(cfg Config) (*Bot, error) {

	roaNavHandler, err := handlers.NewROANavHandler(cfg.ROANavHandlerSaveFile, cfg.ROANavHandlerSuggestionsLimit)
	if err != nil {
		return nil, err
	}

	appCmds := []*discordgo.ApplicationCommand{
		handlers.AppCmdHelp,
		handlers.AppCmdMule,
		roaNavHandler.CommandRouteAdd(),
		roaNavHandler.CommandRoutePrint(),
		roaNavHandler.CommandRouteClear(),
	}

	cmdHandlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		handlers.CmdHelp:       handlers.HandleCmdHelp,
		handlers.CmdMule:       handlers.HandleCmdMule,
		handlers.CmdRouteAdd:   roaNavHandler.HandleCmdRouteAdd,
		handlers.CmdRoutePrint: roaNavHandler.HandleCmdRoutePrint,
		handlers.CmdRouteClear: roaNavHandler.HandleCmdRouteClear,
	}

	reactionAddHandlers := map[string]func(s *discordgo.Session, r *discordgo.MessageReactionAdd){}
	for _, emoji := range handlers.EmojisReactionAddReactionRequired {
		reactionAddHandlers[emoji] = handlers.HandleReactionAddReactionRequired
	}

	onReady := handlers.InitializeOnReadyHandler()
	onMessageCreate := handlers.NewOnMessageCreateHandler()
	onInteractionCreate := handlers.NewOnInteractionCreateHandler(cmdHandlers)
	onMessageReactionAdd := handlers.NewOnMessageReactionAddHandler(reactionAddHandlers)

	return &Bot{
		onReadyFn:              onReady,
		onMessageCreateFn:      onMessageCreate,
		onInteractionCreateFn:  onInteractionCreate,
		onMessageReactionAddFn: onMessageReactionAdd,
		applicationCommands:    appCmds,
		presenceUpdateLoopFn:   presence.PresenceUpdateLoop,
	}, nil
}

func (b *Bot) Run(cfg Config) error {
	s, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		return err
	}
	lg.Info().Msgf("bot created")

	if cfg.GuildID != "" {
		lg.Info().Str("guild_id", cfg.GuildID).Msgf("guild id specified")
	}

	lg.Info().Msgf("adding handlers...")
	s.AddHandler(b.onReadyFn)
	s.AddHandler(b.onMessageCreateFn)
	s.AddHandler(b.onInteractionCreateFn)
	s.AddHandler(b.onMessageReactionAddFn)

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
	for _, v := range b.applicationCommands {
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
		b.presenceUpdateLoopFn(ctx, s)
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
