package main

import (
	"fmt"
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

func contentToHTML(m *discordgo.Message) string {
	wg := sync.WaitGroup{}

	content := m.ContentWithMentionsReplaced()
	data := messageContentTemplateData{
		Attachments: make([]messageAttachmentTemplateData, len(m.Attachments)),
		Embeds:      make([]messageEmbedTemplateData, len(m.Embeds)),
	}

	data.Content = template.HTML(
		parseEmojis(MDtoHTML(content)),
	)

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

			if e.Color == 0 {
				e.Color = 14408667
			}

			embed := messageEmbedTemplateData{
				PillColor: fmt.Sprintf("#%X", e.Color),
				Title:     e.Title,
				TitleURL:  e.URL,
				Description: template.HTML(
					MDtoHTML(
						parseEmojis(
							e.Description,
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

			if e.Video != nil {
				embed.VideoURL = e.Video.ProxyURL
				embed.VideoOriginal = e.Video.URL
			}

			if e.Footer != nil {
				embed.Footer = e.Footer.Text
				embed.FooterIcon = e.Footer.ProxyIconURL
			}

			if e.Timestamp != "" {
				embed.Footer += " - "
				embed.Footer += e.Timestamp
			}

			data.Embeds[i] = embed
		}(i, e)
	}

	wg.Wait()

	return RenderToString(data)
}

func messageToHTML(m *discordgo.Message) string {
	wg := sync.WaitGroup{}

	var (
		color   = 16711422
		author  = safeAuthor(m)
		message = messageTemplateData{}
	)

	wg.Add(1)
	go func(m *discordgo.Message) {
		defer wg.Done()
		author, color = getUserData(m)
	}(m)

	message.ID = m.ID

	if m.Author != nil {
		message.AuthorID = m.Author.ID
		message.AvatarName = m.Author.Avatar
	}

	message.Blocked = func() bool {
		if rstore.Check(m.Author, RelationshipBlocked) {
			return true
		}

		return false
	}()

	message.Timestamp = func() (ts string) {
		stamp, err := m.Timestamp.Parse()
		if err == nil {
			ts = stamp.Format(time.Kitchen)
		}

		return
	}()

	message.Message = template.HTML(
		contentToHTML(m),
	)

	wg.Wait()

	message.NameColor = fmt.Sprintf("#%X", color)
	message.DisplayName = author

	return RenderToString(message)
}
