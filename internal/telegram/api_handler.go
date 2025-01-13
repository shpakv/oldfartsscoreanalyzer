package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type API interface {
	SendMessage(message string) error
	GetAdministrators() ([]User, error)
	GetChatMember(userId string) (*User, error)
}

type (
	User struct {
		User   user   `json:"user"`
		Status string `json:"status"`
	}

	user struct {
		Id        int    `json:"id"`
		IsBot     bool   `json:"is_bot"`
		FirstName string `json:"first_name"`
		Username  string `json:"username"`
	}
)

const APIURLPrefix = "https://api.telegram.org/bot"
const (
	SendMessageAction       = "sendMessage"
	GetAdministratorsAction = "getChatAdministrators"
	GetMemberAction         = "getChatMember"
)

type DefaultAPIHandler struct {
	Bot    *Bot
	ChatId string
}

func NewDefaultAPIHandler(bot *Bot, chatId string) *DefaultAPIHandler {
	return &DefaultAPIHandler{Bot: bot, ChatId: chatId}
}

func (d *DefaultAPIHandler) SendMessage(message string) error {
	url := fmt.Sprintf("%s%s/%s", APIURLPrefix, d.Bot.Token, SendMessageAction)
	body := map[string]string{
		"chat_id":    d.ChatId,
		"text":       message,
		"parse_mode": "Markdown",
	}
	bodyJSON, _ := json.Marshal(body)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(bodyJSON))
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (d *DefaultAPIHandler) GetAdministrators() ([]User, error) {
	url := fmt.Sprintf("%s%s/%s?chat_id=%s", APIURLPrefix, d.Bot.Token, GetAdministratorsAction, d.ChatId)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get administrators: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed with status code: %d", resp.StatusCode)
	}

	var result struct {
		OK     bool   `json:"ok"`
		Result []User `json:"result"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Result, nil
}

func (d *DefaultAPIHandler) GetChatMember(userId string) (*User, error) {
	url := fmt.Sprintf("%s%s/%s?chat_id=%s&user_id=%s", APIURLPrefix, d.Bot.Token, GetMemberAction, d.ChatId, userId)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get administrators: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed with status code: %d", resp.StatusCode)
	}

	var result struct {
		OK     bool   `json:"ok"`
		Result []User `json:"result"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	if len(result.Result) == 0 {
		return nil, fmt.Errorf("member with id %s not found", userId)
	}

	return &result.Result[0], nil
}
