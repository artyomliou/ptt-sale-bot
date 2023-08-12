package core

import (
	"context"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/tj/assert"
	"golang.org/x/net/html"
)

func TestParser(t *testing.T) {
	// Load test data
	data, err := os.ReadFile("./testdata/dc_sale_index.html")
	if err != nil {
		t.Error(err)
	}

	reader := strings.NewReader(string(data))

	htmlNode, err := html.Parse(reader)
	if err != nil {
		t.Error(err)
	}

	// Prepare parser
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	p := NewParser(ctx)
	p.SetInput(make(chan *Page))
	p.Run()

	p.Input <- &Page{
		HtmlNode: htmlNode,
	}
	articles := <-p.Output
	for _, article := range articles {
		log.Printf("{Title: \"%v\", Link: \"%v\"},", article.Title, article.Link)
	}
	assert.GreaterOrEqual(t, len(articles), 16)
}
