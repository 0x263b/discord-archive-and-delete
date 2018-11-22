package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/BurntSushi/toml"
)

const (
	searchURL        = "https://discordapp.com/api/v6/guilds/%s/messages/search?author_id=%s&offset=%d"
	searchChannelURL = "https://discordapp.com/api/v6/guilds/%s/messages/search?author_id=%s&channel_id=%s&offset=%d"
)

type Config struct {
	Guild          string `toml:"server"`
	Channel        string `toml:"channel"`
	UserID         string `toml:"user_id"`
	UserToken      string `toml:"user_token"`
	UserCookie     string `toml:"user_cookie"`
	Referer        string `toml:"referer"`
	SaveImages     bool   `toml:"save_attachments"`
	DeleteMessages bool   `toml:"delete_messages"`
	OnlyChannel    bool   `toml:"only_channel"`
}

type SearchResults struct {
	TotalResults int         `json:"total_results"`
	AnalyticsID  string      `json:"analytics_id"`
	Messages     [][]Message `json:"messages"`
}

type Message struct {
	Attachments     []Attachment `json:"attachments"`
	Tts             bool         `json:"tts"`
	Embeds          []Embed      `json:"embeds"`
	Timestamp       time.Time    `json:"timestamp"`
	MentionEveryone bool         `json:"mention_everyone"`
	ID              string       `json:"id"`
	Pinned          bool         `json:"pinned"`
	EditedTimestamp time.Time    `json:"edited_timestamp"`
	Author          Author       `json:"author"`
	MentionRoles    []string     `json:"mention_roles"`
	Content         string       `json:"content"`
	ChannelID       string       `json:"channel_id"`
	Mentions        []Author     `json:"mentions"`
	Type            int          `json:"type"`
	Hit             bool         `json:"hit,omitempty"`
}

type Attachment struct {
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url"`
	Filename string `json:"filename"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	ID       string `json:"id"`
	Size     int    `json:"size"`
}

type Embed struct {
	Description string `json:"description"`
	Author      struct {
		URL          string `json:"url"`
		IconURL      string `json:"icon_url"`
		ProxyIconURL string `json:"proxy_icon_url"`
		Name         string `json:"name"`
	} `json:"author"`
	URL    string `json:"url"`
	Fields []struct {
		Inline bool   `json:"inline"`
		Name   string `json:"name"`
		Value  string `json:"value"`
	} `json:"fields"`
	Footer struct {
		Text         string `json:"text"`
		ProxyIconURL string `json:"proxy_icon_url"`
		IconURL      string `json:"icon_url"`
	} `json:"footer"`
	Video struct {
		URL    string `json:"url"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	} `json:"video"`
	Type      string `json:"type"`
	Thumbnail struct {
		URL      string `json:"url"`
		Width    int    `json:"width"`
		ProxyURL string `json:"proxy_url"`
		Height   int    `json:"height"`
	} `json:"thumbnail"`
}

type Author struct {
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Bot           bool   `json:"bot"`
	ID            string `json:"id"`
	Avatar        string `json:"avatar"`
}

func CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}

func DownloadFile(filepath string, url string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func main() {

	// Load config
	var config Config
	if _, err := toml.DecodeFile("./config.toml", &config); err != nil {
		fmt.Println(err)
		return
	}

	config.Referer = fmt.Sprintf("https://discordapp.com/channels/%s/%s", config.Guild, config.Channel)

	// Loop through the search results page
	var more bool = true
	var offset int32 = 0

	var messages []Message

	fmt.Println("Searching for messages. This may take several minutes.")
	for more {
		var url string
		if config.OnlyChannel == true {
			url = fmt.Sprintf(searchChannelURL, config.Guild, config.UserID, config.Channel, offset)
		} else {
			url = fmt.Sprintf(searchURL, config.Guild, config.UserID, offset)
		}

		client := &http.Client{}
		request, _ := http.NewRequest("GET", url, nil)
		request.Header.Set("Authorization", config.UserToken)
		request.Header.Set("Cookie", config.UserCookie)
		request.Header.Set("Referer", config.Referer)

		response, _ := client.Do(request)
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)

		var results SearchResults
		json.Unmarshal(body, &results)

		for _, msgs := range results.Messages {
			for _, msg := range msgs {
				if msg.Author.ID == config.UserID {
					messages = append(messages, msg)
				}
			}
		}

		offset += 25

		if len(results.Messages) == 0 {
			more = false
		}

		// Sleep to avoid rate limit
		time.Sleep(1 * time.Second)
	}

	// Unique slice
	var unqiueMessages []Message
	var keys []string
	for _, msg := range messages {
		if stringInSlice(msg.ID, keys) == false {
			unqiueMessages = append(unqiueMessages, msg)
		}
		keys = append(keys, msg.ID)
	}
	// Sort the messages by time
	sort.Slice(unqiueMessages, func(i, j int) bool { return unqiueMessages[i].Timestamp.Before(unqiueMessages[j].Timestamp) })

	// Create directory for logs
	n := time.Now().Unix()
	d := fmt.Sprintf("log_%d", n)

	CreateDirIfNotExist(d)
	if config.SaveImages == true {
		CreateDirIfNotExist(fmt.Sprintf("%s/attachments", d))
	}

	// Export logs to file
	exportData, _ := json.MarshalIndent(unqiueMessages, "", "  ")
	exportFile, _ := os.Create(fmt.Sprintf("%s/logs.json", d))
	defer exportFile.Close()

	exportFile.Write(exportData)
	exportFile.Close()
	fmt.Printf("Saved %d messages to %s\n", len(unqiueMessages), exportFile.Name())

	var saved int32 = 0
	var deleted int32 = 0

	if config.DeleteMessages == true {
		low := time.Duration((int64(len(unqiueMessages)/2) * int64(time.Second)))
		high := time.Duration((int64(len(unqiueMessages)) * int64(time.Second)))
		fmt.Printf("Deleting messages. This will take approximately %s to %s.\n", low, high)
	}

	// loop through the messages
	for _, msg := range unqiueMessages {
		// Download attachments
		if config.SaveImages == true && len(msg.Attachments) > 0 {
			for _, file := range msg.Attachments {
				// Download file
				f := fmt.Sprintf("%d_%s", msg.Timestamp.Unix(), file.Filename)

				err := DownloadFile(fmt.Sprintf("%s/attachments/%s", d, f), file.URL)
				if err != nil {
					panic(err)
				}

				saved += 1
			}
		}

		// Delete message
		if config.DeleteMessages == true {
			client := &http.Client{}
			url := fmt.Sprintf("https://discordapp.com/api/channels/%s/messages/%s", msg.ChannelID, msg.ID)
			request, _ := http.NewRequest("DELETE", url, nil)
			request.Header.Set("Authorization", config.UserToken)
			request.Header.Set("Cookie", config.UserCookie)
			request.Header.Set("Referer", config.Referer)

			response, _ := client.Do(request)
			defer response.Body.Close()
			if response.StatusCode != 204 {
				fmt.Printf("Error deleting message: %s (%d)\n", msg.ID, response.StatusCode)
			}

			deleted += 1

			// Sleep to avoid rate limit
			time.Sleep(500 * time.Millisecond)
		}
	}

	// Done
	fmt.Printf("Deleted %d messages\n", deleted)
	fmt.Printf("Saved %d attachments\n", saved)
}
