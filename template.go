package main

import (
	"bytes"
	"html/template"
	"log"

	"github.com/RumbleFrog/discordgo"
)

var (
	// Data: messageTemplateData
	messageTemplate *template.Template

	// Data: messageContentTemplateData
	messageContentTemplate *template.Template

	// Data: []guildsTemplateData
	guildTemplate *template.Template

	// Data: []categoriesTemplateData
	channelsTemplate *template.Template
)

type messageTemplateData struct {
	ID         int64
	Timestamp  string
	AuthorID   int64
	AvatarName string

	DisplayName string
	NameColor   string

	Blocked bool
	Edited  bool

	Message template.HTML
}

type messageContentTemplateData struct {
	Content template.HTML

	Attachments []messageAttachmentTemplateData
	Embeds      []messageEmbedTemplateData
}

type MediaType string

const (
	TypeFile MediaType = "file"

	TypeAudio MediaType = "audio"

	TypeVideo MediaType = "video"

	TypeImage MediaType = "image"
)

type messageAttachmentTemplateData struct {
	MediaType MediaType
	Name      string
	URL       string
	Original  string
	Size      string
}

type messageEmbedTemplateData struct {
	PillColor string

	Title    string
	TitleURL string

	Author     string
	AuthorURL  string
	AuthorIcon string

	ThumbnailURL      string
	ThumbnailOriginal string
	ThumbnailWidth    int
	ThumbnailHeight   int

	Description template.HTML

	Fields []*discordgo.MessageEmbedField

	// Footer already includes timestamp
	Footer     string
	FooterIcon string

	ImageURL      string
	ImageOriginal string
	ImageWidth    int
	ImageHeight   int

	VideoURL      string
	VideoOriginal string

	// Todo: find out what Provider is
}

type guildsTemplateData struct {
	ID   int64
	URL  string
	Name string
}

type categoriesTemplateData struct {
	CategoryID   int64
	CategoryName string
	Channels     []channelsTemplateData
}

type channelsTemplateData struct {
	ID    int64
	Name  string
	NSFW  bool
	Muted bool
}

// RenderToString converts the appropriate option to template string
// using reflection magic
func RenderToString(data interface{}) string {
	var (
		b   bytes.Buffer
		err error
	)

	switch data.(type) {
	case messageTemplateData:
		err = messageTemplate.Execute(&b, data)
	case messageContentTemplateData:
		err = messageContentTemplate.Execute(&b, data)
	case []guildsTemplateData:
		err = guildTemplate.Execute(&b, data)
	case []categoriesTemplateData:
		err = channelsTemplate.Execute(&b, data)
	default:
		return ""
	}

	if err != nil {
		log.Println(err)
		return ""
	}

	return b.String()
}
