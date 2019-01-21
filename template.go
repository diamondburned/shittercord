package main

import (
	"bytes"
	"html/template"
	"log"
)

var (
	// Data: messageTemplateData
	messageTemplate = template.Must(template.ParseFiles("templates/message.html"))

	// Data: []guildsTemplateData
	guildTemplate = template.Must(template.ParseFiles("templates/guilds.html"))

	// Data: []categoriesTemplateData
	channelsTemplate = template.Must(template.ParseFiles("templates/channels.html"))
)

type messageTemplateData struct {
	ID         int64
	AuthorID   int64
	AvatarName string

	DisplayName string
	NameColor   string

	Content template.HTML

	// Todo: embed
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
