package articles_parser

import (
	"habr-bot/source/models"
	"habr-bot/source/utils"
	"strings"
	"testing"
)

func TestService_ParseCategorySuccess(t *testing.T) {
	service := New("https://habr.com", "https://habr.com/ru/hub", "post__title_link")
	category := "go:top10"
	articles := make(chan models.Article)
	articlesProcessing := models.ProcessingChannels{
		Done:  make(chan bool),
		Error: make(chan error),
	}

	go service.ParseCategory(category, articles, articlesProcessing)

	for {
		select {
		case article := <-articles:
			utils.AssertTrue(category == article.Category, t)
			utils.AssertTrue(strings.Contains(article.Link, "habr.com"), t)
			utils.AssertTrue("" != article.Title, t)
		case <-articlesProcessing.Done:
			return
		case err := <-articlesProcessing.Error:
			t.Logf("Error: %v\n", err)
			t.Fail()
		}
	}
}

func TestService_ParseCategoryError(t *testing.T) {
	service := New("https://habr.com", "bad-url", "")
	category := "go"
	articles := make(chan models.Article)
	articlesProcessing := models.ProcessingChannels{
		Done:  make(chan bool),
		Error: make(chan error),
	}

	go service.ParseCategory(category, articles, articlesProcessing)

	for {
		select {
		case <-articles:
			t.Log("Should not receive article!")
			t.Fail()
			return
		case <-articlesProcessing.Done:
			t.Log("Should not done!")
			t.Fail()
			return
		case err := <-articlesProcessing.Error:
			utils.AssertNotNil(err, t)
			return
		}
	}
}
