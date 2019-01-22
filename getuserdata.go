package main

import (
	"log"
	"sort"

	"github.com/RumbleFrog/discordgo"
)

func getUserData(m *discordgo.Message) (name string, color int) {
	color = 16711422
	name = safeAuthor(m)

	if d == nil {
		return
	}

	channel, err := d.State.Channel(m.ChannelID)
	if err != nil {
		if channel, err = d.Channel(m.ChannelID); err != nil {
			log.Println(err)
			return
		}
	}

	guild, err := d.State.Guild(channel.GuildID)
	if err != nil {
		if guild, err = d.Guild(channel.GuildID); err != nil {
			log.Println(err)
			return
		}
	}

	member, err := d.State.Member(guild.ID, m.Author.ID)
	if err != nil {
		if member, err = d.GuildMember(guild.ID, m.Author.ID); err != nil {
			log.Println(err)
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
			if role.ID == roleID && role.Color != 0 {
				color = role.Color
				return
			}
		}
	}

	return
}
