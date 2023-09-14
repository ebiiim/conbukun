package handlers

import (
	"math/rand"
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
	if detectSayHello(m.Content) {
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
		if rand.Intn(100) < 5 {
			reply = "にゃー" // 5%
		} else {
			reply = "わん"
		}
	case containsWords(m.Content, wordsOha) || containsWords(m.Content, wordsKonn) || containsWords(m.Content, wordsBanwa) || containsWords(m.Content, wordsOyasu) || containsWords(m.Content, wordsOtiru):
		if rand.Intn(100) < 20 {
			reply = "わん！" // 20%
		} else {
			reply = ""
		}
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
