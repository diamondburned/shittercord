package main

import (
	"fmt"
	"html"
	"html/template"
	"log"
	"strings"
	"sync"

	"github.com/RumbleFrog/discordgo"
)

var (
	currentChannel int64
)

func loadMsgs(chID int64) {
	currentChannel = 0

	msgs, err := d.ChannelMessages(chID, 25, 0, 0, 0)
	if err != nil {
		log.Println("Failed to fetch message", err)
		return
	}

	messages := make([]string, len(msgs))

	// reverse
	for i := len(msgs)/2 - 1; i >= 0; i-- {
		opp := len(msgs) - 1 - i
		msgs[i], msgs[opp] = msgs[opp], msgs[i]
	}

	wg := sync.WaitGroup{}

	for i, m := range msgs {
		if rstore.Check(m.Author, RelationshipBlocked) {
			continue
		}

		wg.Add(1)
		go func(m *discordgo.Message, i int) {
			defer wg.Done()

			author, color := getUserData(m) // todo: put this in a goroutine

			data := messageTemplateData{
				ID:          m.ID,
				AuthorID:    m.Author.ID,
				AvatarName:  m.Author.Avatar,
				NameColor:   fmt.Sprintf("#%X", color),
				DisplayName: html.EscapeString(author),
				Content:     template.HTML(MDtoHTML(m.Content)),
			}

			messages[i] = RenderToString(data)
		}(m, i)
	}

	wg.Wait()

	SetHTML(
		GetElementByCSS(".messages"),
		strings.Join(messages, "\n"),
	)

	currentChannel = chID
}
