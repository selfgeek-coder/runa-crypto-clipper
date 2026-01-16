package telegram

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

type telegramRequest struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode,omitempty"`
}

type telegramResponse struct {
	Ok bool `json:"ok"`
}

var client = &http.Client{
	Timeout: 8 * time.Second,
}

// send msg to telegram 
func SendLog(message, chatID, botToken string) bool {
    reqBody := telegramRequest{
        ChatID:    chatID,
        Text:      message,
        ParseMode: "HTML",
    }

    data, err := json.Marshal(reqBody)
    if err != nil {
        return false
    }

    url := "https://api.telegram.org/bot" + botToken + "/sendMessage"

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
    if err != nil {
        return false
    }

    req.Header.Set("Content-Type", "application/json")

    resp, err := client.Do(req)
    if err != nil {
        return false
    }
    defer resp.Body.Close()

    var tgResp telegramResponse
    if err := json.NewDecoder(resp.Body).Decode(&tgResp); err != nil {
        return false
    }

    return tgResp.Ok
}