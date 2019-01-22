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

