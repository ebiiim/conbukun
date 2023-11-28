package handlers

import (
	"fmt"

	trie "github.com/Vivino/go-autocomplete-trie"
	"github.com/bwmarrin/discordgo"

	"github.com/ebiiim/conbukun/pkg/ao/data"
)

type MapNameCompleter struct {
	c   map[string]*discordgo.ApplicationCommandOptionChoice
	t   *trie.Trie
	lim int
}

func NewMapNameCompleter(lim int) *MapNameCompleter {
	c := maps2choices(data.Maps)
	t := trie.New()
	t = t.CaseSensitive()

	for k := range c {
		t.Insert(k)
	}

	return &MapNameCompleter{c: c, t: t, lim: lim}
}

func (c *MapNameCompleter) GetSuggestions(input string) []string {
	return c.t.Search(input, c.lim)
}

func (c *MapNameCompleter) GetChoices(input string) []*discordgo.ApplicationCommandOptionChoice {
	suggestions := c.GetSuggestions(input)
	choices := []*discordgo.ApplicationCommandOptionChoice{}
	for _, s := range suggestions {
		choices = append(choices, c.c[s])
	}
	return choices
}

func maps2choices(maps map[string]data.MapData) map[string]*discordgo.ApplicationCommandOptionChoice {
	choices := map[string]*discordgo.ApplicationCommandOptionChoice{}
	for _, md := range maps {
		choiceName := ""
		choiceValue := md.ID

		switch data.GetMapType(md) {
		case data.MapTypeBlackZone, data.MapTypeRedZone, data.MapTypeYellowZone, data.MapTypeBlueZone, data.MapTypeCity:
			choiceName = md.DisplayName
		case data.MapTypeAvalon:
			choiceName = fmt.Sprintf("%s: %s", data.GetMapShortName(md), md.DisplayName)
		default:
			continue
		}

		choices[choiceName] = &discordgo.ApplicationCommandOptionChoice{Name: choiceName, Value: choiceValue}
	}
	return choices
}
