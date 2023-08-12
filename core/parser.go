package core

import (
	"context"
	"log"

	"github.com/PuerkitoBio/goquery"
)

type parser struct {
	ctx    context.Context
	Input  chan *Page
	Output chan []*Article
}

func NewParser(ctx context.Context) *parser {
	return &parser{
		ctx:    ctx,
		Input:  nil,
		Output: make(chan []*Article, 1),
	}
}

func (p *parser) SetInput(input chan *Page) {
	p.Input = input
}

func (p *parser) Run() {
	go func() {
		for {
			select {
			case page := <-p.Input:
				p.execute(page)
			case <-p.ctx.Done():
				return
			}
		}
	}()
}

func (p *parser) execute(page *Page) {
	articles := []*Article{}

	doc := goquery.NewDocumentFromNode(page.HtmlNode)
	doc.Find(".r-list-container .r-ent .title").Each(func(i int, s *goquery.Selection) {
		link, ok := s.Find("a").Attr("href")
		if !ok {
			log.Println("failed to get a:href")
			return
		}
		title := s.Find("a").Text()
		articles = append(articles, &Article{
			Link:  "https://www.ptt.cc" + link,
			Title: title,
		})
	})

	p.Output <- articles
}
