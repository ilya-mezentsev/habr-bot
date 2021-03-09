package articles

import (
	"github.com/jmoiron/sqlx"
	"habr-bot/source/interfaces"
	"habr-bot/source/internal_errors"
	"habr-bot/source/models"
)

const (
	clearArticlesQuery         = `DELETE FROM articles WHERE id != 0`
	clearCategoryArticlesQuery = `DELETE FROM articles WHERE category = $1`
	getByCategoryQuery         = `SELECT title, category, link FROM articles WHERE category = $1`
	addArticleQuery            = `
	INSERT INTO articles(title, category, link)
	VALUES(:title, :category, :link)`
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) interfaces.ArticlesRepository {
	return Repository{db}
}

func (r Repository) Save(
	articles chan models.Article,
	trySave chan bool,
	processing models.ProcessingChannels,
) {
	tx := r.db.MustBegin()

	for {
		select {
		case article := <-articles:
			_, err := tx.NamedExec(addArticleQuery, article)
			if err != nil {
				processing.Error <- err
				return
			}
		case <-trySave:
			err := tx.Commit()
			if err != nil {
				processing.Error <- err
				return
			} else {
				processing.Done <- true
				return
			}
		}
	}
}

func (r Repository) GetByCategory(category string) ([]models.Article, error) {
	var articles []models.Article
	err := r.db.Select(&articles, getByCategoryQuery, category)
	if len(articles) == 0 {
		err = internal_errors.NoArticlesFound
	}

	return articles, err
}

func (r Repository) ClearArticles() error {
	_, err := r.db.Exec(clearArticlesQuery)

	return err
}

func (r Repository) ClearCategoryArticles(category string) error {
	_, err := r.db.Exec(clearCategoryArticlesQuery, category)

	return err
}
