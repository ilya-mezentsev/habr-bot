package mock

import (
	"github.com/jmoiron/sqlx"
	"habr-bot/source/models"
)

const (
	dropTableQuery   = `DROP TABLE IF EXISTS articles;`
	createTableQuery = `CREATE TABLE IF NOT EXISTS articles(
		id INTEGER PRIMARY KEY,
		title TEXT NOT NULL,
		category TEXT NOT NULL,
		link TEXT NOT NULL
	)`
	addArticleQuery = `
	INSERT INTO articles(title, category, link)
	VALUES(:title, :category, :link)`
)

var (
	articles = []map[string]interface{}{
		{"title": "Title 1", "category": "Category 1", "link": "http://link1.com"},
		{"title": "Title 2", "category": "Category 2", "link": "http://link2.com"},
		{"title": "Title 3", "category": "Category 3", "link": "http://link3.com"},
	}
)

func GetNewArticle() models.Article {
	return models.Article{
		Title:    "Title 4",
		Category: "Category 3",
		Link:     "http://link4.com",
	}
}

func GetFirstArticle() models.Article {
	article := articles[0]
	return models.Article{
		Title:    article["title"].(string),
		Category: article["category"].(string),
		Link:     article["link"].(string),
	}
}

func GetAllArticles() []models.Article {
	var a []models.Article
	for _, article := range articles {
		a = append(a, models.Article{
			Title:    article["title"].(string),
			Category: article["category"].(string),
			Link:     article["link"].(string),
		})
	}

	return a
}

func GetAllCategories() []string {
	var categories []string
	for _, article := range articles {
		categories = append(categories, article["category"].(string))
	}

	return categories
}

func CreateTableIfNotExists(db *sqlx.DB) {
	db.MustExec(createTableQuery)
}

func InitTable(db *sqlx.DB) {
	DropTable(db)
	CreateTableIfNotExists(db)

	tx := db.MustBegin()
	for _, article := range articles {
		_, err := tx.NamedExec(addArticleQuery, article)
		if err != nil {
			panic(err)
		}
	}

	err := tx.Commit()
	if err != nil {
		panic(err)
	}
}

func DropTable(db *sqlx.DB) {
	db.MustExec(dropTableQuery)
}
