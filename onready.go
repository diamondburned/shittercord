package main

import (
	"fmt"
	"log"

	"github.com/RumbleFrog/discordgo"
)

func onReady(s *discordgo.Session, r *discordgo.Ready) {
	rstore.Relationships = r.Relationships

	var tpl = []guildsTemplateData{}

	var guildIDs []int64

	if settings, err := s.UserSettings(); err == nil {
		// switch settings.Status {
		// case discordgo.StatusOnline:
		// 	Busy = false
		// default:
		// 	Busy = true
		// }

		// setStatusButton()

		guildIDs = settings.GuildPositions
	} else {
		for _, g := range d.State.Guilds {
			guildIDs = append(guildIDs, g.ID)
		}
	}

	tpl = append(
		tpl,
		guildsTemplateData{
			ID:   -1,
			URL:  "https://raw.githubusercontent.com/google/material-design-icons/master/social/1x_web/ic_people_black_48dp.png",
			Name: "Direct Messages",
		},
	)

	for _, guildID := range guildIDs {
		guild, err := d.State.Guild(guildID)
		if err != nil {
			// if guild, err = d.Guild(guildID); err != nil {
			log.Println(err)
			continue
			// }
		}

		tpl = append(
			tpl,
			guildsTemplateData{
				ID: guild.ID,
				URL: fmt.Sprintf(
					"https://cdn.discordapp.com/icons/%d/%s.png",
					guild.ID, guild.Icon,
				),
				Name: guild.Name,
			},
		)
	}

	SetHTML(
		GetElementByCSS(".bottom-grid-wrap"),
		RenderToString(tpl),
	)
}
