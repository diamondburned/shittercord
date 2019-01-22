package main

import (
	"encoding/csv"
	"log"
	"strings"

	"github.com/RumbleFrog/discordgo"
)

var (
	latestSentID int64
)

func sendMessage(content string) {
	input := csv.NewReader(strings.NewReader(content))
	input.Comma = ' ' // delimiter
	args, err := input.Read()
	if err != nil {
		log.Println(err)
		return
	}

	// if there's nothing in `content'
	if len(args) < 1 {
		return
	}

	switch args[0] {
	case "/replace", "/replaceAll":
		if latestSentID == 0 {
			return
		}

		if len(args) < 3 {
			WarnDialog("Missing argument(s)! Refer to /help")
		}

		msg, err := d.State.Message(currentChannel, latestSentID)
		if err != nil {
			log.Println(err)
			return
		}

		if msg.Author.ID != d.State.User.ID {
			log.Println("Stored message", latestSentID, "isn't yours.")
			return
		}

		replaced := strings.Replace(msg.Content, args[1], args[2], func() int {
			if strings.HasSuffix(args[0], "All") {
				return -1
			}

			return 0
		}())

		if m, err := d.ChannelMessageEdit(currentChannel, latestSentID, replaced); err != nil {
			log.Println(err)
		} else {
			latestSentID = m.ID
		}
	case "/embed":
		var embedContent = struct {
			Title       string
			Author      string
			AuthorURL   string
			AuthorImage string
			Footer      string
			Thumbnail   string
			Content     []string
		}{}

		for i := 1; i < len(args); i++ {
			switch args[i] {
			case "-t", "--title":
				embedContent.Title = args[i+1]
				i++
			case "-a", "--author":
				embedContent.Author = args[i+1]
				i++
			case "--authorURL":
				embedContent.AuthorURL = args[i+1]
				i++
			case "--authorImage":
				embedContent.AuthorImage = args[i+1]
				i++
			case "-f", "--footer":
				embedContent.Footer = args[i+1]
				i++
			case "-th", "--thumbnail":
				embedContent.Thumbnail = args[i+1]
				i++
			default:
				embedContent.Content = append(embedContent.Content, args[i])
			}
		}

		if _, err := d.ChannelMessageSendEmbed(
			currentChannel,
			&discordgo.MessageEmbed{
				Title: embedContent.Title,
				Author: &discordgo.MessageEmbedAuthor{
					Name:    embedContent.Author,
					URL:     embedContent.AuthorURL,
					IconURL: embedContent.AuthorImage,
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text: embedContent.Footer,
				},
				Description: strings.Join(embedContent.Content, " "),
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: embedContent.Thumbnail,
				},
			},
		); err != nil {
			log.Println("Failed to send embed!", err)
			return
		}

	default:
		if m, err := d.ChannelMessageSend(currentChannel, content); err != nil {
			log.Println(err)

		} else {
			latestSentID = m.ID
		}
	}
}
