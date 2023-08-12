package core

import (
	"context"
	"flag"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/tj/assert"
)

func TestNotifier(t *testing.T) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	// create notifier first
	n := NewNotifier(ctx)
	n.SetInput(make(chan []*Article, 1))
	n.Run()

	// then, mock endpoint
	didSentToMockEndpoint := false

	router := gin.Default()
	router.Any("/*proxyPath", func(ctx *gin.Context) {
		didSentToMockEndpoint = true
		ctx.AbortWithStatus(200)
	})
	server := httptest.NewServer(router)
	flag.Set("tg-host", server.URL)
	flag.Set("tg-bot-key", "12312312323")
	flag.Set("tg-channel-id", "0")

	// test
	n.Input <- []*Article{
		{
			Title: "Testing title",
			Link:  "https://www.ptt.cc/bbs/DC_SALE/M.1691245105.A.729.html",
		},
	}
	success := <-n.Output
	assert.Equal(t, true, success)
	assert.Equal(t, true, didSentToMockEndpoint)
}
