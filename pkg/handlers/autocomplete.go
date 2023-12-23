package handlers

import (
	"fmt"
	"slices"

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

	// NOTE:
	//  To find "Brecilien", it is necessary to add "Brecilien" before
	//  adding other names prefixed with "Brecilien" (e.g. "Brecilien Market").
	//  I don't know why this happens (a bug?), but anyway we can sort the
	//  entries in alphabetical order here to avoid the problem.
	cs := make([]string, len(c))
	i := 0
	for k := range c {
		cs[i] = k
		i++
	}
	slices.Sort(cs)

	for _, k := range cs {
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
