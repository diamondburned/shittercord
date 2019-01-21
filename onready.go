package main

import (
	"fmt"
	"log"

	"github.com/RumbleFrog/discordgo"
)

var (
	UserSettings *discordgo.Settings
)

func onReady(s *discordgo.Session, r *discordgo.Ready) {
	rstore.Relationships = r.Relationships

	var (
		tpl      = []guildsTemplateData{}
		guildIDs []int64
		err      error
	)

	UserSettings, err = s.UserSettings()
	if err == nil {
		// Set avatar status
		if selfAvatar := GetElementByCSS("#self-avatar"); selfAvatar != nil {
			selfAvatar.SetStyle(
				"border-color",
				ReflectStatusColor(UserSettings.Status),
			)

			selfAvatar.SetAttr(
				"src",
				fmt.Sprintf(
					"https://cdn.discordapp.com/avatars/%d/%s.png?size=64",
					r.User.ID,
					r.User.Avatar,
				),
			)
		}

		// Set self username
		SetText(
			GetElementByCSS("#self-username"),
			r.User.Username+"#"+r.User.Discriminator,
		)

		guildIDs = UserSettings.GuildPositions

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

// ReflectStatusColor converts Discord status to HEX colors (#RRGGBB)
func ReflectStatusColor(status discordgo.Status) string {
	switch status {
	case discordgo.StatusOnline:
		return "#43b581"
	case discordgo.StatusDoNotDisturb:
		return "#f04747"
	case discordgo.StatusIdle:
		return "#faa61a"
	default: // includes invisible and offline
		return "#747f8d"
	}
}
