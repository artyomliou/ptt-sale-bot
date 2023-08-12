package core

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

func NewCrawler(ctx context.Context) *crawler {
	return &crawler{
		ctx:      ctx,
		Interval: 5 * time.Minute,
		targets:  []string{},
		Output:   make(chan *Page, 1),
	}
}

type crawler struct {
	ctx      context.Context
	Interval time.Duration
	targets  []string
	Output   chan *Page
}

func (c *crawler) AddTarget(url string) {
	c.targets = append(c.targets, url)
}

func (c *crawler) Run() {
	go func() {
		ticker := time.NewTicker(c.Interval)
		for {
			select {
			case <-ticker.C:
				c.execute()
			case <-c.ctx.Done():
				return
			}
		}
	}()
}

func (c *crawler) execute() {
	log.Println("[crawler] executing")

	if len(c.targets) == 0 {
		log.Println("[crawler] no target, skip")
		return
	}

	// TODO request in parellel
	for _, target := range c.targets {
		log.Printf("[crawler] fetch %s", target)
		req, err := http.NewRequestWithContext(c.ctx, http.MethodGet, target, nil)
		if err != nil {
			log.Printf("failed to init a request: %v", err)
			continue
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("failed to request a target: %v", err)
			continue
		}
		if !strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
			log.Printf("Response of target is not HTML: %v", err)
			continue
		}

		htmlNode, err := html.Parse(resp.Body)
		if err != nil {
			log.Printf("failed to parse response as HTML node: %v", err)
			continue
		}

		c.Output <- &Page{
			HtmlNode: htmlNode,
		}
	}
}
