package core

import (
	"context"
	"log"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
)

type filterer struct {
	ctx              context.Context
	Input            chan []*Article
	Output           chan []*Article
	interestTopics   []InterestTopic
	deduplicateLinks mapset.Set[string]
}

func NewFilterer(ctx context.Context) *filterer {
	return &filterer{
		ctx:              ctx,
		Input:            nil,
		Output:           make(chan []*Article),
		interestTopics:   []InterestTopic{},
		deduplicateLinks: mapset.NewSetWithSize[string](100),
	}
}

func (f *filterer) SetInput(input chan []*Article) {
	f.Input = input
}

func (f *filterer) AddInterestTopic(topic InterestTopic) {
	f.interestTopics = append(f.interestTopics, topic)
}

func (f *filterer) Run() {
	go func() {
		for {
			select {
			case articles := <-f.Input:
				f.handle(articles)
			case <-f.ctx.Done():
				log.Println("[filterer] stop")
				return
			}
		}
	}()
}

func (f *filterer) handle(articles []*Article) {
	matchedArticles := []*Article{}

	for _, article := range articles {

		if f.skipIfArticleHasBeenFiltered(article) {
			continue
		}
		if f.skipIfTitleContainKeyword(article) {
			continue
		}

		// for each topic, use multiple pattern for filtering
		for _, topic := range f.interestTopics {
			matchedCount := 0
			for _, pattern := range topic.CompiledPatterns {
				if !pattern.MatchString(article.Title) {
					break // no further evaluating with this topic
				}
				matchedCount++
			}
			if matchedCount == len(topic.CompiledPatterns) {
				matchedArticles = append(matchedArticles, article)
			}
		}
	}

	// insert matchedArticle into set
	for _, article := range matchedArticles {
		f.deduplicateLinks.Add(article.Link)
	}

	f.Output <- matchedArticles
}

func (f *filterer) skipIfArticleHasBeenFiltered(article *Article) bool {
	return f.deduplicateLinks.Contains(article.Link)
}

func (f *filterer) skipIfTitleContainKeyword(article *Article) bool {
	skipKeywords := []string{
		"公告",
		"黑名",
		"遭竊",
	}
	for _, skipKeyword := range skipKeywords {
		if strings.Contains(article.Title, skipKeyword) {
			return true
		}
	}
	return false
}
