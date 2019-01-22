# shittercord

Discord client in Sciter native GUI framework

![](https://cdn.discordapp.com/attachments/361920025051004939/537151107680698370/unknown.png)
![](https://media.discordapp.net/attachments/361920025051004939/537172219328069632/unknown.png)

## Installation

1. `git clone https://gitlab.com/diamondburned/shittercord`
2. `cd shittercord`
3. `go get`
4. `go run *.go -t <DISCORD TOKEN>` (if you have a Keyring Manager, you only need to pass in the token once)
5. Optional: `go build && ./shittercord`

## Todo

- [ ] Proper dialogs to warn if messages are sent or not
- [ ] Fix possible graphical glitches
- [ ] Min and max height in plaintext
- [ ] `/` for commands (untested, only `embed` and `replace(All)` for now)
- [ ] Add a small subtitle for current presence status
- [ ] Scroll up to load more
- [ ] Fix `<pre><code>` having white spaces
- [ ] Fix thumbnail-only embeds
- [ ] Implement `MessageAck` (acknowledgement) for read/unread messages
- [ ] Add in mini handlers to change the global variables eg. Settings
- [ ] Add in a method of guaranteeing message order, as parsing messages with code is considerably slower
- [ ] Hide channels that can't be seen
- [ ] A different symbol for NSFW channels
- [ ] Pack all HTML/CSS materials in the binary
- [ ] Better timestamping, remove the edited part from timestamps and use `.message-content:edited`
- [x] Added syntax highlighting
- [x] **CRASH-RELATED** Ideally, messageEdit should only edit the `.message-content` part, as Author pointer could be `nil`
- [x] Embeds and attachments support
- [x] ~~Proper spinner for messages~~ Spinners weren't added, but idle messages

## Things I don't plan on doing

- User settings (except for a few commands eg `/nick`)
- Server settings 
- Syntax highlighting
	- Sciter isn't a browser (like Electron), thus stuff like hljs won't work
	- I might implement this when there's a plugin for it in Go

## Documentation

### [Inline commands](https://gitlab.com/diamondburned/shittercord/blob/master/sendmessage.go)

#### `/embed`

- Usage: `/embed [OPTIONS] [DESCRIPTION]`
- Options: 
	- `-t|--title` - title
	- `-a|--author` - author's name
	- `--authorURL` - author's URL (hyperlink)
	- `--authorImage` - author's Image
	- `-f|--footer` - footer text
	- `-th|--thumbnail` - embed thumbnail, URL

#### `/replace|/replaceAll`

Replaces the latest message of yours, same as Discord's `s/string1/string2`

- Usage: `/replace[All] [string] [string to replace with]`

### CSS

#### Injecting

You can currently inject a CSS file with the `-css` flag, followed by the file path. Do note that this feature is untested, and you might need `!important`. It is also important to point out that Sciter's CSS is not like other browser's CSS.

##### Resources

- [Almost all CSS properties](https://sciter.com/docs/content/css/cssmap.html)
- [HTML's flex vs Sciter's flow](https://terrainformatica.com/w3/flex-layout/flex-vs-flexbox.htm)

#### Background color

```css
html {
	/* 
		This controls the entire window's background color,
		including chat messages, channel boxes, etc
	*/
	background-color: rgba(25, 25, 25, 0.75);

	/*
		This element is redundant on Linux, not sure what
		it does on Windows 
	*/
	window-frame: "transparent";
}
```

#### Fonts

```css
/* Formal fonts for the app */
body {
	/*
		This is the default font list for the UI. Worth pointing
		out that Sciter is a bit retarded with the font fallback,
		so things might be missing glyphs. You might also need an
		!important on this one.
	*/
	font-family: "Overpass", "Segoe UI", "Helvetica", "Source Sans Pro", "Noto Sans", sans-serif;
}

/* Monospaced fonts for code/fenced code */
pre, code {
	/* 
		Might also need !important
	*/
	font-family: "Noto Sans Mono", monospace;
}
```

#### Top bar color

```css
/* ID for topbar */
#topbar {
	/*
		The top bar's background color. It has alpha, as
		the default style is a transparent one.
	*/
	background-color: rgba(45, 45, 45, 0.45);

	/*
		Self-explanatory, it's the same as HTML's CSS
	*/
	box-shadow: 0px 0px 5px 0px rgba(0, 0, 0, 0.35);
}
```

#### Hiding blocked messages

```css
.message.blocked {
	/*
		The original properties for blocked messages:
			height: 2px;
			overflow: hidden;

		This hides the blocked messages in a barely
		distinguishable/clickable box cleanly.
	*/

	/* The actual code to hide blocked messages */
	display: none;
}
```