package main

import (
	"log"
	"testing"
)

func TestMDtoHTML(t *testing.T) {
	log.Println(MDtoHTML(`"this is a test!" says the guy's phone` + "```bash\n" + `
echo test
` + "```"))
}
