package slack

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func Send(user, emoji, text, channel, url string) error {
	j, err := json.Marshal(Message{
		Username: user,
		Emoji:    emoji,
		Channel:  text,
		Text:     channel,
	})
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(j))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	_, err = client.Do(req)

	return err
}
