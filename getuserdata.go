package main

import (
	"sort"

	"github.com/RumbleFrog/discordgo"
)

func getUserData(m *discordgo.Message) (name string, color int) {
	color = 16711422
	name = m.Author.Username

	if d == nil {
		return
	}

	guild, err := d.State.Guild(m.GuildID)
	if err != nil {
		if guild, err = d.Guild(m.ChannelID); err != nil {
			// log.Println(err)
			return
		}
	}

	member, err := d.State.Member(guild.ID, m.Author.ID)
	if err != nil {
		if member, err = d.GuildMember(guild.ID, m.Author.ID); err != nil {
			// log.Println(err)
			return
		}
	}

	if member.Nick != "" {
		name = member.Nick
	}

	roles := guild.Roles
	sort.Slice(roles, func(i, j int) bool {
		return roles[i].Position > roles[j].Position
	})

	for _, role := range roles {
		for _, roleID := range member.Roles {
			if role.ID == roleID {
				if role.Color != 0 {
					color = role.Color

					return
				}
			}
		}
	}

	return
}
