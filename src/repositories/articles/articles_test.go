package articles

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"interfaces"
	"internal_errors"
	"mock"
	"models"
	"os"
	"plugins/config"
	"testing"
	"utils"
)

var (
	db         *sqlx.DB
	repository interfaces.ArticlesRepository
)

func init() {
	dbFile, err := config.GetDBPath()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	db, err = sqlx.Open("sqlite3", dbFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	repository = New(db)
}

func TestRepository_GetByCategorySuccess(t *testing.T) {
	mock.InitTable(db)
	defer mock.DropTable(db)

	category := mock.GetFirstArticle().Category
	articles, err := repository.GetByCategory(category)

	utils.AssertNil(err, t)
	for _, article := range articles {
		utils.AssertEqual(category, article.Category, t)
	}
}

func TestRepository_GetByCategoryNoArticles(t *testing.T) {
	mock.InitTable(db)
	defer mock.DropTable(db)

	_, err := repository.GetByCategory("")
	utils.AssertErrorsEqual(internal_errors.NoArticlesFound, err, t)
}

func TestRepository_GetByCategorySomeError(t *testing.T) {
	mock.DropTable(db)

	_, err := repository.GetByCategory(mock.GetFirstArticle().Category)

	utils.AssertNotNil(err, t)
}

func TestRepository_SaveSuccess(t *testing.T) {
	mock.InitTable(db)
	defer mock.DropTable(db)

	article := mock.GetNewArticle()
	articles, _ := repository.GetByCategory(article.Category)
	expectedArticlesCount := len(articles) + 1
	articlesChan := make(chan models.Article)
	trySave := make(chan bool)
	savingProcessing := models.ProcessingChannels{
		Done:  make(chan bool),
		Error: make(chan error),
	}

	go repository.Save(
		articlesChan,
		trySave,
		savingProcessing,
	)

	articlesChan <- article
	trySave <- true

	for {
		select {
		case err := <-savingProcessing.Error:
			t.Logf("Error: %v\n", err)
			t.Fail()
			return
		case <-savingProcessing.Done:
			articles, _ := repository.GetByCategory(article.Category)

			utils.AssertEqual(expectedArticlesCount, len(articles), t)
			return
		}
	}
}
