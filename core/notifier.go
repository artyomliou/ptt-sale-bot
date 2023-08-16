package core

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
)

type notifier struct {
	ctx               context.Context
	telegramHost      string
	telegramBotKey    string
	telegramChannelId int
	Input             chan []*Article
	Output            chan bool
}

func NewNotifier(ctx context.Context) *notifier {
	return &notifier{
		ctx:    ctx,
		Input:  nil,
		Output: make(chan bool),
	}
}

func (n *notifier) SetInput(input chan []*Article) {
	n.Input = input
}

func (n *notifier) Run() {
	flag.StringVar(&n.telegramHost, "tg-host", "https://api.telegram.org", "Telegram hostname")
	flag.StringVar(&n.telegramBotKey, "tg-bot-key", "", "Telegram bot key")
	flag.IntVar(&n.telegramChannelId, "tg-channel-id", 0, "Telegram channel ID")
	flag.Parse()

	go func() {
		for {
			select {
			case articles := <-n.Input:
				n.handle(articles)
			case <-n.ctx.Done():
				log.Println("[notifier] stop")
				return
			}
		}
	}()
}

func (n *notifier) handle(articles []*Article) {

	if len(articles) == 0 {
		return
	}
	log.Printf("[Notifier] sending %d articles to telegram channel", len(articles))

	req, err := n.newTelegramSendMessageRequest(n.articleToMessageText(articles))
	if err != nil {
		log.Printf("[Notifier] failed when preparing request: %s", err.Error())
		n.Output <- false
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("[Notifier] failed when sending request: %s", err.Error())
		n.Output <- false
		return
	}
	if resp.StatusCode != 200 {
		respBodyString, _ := io.ReadAll(resp.Body)
		log.Printf("[Notifier] response code is not 200 but %d and got response: %s", resp.StatusCode, respBodyString)
		n.Output <- false
		return
	}

	n.Output <- true
}

func (n *notifier) articleToMessageText(articles []*Article) string {
	text := ""
	for _, article := range articles {
		text += fmt.Sprintf("<a href=\"%s\">%s</a>\n", article.Link, article.Title)
	}
	return text
}

func (n *notifier) newTelegramSendMessageRequest(text string) (*http.Request, error) {
	var buf bytes.Buffer
	body := TelegramSendMessagepayload{
		ChatId:                int(n.telegramChannelId),
		Text:                  text,
		ParseMode:             "HTML",
		DisableWebPagePreview: true,
	}
	if err := json.NewEncoder(&buf).Encode(&body); err != nil {
		return nil, fmt.Errorf("[Notifier] failed to do json-masrshalling: %s", err.Error())
	}

	url := fmt.Sprintf("%s/bot%s/sendMessage", n.telegramHost, n.telegramBotKey)
	req, err := http.NewRequestWithContext(n.ctx, http.MethodPost, url, &buf)
	if err != nil {
		return nil, fmt.Errorf("[Notifier] failed to make request: %s", err.Error())
	}

	req.Header.Set("Content-Type", "application/json")

	return req, nil
}
