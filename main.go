package main

import (
	"flag"
	"html/template"
	"log"
	"strconv"
	"strings"

	"github.com/RumbleFrog/discordgo"
	packr "github.com/gobuffalo/packr/v2"
	"github.com/pkg/browser"

	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
	keyring "github.com/zalando/go-keyring"
)

const (
	// AppName is the application name for keyring storage
	AppName = "scitercord"
)

var (
	// Token contains Discord token
	Token string

	w *window.Window

	d *discordgo.Session

	tplBox = packr.New("Templates", "./templates")
	srcBox = packr.New("Sources", "./src")
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	var (
		messageTemplateFile, _        = tplBox.FindString("message.html")
		messageContentTemplateFile, _ = tplBox.FindString("message-content.html")
		guildTemplateFile, _          = tplBox.FindString("guilds.html")
		channelsTemplateFile, _       = tplBox.FindString("channels.html")
	)

	messageTemplate = template.Must(
		template.New("messageTemplate").Parse(messageTemplateFile),
	)

	messageContentTemplate = template.Must(
		template.New("messageContentTemplate").Parse(messageContentTemplateFile),
	)

	guildTemplate = template.Must(
		template.New("guildTemplate").Parse(guildTemplateFile),
	)

	channelsTemplate = template.Must(
		template.New("channelsTemplate").Parse(channelsTemplateFile),
	)
}

func main() {
	var (
		token       = flag.String("t", "", "Discord token (1)")
		debug       = flag.Bool("d", false, "Enables the debugging Websocket for Inspector (Sciter developers only!)")
		customCSS   = flag.String("css", "", "Location to custom CSS (ONLY Sciter CSS! https://sciter.com/docs/content/css/cssmap.html)")
		deleteToken = flag.Bool("deletetoken", false, "Deletes token")
	)

	flag.Parse()

	if *deleteToken {
		log.Fatalln(keyring.Delete(AppName, "token"))
	}

	k, err := keyring.Get(AppName, "token")
	if err != nil {
		if err != keyring.ErrNotFound {
			log.Println(err)
		}

		if *token == "" {
			log.Fatalln("No tokens provided! Use `-t'!")
		}

		Token = *token

		log.Println("Storing token inside keyring...")
		if err := keyring.Set(AppName, "token", Token); err != nil {
			log.Println("Failed to set keyring! Continuing anyway...", err.Error())
		}

	} else {
		Token = k
	}

	if *debug {
		sciter.SetOption(
			sciter.SCITER_SET_SCRIPT_RUNTIME_FEATURES,
			sciter.ALLOW_FILE_IO|
				sciter.ALLOW_SOCKET_IO|
				sciter.ALLOW_EVAL|
				sciter.ALLOW_SYSINFO,
		)
	}

	w, err = window.New(sciter.DefaultWindowCreateFlag, &sciter.Rect{
		Left:   0,
		Top:    0,
		Right:  800,
		Bottom: 600,
	})

	if err != nil {
		log.Fatal(err)
	}

	var (
		html, _ = srcBox.FindString("index.html")
		css, _  = srcBox.FindString("style.css")
	)

	w.LoadHtml(html, "")
	sciter.AppendMasterCSS(css)

	if *customCSS != "" {
		log.Println("Appending custom Sciter CSS...")
		log.Println("Status:", sciter.AppendMasterCSS(*customCSS))
	}

	w.DefineFunction("loadguild", func(args ...*sciter.Value) *sciter.Value {
		if len(args) < 1 {
			return nil
		}

		gid, err := strconv.ParseInt(args[0].String(), 10, 64)
		if err != nil {
			log.Println(err)
			return nil
		}

		go loadGuild(gid)

		return nil
	})

	w.DefineFunction("loadchannel", func(args ...*sciter.Value) *sciter.Value {
		if len(args) < 1 {
			return nil
		}

		chid, err := strconv.ParseInt(args[0].String(), 10, 64)
		if err != nil {
			log.Println(err)
			return nil
		}

		go func() {
			loadMsgs(chid)
		}()

		return nil
	})

	w.DefineFunction("sendmessage", func(args ...*sciter.Value) *sciter.Value {
		if len(args) < 1 {
			return nil
		}

		if currentChannel < 1 {
			return nil
		}

		go sendMessage(args[0].String())

		return nil
	})

	w.DefineFunction("openURL", func(args ...*sciter.Value) *sciter.Value {
		if len(args) < 1 {
			return nil
		}

		log.Println(args[0].String())

		if e := browser.OpenURL(args[0].String()); e != nil {
			log.Println(e)
		}

		return nil
	})

	buffer := make(chan string, 2048)

	go func(buffer chan string) {
		for input := range buffer {
			fields := strings.Fields(input)

			if len(fields) == 0 {
				continue
			}

			i := fields[len(fields)-1]

			switch {
			case strings.HasPrefix(i, ":"):
				if len(i) < 2 {
					GetElementByCSS(".autosuggestions").Clear()
					continue
				}

				handleEmojis(i)
			}
		}
	}(buffer)

	w.DefineFunction("fuzzy", func(args ...*sciter.Value) *sciter.Value {
		if len(args) < 1 {
			return nil
		}

		buffer <- args[0].String()

		return nil
	})

	w.SetTitle("tdeo is homo")

	d, err = discordgo.New(Token)
	if err != nil {
		panic(err)
	}

	// Called to load guild list
	d.AddHandler(onReady)

	// The following function drops when the Channel ID
	// doesn't match the global variable

	// Called when there's a new message
	d.AddHandler(messageCreate)

	// Called when there's a message edited
	d.AddHandler(messageUpdate)

	// Called when there's a message deleted
	d.AddHandler(messageDelete)

	// Called when messages are deleted in bulk
	// This is called when a bot/self-bot uses the API
	// endpoint specifically for deleting in bulk
	d.AddHandler(messageDeleteBulk)

	// These 2 following functions are currently deprecated
	// They're never called, and they don't have an use
	d.AddHandler(presenceUpdate)
	d.AddHandler(presencesReplace)

	// This function is called when the user updates their settings
	// That includes status changes (online, busy, etc)
	d.AddHandler(userSettingsUpdate)

	if err := d.Open(); err != nil {
		panic(err) // not panic when I find a way to do modal dialog
	}

	// Set a max message count to cache
	d.State.MaxMessageCount = 50

	defer d.Close()

	// go WatchCSS()

	w.Show()
	w.Run()
}
