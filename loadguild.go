package main

import (
	"log"
	"sort"
	"strings"

	"github.com/RumbleFrog/discordgo"
)

type sortCats struct {
	Channels      []*discordgo.Channel
	CategoryOrder int
	CategoryName  string
	CategoryID    int64
}

func loadGuild(gid int64) {
	datum := []categoriesTemplateData{}

	if gid == -1 {
		datum = append(
			datum,
			categoriesTemplateData{
				CategoryID:   -1,
				CategoryName: "",
				Channels:     []channelsTemplateData{},
			},
		)

		for _, dmCh := range d.State.PrivateChannels {
			var userstr []string
			for _, u := range dmCh.Recipients {
				userstr = append(userstr, u.Username)
			}

			datum[0].Channels = append(
				datum[0].Channels,
				channelsTemplateData{
					ID:   dmCh.ID,
					Name: strings.Join(userstr, ", "),
				},
			)
		}

	} else {
		guild, err := d.State.Guild(gid)
		if err != nil {
			if guild, err = d.Guild(gid); err != nil {
				log.Println(err)
				return
			}
		}

		// map[ch_id][]channels
		categorizedChannels := make(map[int64]sortCats)
		categorizedChannels[0] = sortCats{
			CategoryName:  "",
			CategoryOrder: -1,
			Channels:      []*discordgo.Channel{},
		}

		for _, channel := range guild.Channels {
			switch channel.Type {
			case discordgo.ChannelTypeGuildVoice:
				continue
			case discordgo.ChannelTypeGuildCategory:
				t := categorizedChannels[channel.ID]
				t.CategoryName = channel.Name
				t.CategoryOrder = channel.Position
				t.CategoryID = channel.ID
				categorizedChannels[channel.ID] = t
			}
		}

		for _, channel := range guild.Channels {
			switch channel.Type {
			case discordgo.ChannelTypeGuildText:
				t := categorizedChannels[channel.ParentID]
				t.Channels = append(
					t.Channels,
					channel,
				)
				categorizedChannels[channel.ParentID] = t
			}
		}

		var sortedCats []sortCats
		for _, category := range categorizedChannels {
			sortedCats = append(sortedCats, category)
		}

		sort.SliceStable(sortedCats, func(i, j int) bool {

			return sortedCats[i].CategoryOrder < sortedCats[j].CategoryOrder
		})

		for _, category := range sortedCats {
			if len(category.Channels) < 1 {
				continue
			}

			data := categoriesTemplateData{
				CategoryID:   category.CategoryID,
				CategoryName: category.CategoryName,
				Channels:     []channelsTemplateData{},
			}

			sort.SliceStable(category.Channels, func(i, j int) bool {
				return category.Channels[i].Position < category.Channels[j].Position
			})

			for _, channel := range category.Channels {
				data.Channels = append(
					data.Channels,
					channelsTemplateData{
						ID:    channel.ID,
						Name:  channel.Name,
						NSFW:  channel.NSFW,
						Muted: false,
					},
				)
			}

			datum = append(datum, data)
		}
	}

	SetHTML(
		GetElementByCSS(".channel-list"),
		RenderToString(datum),
	)
}
