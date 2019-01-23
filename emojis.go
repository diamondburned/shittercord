package main

import (
	"fmt"
	"strings"

	"github.com/RumbleFrog/discordgo"
	"github.com/sahilm/fuzzy"
)

// Emojis is a slice
type Emojis []*discordgo.Emoji

// String returns the fuzzy search part of the struct
func (e Emojis) String(i int) string {
	return e[i].Name
}

// Len returns the length of the Emojis slice
func (e Emojis) Len() int {
	return len(e)
}

// FuzzyRemoveDups fuzzy searches the list of emojis and returns the slice of results
func FuzzyRemoveDups(pattern string, emojis []*discordgo.Emoji) (nodups Emojis) {
	_e := Emojis(emojis)
	matches := fuzzy.FindFrom(pattern, _e)

	for i := 0; i < len(matches)-1; i++ {
		if strings.ToUpper(matches[i].Str) == strings.ToUpper(matches[i+1].Str) {
			continue
		}

		nodups = append(nodups, emojis[matches[i].Index])
	}

	return
}

/* HTML STRUCTURE FOR AUTOSUGGESTION ENTRY
[div.autosuggestions]
|	[img.emoji]
|	TEXT
|	[p.autosuggestion-subtitle]
*/

const EmojiSuggestionHTML = `
<div class="autosuggestion" id="%s">
	<img class="autosuggestion-emoji emoji" src="https://cdn.discordapp.com/emojis/%d.png" />
	%s
</div>
`

// ConstructAutocompletions builds the list of divs for emoji completion
// The maximum entry it will do is 8
func (e Emojis) ConstructAutocompletions() (html string) {
	for i := 0; i < 8 && i < len(e); i++ {
		emoji := e[i]

		html += fmt.Sprintf(
			EmojiSuggestionHTML,
			emoji.MessageFormat(), emoji.ID, emoji.Name,
		)
	}

	return
}
