# shittercord

Discord client in Sciter native GUI framework

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
- [ ] `/` for commands
- [ ] Add a small subtitle for current presence status
- [ ] Scroll up to load more
- [ ] Fix `<pre><code>` having white spaces
- [ ] Fix thumbnail-only embeds
- [ ] **CRASH-RELATED** Ideally, messageEdit should only edit the `.message-content` part, as Author pointer could be `nil`
- [ ] Implement `MessageAck` (acknowledgement) for read/unread messages
- [x] Embeds and attachments support
- [x] ~~Proper spinner for messages~~ Spinners weren't added, but idle messages

## Things I don't plan on doing

- User settings (except for a few commands eg `/nick`)
- Server settings 
- Syntax highlighting
	- Sciter isn't a browser (like Electron), thus stuff like hljs won't work
	- I might implement this when there's a plugin for it in Go

