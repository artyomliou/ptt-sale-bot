package core

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tj/assert"
)

func TestPttCrawler(t *testing.T) {
	// testing criteria
	targetWasRequested := false

	// setup testing server
	router := gin.Default()
	router.GET("/helloworld", func(ctx *gin.Context) {
		targetWasRequested = true
		ctx.Writer.Header().Add("Content-Type", "text/html")
		ctx.Writer.WriteString("<html></html>")
		ctx.AbortWithStatus(200)
	})
	server := httptest.NewServer(router)
	testingUrl := fmt.Sprintf("%s/helloworld", server.URL)

	// test crawler to request target in testing server
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	crawler := NewCrawler(ctx)
	crawler.Interval = 1 * time.Second
	crawler.AddTarget(testingUrl)
	crawler.Run()

	page := <-crawler.Output
	assert.Equal(t, true, targetWasRequested)
	assert.Equal(t, true, page != nil)
}
