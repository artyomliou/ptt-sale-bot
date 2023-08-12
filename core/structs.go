package core

import (
	"regexp"

	"golang.org/x/net/html"
)

type Page struct {
	HtmlNode *html.Node
}

type Article struct {
	Link  string
	Title string
}

type InterestTopic struct {
	Name             string
	CompiledPatterns []*regexp.Regexp
}

type TelegramSendMessagepayload struct {
	ChatId                int    `json:"chat_id"`
	Text                  string `json:"text"`
	ParseMode             string `json:"parse_mode"`
	DisableWebPagePreview bool   `json:"disable_web_page_preview"`
}
