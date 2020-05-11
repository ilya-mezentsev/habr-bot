package articles

import (
	"interfaces"
	"io/ioutil"
	"log"
	"mock"
	"os"
	"services/errors"
	"testing"
	"utils"
)

var (
	repository = mock.RepositoryMock{}
	parser     = mock.ParserMock{}
	service    interfaces.ArticlesService
)

func init() {
	repository.ReInit()
	parser.ReInit()

	service = New(&repository, &parser, mock.GetAllCategories())
}

func resetMocks() {
	repository.ReInit()
	parser.ReInit()
}

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func TestService_GetArticlesSuccess(t *testing.T) {
	defer resetMocks()

	category := mock.GetAllCategories()[0]
	articles, err := service.GetArticles(category)
	expectedArticles, _ := repository.GetByCategory(category)

	utils.AssertNil(err, t)
	utils.AssertEqual(expectedArticles[0], articles[0], t)
}

func TestService_GetArticlesNoArticlesFound(t *testing.T) {
	defer resetMocks()

	_, err := service.GetArticles("blah-blah")

	utils.AssertErrorsEqual(errors.NoArticlesForCategory, err, t)
}

func TestService_GetArticlesInternalError(t *testing.T) {
	defer resetMocks()

	category := mock.GetBadArticle().Category
	_, err := service.GetArticles(category)

	utils.AssertErrorsEqual(errors.InternalError, err, t)
}

func TestService_GetCategories(t *testing.T) {
	categories := service.GetCategories()

	for categoryIndex, expectedCategory := range mock.GetAllCategories() {
		utils.AssertEqual(expectedCategory, categories[categoryIndex], t)
	}
}

func TestService_ParseCategorySuccess(t *testing.T) {
	defer resetMocks()

	category := mock.GetAllCategories()[0]
	err := service.ParseCategory(category)
	expectedArticlesCount := len(mock.GetAllArticles())

	utils.AssertNil(err, t)
	utils.AssertEqual(expectedArticlesCount, len(repository.Articles), t)
}

func TestService_ParseCategorySavingError(t *testing.T) {
	defer resetMocks()

	err := service.ParseCategory(mock.GetBadArticle().Category)

	utils.AssertErrorsEqual(errors.SavingArticlesError, err, t)
}

func TestService_ParseCategoryParsingError(t *testing.T) {
	defer resetMocks()

	parser.ParsingError = errors.InternalError
	err := service.ParseCategory(mock.GetAllCategories()[0])

	utils.AssertErrorsEqual(errors.ParsingCategoryError, err, t)
}

func TestService_ParseAllSuccess(t *testing.T) {
	defer resetMocks()

	expectedArticlesCount := len(repository.Articles)
	err := service.ParseAll()

	utils.AssertNil(err, t)
	utils.AssertEqual(expectedArticlesCount, len(repository.Articles), t)
}

func TestService_ParseAllParsingError(t *testing.T) {
	defer resetMocks()

	parser.ParsingError = errors.InternalError
	err := service.ParseAll()

	utils.AssertErrorsEqual(errors.ParsingCategoryError, err, t)
}

func TestService_ParseAllClearArticlesError(t *testing.T) {
	defer resetMocks()

	repository.RepositoryError = errors.InternalError
	err := service.ParseAll()

	utils.AssertErrorsEqual(errors.ClearArticlesError, err, t)
}
