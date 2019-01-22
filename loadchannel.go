package main

import (
	"log"
	"sort"
	"strings"
	"sync"

	"github.com/RumbleFrog/discordgo"
)

var (
	currentChannel int64
	currentGuild   int64
)

func loadMsgs(chID int64) {
	currentChannel = 0

	go func(chID int64) {
		ch, err := d.State.Channel(chID)
		if err == nil {
			SetText(GetElementByCSS("#channel-name"), ch.Name)
			currentGuild = ch.GuildID
		}
	}(chID)

	msgs := []*discordgo.Message{}

	ch, err := d.State.Channel(chID)
	if err == nil {
		if len(ch.Messages) > 25 {
			msgs = ch.Messages
			goto Continue
		}
	}

	msgs, err = d.ChannelMessages(chID, 50, 0, 0, 0)
	if err != nil {
		log.Println("Failed to fetch message", err)
		return
	}

Continue:
	messages := make([]string, len(msgs))

	// reverse
	sort.Slice(msgs, func(i, j int) bool {
		return msgs[i].ID < msgs[j].ID
	})

	//	for i := len(msgs)/2 - 1; i >= 0; i-- {
	//		opp := len(msgs) - 1 - i
	//		msgs[i], msgs[opp] = msgs[opp], msgs[i]
	//	}

	wg := sync.WaitGroup{}

	for i, m := range msgs {
		wg.Add(1)
		go func(m *discordgo.Message, i int) {
			defer wg.Done()
			messages[i] = messageToHTML(m)
		}(m, i)
	}

	wg.Wait()

	SetHTML(
		GetElementByCSS(".messages"),
		strings.Join(messages, "\n"),
	)

	currentChannel = chID

	if _, err := w.Call(
		"scrollToBottom",
	); err != nil {
		log.Println(err)
		return
	}
}
