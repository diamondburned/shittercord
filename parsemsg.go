package main

import (
	"fmt"
	"html"
	"html/template"
	"log"
	"regexp"
	"strings"
	"sync"

	"github.com/RumbleFrog/discordgo"
)

var (
	// EmojiRegex to get emoji IDs
	EmojiRegex = regexp.MustCompile(`&lt;(.*?):.*?:(\d+)&gt;`)
)

func messageToHTML(m *discordgo.Message) string {
	wg := sync.WaitGroup{}

	var (
		color   = 16711422
		author  = m.Author.Username
		content = html.EscapeString(m.ContentWithMentionsReplaced())
	)

	wg.Add(1)
	go func(m *discordgo.Message) {
		defer wg.Done()
		author, color = getUserData(m)
	}(m)

	// Todo: get a better ContentWithMentionsReplaced
	emojiIDs := EmojiRegex.FindAllStringSubmatch(content, -1)
	for _, nameandID := range emojiIDs {
		var format = "png"
		if nameandID[1] != "" {
			format = "gif"
		}

		content = strings.Replace(
			content,
			nameandID[0],
			fmt.Sprintf(
				`<img class="emoji" src="https://cdn.discordapp.com/emojis/%s.%s?v=1" />`,
				nameandID[2], format,
			), 1,
		)

		log.Println(fmt.Sprintf(
			`<img class="emoji" src="https://cdn.discordapp.com/emojis/%s.%s?v=1" />`,
			nameandID[2], format,
		))
	}

	wg.Wait()

	data := messageTemplateData{
		ID:          m.ID,
		AuthorID:    m.Author.ID,
		AvatarName:  m.Author.Avatar,
		NameColor:   fmt.Sprintf("#%X", color),
		DisplayName: html.EscapeString(author),
		Content:     template.HTML(MDtoHTML(content)),
	}

	return RenderToString(data)
}
