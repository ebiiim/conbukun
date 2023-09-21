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
	isReply := false
	if m.MessageReference != nil && m.ReferencedMessage.Author.ID == s.State.Application.ID {
		isReply = true
	}

	// check function
	funcName := ""
	if detectSayHello(m.Content) {
		funcName = FuncMessageCreateSayHello
	}
	if isReply {
		funcName = FuncMessageCreateReplyHello
	}

	lg.Debug().Bool(lkDM, isDM).Msgf("isMention=%v isThread=%v hasRef=%v isReply=%v content=%s", isMention, isThread, hasRef, isReply, m.Content)

	if funcName == "" {
		return
	}

	lg.Info().Str(lkFunc, funcName).Bool(lkDM, isDM).Msg("OnMessageCreate")
	switch funcName {
	case FuncMessageCreateSayHello:
		handleMessageCreateSayHello(s, m)
	case FuncMessageCreateReplyHello:
		handleMessageCreateReplyHello(s, m)
	default:
		return
	}
}

func detectSayHello(s string) bool {
	return false ||
		containsWords(s, wordsConbu) ||
		containsWords(s, wordsOha) ||
		containsWords(s, wordsKonn) ||
		containsWords(s, wordsBanwa) ||
		containsWords(s, wordsOyasu) ||
		containsWords(s, wordsOtiru) ||
		false
}

var (
	wordsConbu = []string{"こんぶくん"}
	wordsOha   = []string{"おはよ", "おはです", "おはます"}
	wordsKonn  = []string{"こんにちは", "こんにちわ", "こんちは", "こんちわ", "こんちゃ"}
	wordsBanwa = []string{"こんばんは", "こんばんわ", "こんばわ", "ばんちゃ"}
	wordsOyasu = []string{"おやす", "寝ま"}
	wordsOtiru = []string{"おちま", "落ちま", "おちる", "落ちる", "お先", "おさき"}
)

func containsWords(s string, words []string) bool {
	for _, word := range words {
		if strings.Contains(s, word) {
			return true
		}
	}
	return false
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

	reply := ""

	switch {
	case containsWords(m.Content, wordsConbu):
		reply = pickOne([]choice{
			{80, "わん"},
			{5, "にゃー"},
			{15, "ぱぱ？"},
		})
	case containsWords(m.Content, wordsOha) || containsWords(m.Content, wordsKonn) || containsWords(m.Content, wordsBanwa) || containsWords(m.Content, wordsOyasu) || containsWords(m.Content, wordsOtiru):
		reply = pickOne([]choice{
			{17, "わん！"},
			{3, "にゃー！"},
			{80, ""},
		})
	}

	if reply == "" {
		return
	}

	// send reply
	if err := s.ChannelTyping(m.ChannelID); err != nil {
		lg.Error().Err(err).Msg("could not send typing")
	}
	time.Sleep(time.Second)
	_, err := sendSilentMessage(s, m.ChannelID, &discordgo.MessageSend{
		Content:   reply,
		Reference: m.Reference(),
	})
	if err != nil {
		lg.Error().Err(err).Msg("could not send msg")
	}
}

func handleMessageCreateReplyHello(s *discordgo.Session, m *discordgo.MessageCreate) {
	lg := lg.With().
		Str(lkFunc, FuncMessageCreateReplyHello).
		Str(lkMID, m.ID).
		Str(lkGuild, m.GuildID).
		Str(lkCh, m.ChannelID).
		Str(lkUsr, m.Author.ID).
		Str(lkName, m.Author.Username).
		Logger()

	lg.Info().Msgf("MessageCreateReplyHello: called")

	reply := ""

	switch {
	default:
		reply = pickOne([]choice{
			{25, "わん！！"},
			{15, "にゃー？"},
			{60, "ぱぱ！"},
		})
	}

	if reply == "" {
		return
	}

	// send reply
	if err := s.ChannelTyping(m.ChannelID); err != nil {
		lg.Error().Err(err).Msg("could not send typing")
	}
	time.Sleep(time.Second)
	_, err := sendSilentMessage(s, m.ChannelID, &discordgo.MessageSend{
		Content:   reply,
		Reference: m.Reference(),
	})
	if err != nil {
		lg.Error().Err(err).Msg("could not send msg")
	}
}
