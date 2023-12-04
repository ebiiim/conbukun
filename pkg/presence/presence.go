package presence

import (
	"context"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Version = "dev" // main.go injects this value

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

			a0, a1, a2, a3, a4 := generatePresence(time.Now())
			if err := SetSelfPresence(s, a0, a1, a2, a3); err != nil {
				lg.Error().Err(err).Msgf("SetSelfPresence")
			}
			lg.Debug().Msgf("PresenceUpdateLoop wait for %s", a4)
			time.AfterFunc(a4, func() {
				next <- struct{}{}
			})

		}
	}
}

var (
	jst = time.FixedZone("Asia/Tokyo", 9*60*60)
)

const (
	actWalking  = "お散歩"
	actToilet   = "トイレ"
	actNap      = "お昼寝"
	actFood     = "ごはん"
	actSnack    = "おやつ"
	actMTG      = "MTG"
	actTraining = "運動"
	actSleep    = "すやすや"
	actAO       = "Albion Online"
)

var (
	actPrefix = map[string]bool{
		actWalking:  true,
		actToilet:   false,
		actNap:      false,
		actFood:     false,
		actSnack:    false,
		actMTG:      true,
		actTraining: true,
		actSleep:    false,
		actAO:       false,
	}
	actDuration = map[string]time.Duration{
		actWalking:  time.Minute * 25,
		actToilet:   time.Minute * 5,
		actNap:      time.Minute * 45,
		actFood:     time.Minute * 25,
		actSnack:    time.Minute * 15,
		actMTG:      time.Minute * 30,
		actTraining: time.Minute * 30,
		actSleep:    time.Minute * 80,
		actAO:       time.Minute * 50,
	}

	actMorning = []string{actWalking, actToilet, actFood, actSnack, actTraining, actAO}
	actEvening = []string{actWalking, actToilet, actNap, actFood, actSnack, actMTG, actAO}
	actNight   = []string{actWalking, actToilet, actFood, actSnack, actMTG, actSleep, actAO}
)

func generatePresence(t time.Time) (status discordgo.Status, actName, actState string, actType discordgo.ActivityType, waitNext time.Duration) {
	status = discordgo.StatusOnline
	actState = ""
	actType = discordgo.ActivityTypeGame
	waitNext = time.Minute * 30

	w0 := "時間帯不明" // 朝
	w1 := "行方不明"  // お散歩

	h := t.In(jst).Hour()
	// m := t.In(jst).Minute()

	switch {
	case h >= 0 && h < 4: // 0:00-3:59
		w0 = "深夜"
		w1 = actNight[rand.Intn(len(actNight))]
	case h >= 4 && h < 6: // 4:00-5:59
		w0 = "早朝"
		w1 = actMorning[rand.Intn(len(actMorning))]
	case h >= 6 && h < 11: // 6:00-10:59
		w0 = "朝"
		w1 = actMorning[rand.Intn(len(actMorning))]
	case h >= 11 && h < 15: // 11:00-14:59
		w0 = "昼"
		w1 = actEvening[rand.Intn(len(actEvening))]
	case h >= 15 && h < 19: // 15:00-18:59
		w0 = "夕方"
		w1 = actEvening[rand.Intn(len(actEvening))]
	case h >= 19 && h < 24: // 19:00-23:59
		w0 = "夜"
		w1 = actNight[rand.Intn(len(actNight))]
	}

	if dur, ok := actDuration[w1]; ok {
		waitNext = dur
	}
	if hasPrefix, ok := actPrefix[w1]; ok && hasPrefix {
		actName = w0 + "の" + w1 // 朝のお散歩
	} else {
		actName = w1 // お昼寝
	}

	return
}

func SetSelfPresence(s *discordgo.Session, status discordgo.Status, actName, actState string, actType discordgo.ActivityType) error {
	lg.Info().Msgf("set self presence status=%v actName=%v actState=%v", status, actName, actState)

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
