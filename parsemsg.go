package main

import (
	"fmt"
	"html"
	"html/template"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/RumbleFrog/discordgo"
	humanize "github.com/dustin/go-humanize"
)

var (
	// EmojiRegex to get emoji IDs
	EmojiRegex = regexp.MustCompile(`&lt;(.*?):.*?:(\d+)&gt;`)
)

func parseEmojis(content string) string {
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
	}

	return content
}

func messageToHTML(m *discordgo.Message) string {
	wg := sync.WaitGroup{}

	var (
		color   = 16711422
		author  = safeAuthor(m)
		content = html.EscapeString(m.ContentWithMentionsReplaced())
		data    = messageTemplateData{
			Attachments: make([]messageAttachmentTemplateData, len(m.Attachments)),
			Embeds:      make([]messageEmbedTemplateData, len(m.Embeds)),
		}
	)

	wg.Add(1)
	go func(m *discordgo.Message) {
		defer wg.Done()
		author, color = getUserData(m)
	}(m)

	content = parseEmojis(content)

	for i, a := range m.Attachments {
		wg.Add(1)

		go func(i int, a *discordgo.MessageAttachment) {
			defer wg.Done()

			data.Attachments[i] = messageAttachmentTemplateData{
				Name:     a.Filename,
				URL:      a.ProxyURL,
				Original: a.URL,
				Size:     humanize.Bytes(uint64(a.Size)),
				MediaType: func() MediaType {
					if a.Height > 0 && a.Width > 0 {
						return TypeImage
					}

					// Returning File as-is for now
					return TypeFile
				}(),
			}
		}(i, a)
	}

	for i, e := range m.Embeds {
		wg.Add(1)

		go func(i int, e *discordgo.MessageEmbed) {
			defer wg.Done()

			embed := messageEmbedTemplateData{
				PillColor: fmt.Sprintf("#%X", e.Color),
				Title:     e.Title,
				TitleURL:  e.URL,
				Description: template.HTML(
					MDtoHTML(
						parseEmojis(
							html.EscapeString(
								e.Description,
							),
						),
					),
				),
				Fields: e.Fields,
			}

			if e.Author != nil {
				embed.Author = e.Author.Name
				embed.AuthorURL = e.Author.URL
				embed.AuthorIcon = e.Author.ProxyIconURL
			}

			if e.Thumbnail != nil {
				embed.ThumbnailURL = e.Thumbnail.ProxyURL
				embed.ThumbnailOriginal = e.Thumbnail.URL
				embed.ThumbnailWidth = e.Thumbnail.Width
				embed.ThumbnailHeight = e.Thumbnail.Height
			}

			if e.Image != nil {
				embed.ImageURL = e.Image.ProxyURL
				embed.ImageOriginal = e.Image.URL
				embed.ImageWidth = e.Image.Width
				embed.ImageHeight = e.Image.Height
			}

			if e.Footer != nil {
				embed.Footer = e.Footer.Text
				embed.FooterIcon = e.Footer.ProxyIconURL
			}

			data.Embeds[i] = embed
		}(i, e)
	}

	wg.Wait()

	data.ID = m.ID
	data.AuthorID = m.Author.ID
	data.AvatarName = m.Author.Avatar
	data.NameColor = fmt.Sprintf("#%X", color)
	data.DisplayName = html.EscapeString(author)
	data.Content = template.HTML(MDtoHTML(content))

	data.Blocked = func() bool {
		if rstore.Check(m.Author, RelationshipBlocked) {
			return true
		}

		return false
	}()

	data.Timestamp = func() (ts string) {
		stamp, err := m.Timestamp.Parse()
		if err == nil {
			ts = stamp.Format(time.Kitchen)

			if m.EditedTimestamp != "" {
				if stamp, err := m.EditedTimestamp.Parse(); err == nil {
					ts += " (edited: " + stamp.Format(time.Kitchen) + ")"
				}
			}
		}

		return
	}()

	return RenderToString(data)
}
