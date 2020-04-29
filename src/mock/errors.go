package mock

import (
	"errors"
	"models"
)

var (
	someRepositoryError = errors.New("some repository error")
	someParsingError    = errors.New("some parsing error")
)

func GetBadArticle() models.Article {
	return models.Article{}
}
