package main

import (
	"html"
	"log"

	"github.com/sciter-sdk/go-sciter"
)

// GetElementByCSS gets element by CSS identifier (.css, etc)
func GetElementByCSS(css string) *sciter.Element {
	root, err := w.GetRootElement()
	if err != nil {
		return nil
	}

	elem, err := root.SelectUnique(css)
	if err != nil {
		return nil
	}

	return elem
}

// SetHTML sets the innerHTML for the element
func SetHTML(elem *sciter.Element, html string) {
	if elem == nil {
		log.Println("elem is nil")
	}

	elem.SetHtml(html, sciter.SIH_REPLACE_CONTENT)
	elem.Update(false) // screen sometimes flashes white with false
	// we'll see if this fixes it
}

// SetText uses SetHTML with html EscapeString instead of
// Sciter's SetText function
func SetText(elem *sciter.Element, text string) {
	SetHTML(elem, html.EscapeString(text))
}
