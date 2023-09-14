package handlers

import (
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	lg := lg.With().
		Str(lkHandler, "MessageCreate").
		Str(lkMID, m.ID).
		Str(lkGuild, m.GuildID).
		Str(lkCh, m.ChannelID).
		Str(lkUsr, m.Author.ID).
		Str(lkName, m.Author.Username).
		Logger()

	// ignore bot messages
	if m.Author.Bot {
		return
	}

	// check DM
	isDM := false
	if m.GuildID == "" {
		isDM = true
	}

	// check mention
	isMention := false
	for _, u := range m.Mentions {
		if u.ID == s.State.Application.ID {
			isMention = true
		}
	}

	// check thread
	isThread := false
	ch, err := s.Channel(m.ChannelID)
	if err != nil {
		lg.Error().Err(err).Msg("could not get channel")
	}
	if ch.IsThread() {
		isThread = true
	}

	// check ref
	hasRef := false
	if m.MessageReference != nil {
		hasRef = true
	}

	// check function
	funcName := ""
	if containsConbukun(m.Content) {
		funcName = FuncMessageCreateSayHello
	}
	if isMention {
		funcName = FuncMessageCreateSayHello
	}

	lg.Debug().Bool(lkDM, isDM).Msgf("isMention=%v isThread=%v hasRef=%v content=%s", isMention, isThread, hasRef, m.Content)

	if funcName == "" {
		return
	}

	lg.Info().Str(lkFunc, funcName).Bool(lkDM, isDM).Msg("OnMessageCreate")
	switch funcName {
	case FuncMessageCreateSayHello:
		handleMessageCreateSayHello(s, m)
	default:
		return
	}
}

func containsConbukun(s string) bool {
	ss := strings.ToLower(s)
	return strings.Contains(ss, "こんぶくん")
}

func handleMessageCreateSayHello(s *discordgo.Session, m *discordgo.MessageCreate) {
	lg := lg.With().
		Str(lkFunc, FuncMessageCreateSayHello).
		Str(lkMID, m.ID).
		Str(lkGuild, m.GuildID).
		Str(lkCh, m.ChannelID).
		Str(lkUsr, m.Author.ID).
		Str(lkName, m.Author.Username).
		Logger()

	lg.Info().Msgf("MessageCreateSayHello: called")

	if err := s.ChannelTyping(m.ChannelID); err != nil {
		lg.Error().Err(err).Msg("could not send typing")
	}

	time.Sleep(time.Second)

	reply := "わん"
	_, err := sendSilentMessage(s, m.ChannelID, &discordgo.MessageSend{
		Content:   reply,
		Reference: m.Reference(),
	})
	if err != nil {
		lg.Error().Err(err).Msg("could not send msg")
	}
}
