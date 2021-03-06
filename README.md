# Discord archive and delete

Archives all posts and file attachments in a given discord server or channel

### Building

```
go get -u github.com/BurntSushi/toml
go build discord-archive-and-delete.go
```

### Example configuration

Place a file called `config.toml` in the same directory as discord-archive-and-delete with the following values

```toml
# The digit ID of the server you wish to archive
server = "..."

# The digit ID of the channel you wish to archive
# If you want to archive every channel, set this to the main channel ID
channel = "..."

# Your discord account's digit ID
user_id = "..."

# Your user token
user_token = "..."

# Your browser cookie
user_cookie = "__cfduid=..."

# Should we backup any images or files you uploaded? (true/false)
save_attachments = false

# Should we delete messages after archiving? (true/false)
delete_messages = false

# Should we only go through one channel? (true/false)
only_channel = true
```

### Configuring

1. Enable ["Developer Mode"](https://support.discordapp.com/hc/article_attachments/115002742731/mceclip0.png)
2. Open the Developer Console (`Ctrl + Shift + I`)
3. In the console go to Application tab, then -> Storage -> [Local Storage](https://i.imgur.com/XtnjUB3.png)
4. Refresh Discord (`Ctrl + R`)
5. Copy the value for [`token`](https://i.imgur.com/RcQHOWJ.png)
6. Paste as the value for `user_token` in config.toml
7. Go to [Cookies](https://i.imgur.com/aLiPDGw.png) in the same Application tab
8. Copy the value for [`__cfduid`](https://i.imgur.com/ZM0t2hE.png)
9. Paste this as the value for `user_cookie` in config.toml after `__cfduid=`
10. Right click the icon of the desired server and click [Copy ID](https://i.imgur.com/ocXjHqq.png)
11. Paste this as the value for `server` in config.toml
12. Right click the channel you want to archive and click [Copy ID](https://i.imgur.com/wCsqtqD.png) (if you wish to archive the whole server, select any channel)
13. Paste this as the value for `channel` in config.toml
14. Right click on your avatar in one of your posts and click [Copy ID](https://i.imgur.com/lTtcmoJ.png)
15. Paste this as the value for `user_id` in config.toml
16. If you want to save the files/attachments you've uploaded, set `save_attachments` to `true`
17. If you want to delete all your messages, set `delete_messages` to `true`
18. If you want to only archive one channel, set `only_channel` to `true`. If you want to archive every channel, set it to `false`

### Usage

**MacOS or GNU/Linux**

Open a terminal window and `cd` into the directory containing `discord-archive-and-delete` and `config.toml`

```
chmod u+x discord-archive-and-delete
./discord-archive-and-delete
```

**Windows**

Open a PowerShell or Command Prompt window then `cd` into the folder containing `discord-archive-and-delete.exe` and `config.toml`

```
.\discord-archive-and-delete.exe
```
