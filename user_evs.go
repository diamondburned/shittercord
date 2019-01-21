package main

import (
	"github.com/RumbleFrog/discordgo"
	"github.com/davecgh/go-spew/spew"
)

// redundant for now because of ID check
// todo: refer to README
func presenceUpdate(s *discordgo.Session, presence *discordgo.PresenceUpdate) {
	if presence.User.ID != d.State.User.ID {
		return
	}

	UserSettings.Status = presence.Status

	GetElementByCSS("#self-avatar").SetStyle(
		"border-color",
		ReflectStatusColor(presence.Status),
	)
}

// redundant for now because of ID check
func presencesReplace(s *discordgo.Session, pReplace *discordgo.PresencesReplace) {
	for _, presence := range *pReplace {
		if presence.User.ID != d.State.User.ID {
			return
		}

		UserSettings.Status = presence.Status

		GetElementByCSS("#self-avatar").SetStyle(
			"border-color",
			ReflectStatusColor(presence.Status),
		)
	}
}

func userSettingsUpdate(s *discordgo.Session, settings *discordgo.UserSettingsUpdate) {
	go spew.Dump(settings) // verbose debugging

	if settings == nil {
		return
	}

	_settings := *settings

	if status, ok := _settings["status"]; ok {
		if str, ok := status.(string); ok {
			st := discordgo.Status(str)
			UserSettings.Status = st

			GetElementByCSS("#self-avatar").SetStyle(
				"border-color",
				ReflectStatusColor(st),
			)
		}
	}
}

func userisBusy() bool {
	if UserSettings == nil {
		return false
	}

	switch UserSettings.Status {
	case discordgo.StatusDoNotDisturb:
		return true
	default:
		return false
	}
}
