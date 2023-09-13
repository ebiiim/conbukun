package presence

import (
	"context"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var lg zerolog.Logger = log.With().Str("component", "conbukun/pkg/presence").Logger()

// log keys
const (
	lkGuild = "guild"
)

func PresenceUpdateLoop(ctx context.Context, s *discordgo.Session) {
	next := make(chan struct{})
	go func() {
		next <- struct{}{}
	}()
	for {
		select {
		case <-ctx.Done():
			lg.Debug().Msgf("PresenceUpdateLoop <-ctx.Done()")
			if err := SetSelfPresence(s, discordgo.StatusInvisible, "", "", 0); err != nil {
				lg.Error().Err(err).Msgf("SetSelfPresence")
			}
			return
		case <-next:
			lg.Debug().Msgf("PresenceUpdateLoop <-next")

			// TODO
			if err := SetSelfPresence(s, discordgo.StatusOnline, "お散歩", "", discordgo.ActivityTypeGame); err != nil {
				lg.Error().Err(err).Msgf("SetSelfPresence")
			}
			time.AfterFunc(time.Minute*30, func() {
				next <- struct{}{}
			})

		}
	}
}

func SetSelfPresence(s *discordgo.Session, status discordgo.Status, actName, actState string, actType discordgo.ActivityType) error {
	lg.Debug().Msgf("set self presence status=%v actName=%v actState=%v", status, actName, actState)

	// FIXME: status doesn't change
	// See also: s.UpdateGameStatus
	st := discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{},
		Status:     string(status),
	}

	// NOTE:
	// ActivityTypeCompeting ～に参戦中です (～に参戦中です)
	// ActivityTypeGame	～をプレイ中 (ゲームのプレイ)
	// ActivityTypeStreaming ～をプレイ中 (TWITCHでライブ)
	// ActivityTypeListening ～を再生中 (～を再生中)
	// ActivityTypeWatching ～を視聴中 (～を視聴中)
	// ActivityTypeCustom State (State)
	if actName != "" {
		st.Activities = append(st.Activities, &discordgo.Activity{
			Name:  actName,
			State: actState,
			Type:  actType,
		})
	}

	return s.UpdateStatusComplex(st)
}
