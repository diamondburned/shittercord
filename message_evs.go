package main

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/RumbleFrog/discordgo"
	"github.com/gen2brain/beeep"
	sciter "github.com/sciter-sdk/go-sciter"
)

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if userisBusy() {
		go func(m *discordgo.Message) {
			for _, mention := range m.Mentions {
				if mention.ID == d.State.User.ID {
					if err := beeep.Notify(
						m.Author.Username+" mentioned you",
						m.ContentWithMentionsReplaced(), "",
					); err != nil {
						log.Println(err)
					}
				}
			}
		}(m.Message)
	}

	if m.ChannelID != currentChannel {
		GetElementByCSS
		return
	}

	html := messageToHTML(m.Message)

	if _, err := w.Call(
		"appendHTMLMessage",
		sciter.NewValue(html),
	); err != nil {
		log.Println(err)
		return
	}
}

func messageUpdate(s *discordgo.Session, m *discordgo.MessageUpdate) {
	if m.ChannelID != currentChannel {
		return
	}

	var (
		elem = GetElementByCSS(fmt.Sprintf("#%d", m.ID))
		html = messageToHTML(m.Message)
	)

	if elem == nil {
		return
	}

	_lines := strings.Split(html, "\n")

	SetHTML(
		elem,
		strings.Join(_lines[1:len(_lines)-2], "\n"),
	)
}

func messageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
	if m.ChannelID != currentChannel {
		return
	}

	deleteMessage(m.ID)
}

func messageDeleteBulk(s *discordgo.Session, bulk *discordgo.MessageDeleteBulk) {
	if bulk.ChannelID != currentChannel {
		return
	}

	wg := sync.WaitGroup{}

	for _, msgID := range bulk.Messages {
		wg.Add(1)

		go func(ID int64) {
			defer wg.Done()
			deleteMessage(ID)
		}(msgID)
	}

	wg.Wait()
}

// Get message from DOM and delete it
// Function because it's re-used twice
func deleteMessage(msgID int64) {
	elem := GetElementByCSS(fmt.Sprintf("#%d", msgID))
	if elem == nil {
		return
	}

	if err := elem.Delete(); err != nil {
		log.Println(err)
	}
}
