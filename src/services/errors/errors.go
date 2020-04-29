package errors

import "errors"

var (
	InternalError         = errors.New("internal-error")
	NoArticlesForCategory = errors.New("no-articles-for-category")
	SavingArticlesError   = errors.New("saving-articles-error")
	ParsingCategoryError  = errors.New("parsing-category-error")
)
