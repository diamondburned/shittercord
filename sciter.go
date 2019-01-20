package main

import (
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

func SetHTML(elem *sciter.Element, html string) {
	elem.SetHtml(html, sciter.SIH_REPLACE_CONTENT)
	elem.Update(false)
}
