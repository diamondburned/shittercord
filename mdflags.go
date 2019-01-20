package main

import "github.com/russross/blackfriday"

const (
	mdHTMLFlags = 0 |
		blackfriday.HTML_SKIP_STYLE |
		blackfriday.HTML_SKIP_IMAGES |
		blackfriday.HTML_SKIP_LINKS |
		blackfriday.HTML_SKIP_HTML

	mdExtensions = 0 |
		blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
		blackfriday.EXTENSION_FENCED_CODE |
		blackfriday.EXTENSION_AUTOLINK |
		blackfriday.EXTENSION_STRIKETHROUGH |
		blackfriday.EXTENSION_HARD_LINE_BREAK |
		blackfriday.EXTENSION_JOIN_LINES
)

// MDtoHTML .
func MDtoHTML(md string) string {
	return string(
		blackfriday.MarkdownOptions(
			[]byte(md),
			blackfriday.HtmlRenderer(mdHTMLFlags, "", ""),
			blackfriday.Options{
				Extensions: mdExtensions,
			},
		),
	)
}
