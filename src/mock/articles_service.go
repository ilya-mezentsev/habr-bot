package mock

import (
	"internal_errors"
	"models"
)

type RepositoryMock struct {
	Articles        []models.Article
	RepositoryError error
}

func (m *RepositoryMock) ReInit() {
	m.Articles = GetAllArticles()
	m.RepositoryError = nil
}

func (m *RepositoryMock) Save(
	articles chan models.Article,
	trySave chan bool,
	savingProcessing models.ProcessingChannels,
) {
	for {
		select {
		case article := <-articles:
			if article.Category == GetBadArticle().Category {
				savingProcessing.Error <- someRepositoryError
				return
			}

			m.Articles = append(m.Articles, article)
		case <-trySave:
			savingProcessing.Done <- true
		}
	}
}

func (m *RepositoryMock) GetByCategory(category string) ([]models.Article, error) {
	if category == GetBadArticle().Category {
		return nil, someRepositoryError
	}

	var articles []models.Article
	for _, article := range m.Articles {
		if article.Category == category {
			articles = append(articles, article)
		}
	}

	if len(articles) == 0 {
		return nil, internal_errors.NoArticlesFound
	}
	return articles, nil
}

func (m *RepositoryMock) ClearArticles() error {
	if m.RepositoryError != nil {
		return m.RepositoryError
	}

	m.Articles = []models.Article{}
	return nil
}

type ParserMock struct {
	ParsingError error
}

func (m *ParserMock) ReInit() {
	m.ParsingError = nil
}

func (m *ParserMock) ParseCategory(
	category string,
	articles chan<- models.Article,
	parsingProcessing models.ProcessingChannels,
) {
	if m.ParsingError != nil {
		parsingProcessing.Error <- m.ParsingError
	}

	article := GetFirstArticle()
	article.Category = category
	articles <- article
	parsingProcessing.Done <- true
}
