package main

import "log"

func sendMessage(content string) {

	if _, err := d.ChannelMessageSend(currentChannel, content); err != nil {
		log.Println(err)
	}
}
