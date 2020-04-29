package articles_parser

import (
	"fmt"
	"golang.org/x/net/html"
	"interfaces"
	"io"
	"models"
	"net/http"
	"strings"
)

type Service struct {
	habrHubUrl           string
	articleLinkClassName string
	filter               string
}

func New(url, className, filter string) interfaces.ArticlesParserService {
	return Service{url, className, filter}
}

func (s Service) ParseCategory(
	category string,
	articles chan<- models.Article,
	articlesProcessing models.ProcessingChannels,
) {
	responseBody, err := s.getCategoryContent(category)
	if err != nil {
		articlesProcessing.Error <- err
		return
	}

	tokenizer := html.NewTokenizer(responseBody)
	for {
		tokenType := tokenizer.Next()

		switch {
		case tokenType == html.ErrorToken:
			articlesProcessing.Done <- true
			return
		case tokenType == html.StartTagToken:
			token := tokenizer.Token()

			if s.isArticleLink(token) {
				tokenizer.Next()
				article := s.prepareArticle(category, token)
				article.Title = string(tokenizer.Text())

				articles <- article
				continue
			}
		}
	}
}

func (s Service) getCategoryContent(category string) (io.ReadCloser, error) {
	response, err := http.Get(s.getURL(category))
	if err != nil {
		return nil, err
	}

	return response.Body, nil
}

func (s Service) getURL(category string) string {
	return fmt.Sprintf("%s/%s/%s", s.habrHubUrl, category, s.filter)
}

func (s Service) isArticleLink(token html.Token) bool {
	var (
		isAnchor   = token.Data == "a"
		isPostLink bool
	)
	for _, attr := range token.Attr {
		if attr.Key == "class" {
			isPostLink = strings.Contains(attr.Val, s.articleLinkClassName)
			break
		}
	}

	return isAnchor && isPostLink
}

func (s Service) prepareArticle(category string, token html.Token) models.Article {
	article := models.Article{
		Category: category,
	}
	for _, attr := range token.Attr {
		if attr.Key == "href" {
			article.Link = attr.Val
		}
	}

	return article
}
