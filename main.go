package main

import (
	"flag"
	"log"
	"strconv"

	"github.com/RumbleFrog/discordgo"

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
)

func main() {
	token := flag.String("t", "", "Discord token (1)")

	flag.Parse()

	k, err := keyring.Get(AppName, "token")
	if err != nil {
		if err != keyring.ErrNotFound {
			log.Fatalln(err.Error())
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

	sciter.SetOption(
		sciter.SCITER_SET_SCRIPT_RUNTIME_FEATURES,
		sciter.ALLOW_FILE_IO|
			sciter.ALLOW_SOCKET_IO|
			sciter.ALLOW_EVAL|
			sciter.ALLOW_SYSINFO,
	)

	w, err = window.New(sciter.DefaultWindowCreateFlag, &sciter.Rect{
		Left:   0,
		Top:    0,
		Right:  800,
		Bottom: 600,
	})

	if err != nil {
		log.Fatal(err)
	}

	w.LoadFile("src/index.html")

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

		go loadMsgs(chid)

		return nil
	})

	w.DefineFunction("sendmessage", func(args ...*sciter.Value) *sciter.Value {
		if len(args) < 1 {
			return nil
		}

		if currentChannel < 1 {
			return nil
		}

		go func() {
			if _, err := d.ChannelMessageSend(currentChannel, args[0].String()); err != nil {
				log.Println(err)
			}
		}()

		return nil
	})

	w.SetTitle("tdeo is homo")

	d, err = discordgo.New(Token)
	if err != nil {
		panic(err)
	}

	d.AddHandler(onReady)
	d.AddHandler(messageCreate)

	if err := d.Open(); err != nil {
		panic(err)
	}

	defer d.Close()

	// go WatchCSS()

	w.Show()
	w.Run()
}
