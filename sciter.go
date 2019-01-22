package main

import (
	"fmt"
	"html"
	"log"
	"runtime"
	"strings"
	"time"

	"github.com/sciter-sdk/go-sciter"
)

// GetElementByCSS gets element by CSS identifier (.css, etc)
func GetElementByCSS(css string) (elem *sciter.Element) {
	for i := 1; i < 3; i++ {
		root, err := w.GetRootElement()
		if err == nil {
			elem, err = root.SelectUnique(css)
			if err == nil {
				break
			}

			log.Println(css, err)
			goto Fail
		}

		log.Println(err)

	Fail:
		time.Sleep(time.Millisecond * 200)
	}

	return
}

// SetHTML sets the innerHTML for the element
func SetHTML(elem *sciter.Element, html string) {
	if elem == nil {
		log.Println("elem is nil")
		return
	}

	elem.SetHtml(html, sciter.SIH_REPLACE_CONTENT)
	elem.Update(true) // screen sometimes flashes white with false
	// we'll see if this fixes it
}

// SetText uses SetHTML with html EscapeString instead of
// Sciter's SetText function
func SetText(elem *sciter.Element, text string) {
	SetHTML(elem, html.EscapeString(text))
}

// WarnDialog .
func WarnDialog(i ...interface{}) {
	var (
		printed []string
		line    int
		fn      string
		unknown []interface{}
	)

	_, fn, line, _ = runtime.Caller(1)

	for _, thing := range i {
		switch thing.(type) {
		case string:
			printed = append(
				printed,
				thing.(string),
			)
		case error:
			printed = append(
				printed,
				thing.(error).Error(),
			)
		default:
			unknown = append(
				unknown,
				thing,
			)
		}
	}

	if len(unknown) > 1 {
		log.Printf("%s:%d", fn, line)
		log.Println(unknown...)
	}

	log.Println(
		fmt.Sprintf(
			"%s:%d: %s",
			fn, line,
			strings.Join(printed, "\n"),
		),
	)
}
