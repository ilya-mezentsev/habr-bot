package cli

import (
	"fmt"
	"habr-bot/source/interfaces"
	"habr-bot/source/models"
)

type Presenter struct {
}

func New() interfaces.ArticlesPresenter {
	return Presenter{}
}

func (p Presenter) ShowAsButtons(categories []string) error {
	fmt.Println("All categories:")
	for idx, category := range categories {
		fmt.Printf("\t%d. %s\n", idx+1, category)
	}

	return nil
}

func (p Presenter) ShowArticles(articles []models.Article) error {
	fmt.Printf("Articles of `%s` category\n", articles[0].Category)
	for idx, article := range articles {
		fmt.Printf("\t%d. %s (%s)\n", idx+1, article.Title, article.Link)
	}

	return nil
}
