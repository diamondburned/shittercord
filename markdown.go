package main

import (
	"io"
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	blackfriday "gopkg.in/russross/blackfriday.v2"
)

/*
	Credits: https://eddieflores.com/tech/blackfriday-chroma/
	Code has slight modification
*/

const (
	// Both are unused atm
	mdHTMLFlags = 0 |
		blackfriday.SkipHTML |
		blackfriday.SkipLinks

	mdExtensions = 0 |
		blackfriday.NoIntraEmphasis |
		blackfriday.FencedCode |
		blackfriday.Autolink |
		blackfriday.Strikethrough |
		blackfriday.HardLineBreak
)

var (
	// Renderer is the renderer to the syntax highlighter
	// TODO: Make this a flag
	Renderer = NewChromaRenderer("solarized-dark256")
)

// MDtoHTML converts Markdown to HTML
// with syntax highlighting using the Chroma
// renderer
func MDtoHTML(md string) string {
	return strings.TrimSpace(
		string(
			blackfriday.Run(
				[]byte(md),
				blackfriday.WithNoExtensions(),
				blackfriday.WithExtensions(mdExtensions),
				blackfriday.WithRenderer(Renderer),
			),
		),
	)
}

// NewChromaRenderer creates a new syntax highlight renderer
func NewChromaRenderer(theme string) *ChromaRenderer {
	return &ChromaRenderer{
		html: &blackfriday.HTMLRenderer{
			HTMLRendererParameters: blackfriday.HTMLRendererParameters{
				Flags: mdHTMLFlags,
			},
		},
		theme: theme,
	}
}

type ChromaRenderer struct {
	html  *blackfriday.HTMLRenderer
	theme string
}

func (r *ChromaRenderer) RenderNode(w io.Writer, node *blackfriday.Node,
	entering bool) blackfriday.WalkStatus {
	switch node.Type {
	case blackfriday.CodeBlock:
		var lexer chroma.Lexer

		lang := string(node.CodeBlockData.Info)
		if lang != "" {
			lexer = lexers.Get(lang)
		} else {
			lexer = lexers.Analyse(string(node.Literal))
		}

		if lexer == nil {
			lexer = lexers.Fallback
		}

		style := styles.Get(r.theme)
		if style == nil {
			style = styles.Fallback
		}

		iterator, err := lexer.Tokenise(nil, string(node.Literal))
		if err != nil {
			panic(err)
		}

		formatter := html.New()

		err = formatter.Format(w, style, iterator)
		if err != nil {
			panic(err)
		}

		return blackfriday.GoToNext
	}

	return r.html.RenderNode(w, node, entering)
}

func (r *ChromaRenderer) RenderHeader(w io.Writer, ast *blackfriday.Node) {}
func (r *ChromaRenderer) RenderFooter(w io.Writer, ast *blackfriday.Node) {}
