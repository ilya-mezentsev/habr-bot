package articles

import (
	"interfaces"
	"internal_errors"
	"models"
	"plugins/logger"
	"services/category_format"
	"services/errors"
)

type Service struct {
	repository interfaces.ArticlesRepository
	parser     interfaces.ArticlesParserService
	categories []string
	filters    []string
}

func New(
	repository interfaces.ArticlesRepository,
	parser interfaces.ArticlesParserService,
	categories []string,
	filters []string,
) interfaces.ArticlesService {
	return Service{repository, parser, categories, filters}
}

func (s Service) ParseAll() error {
	parsedCategoriesCount := 0
	categoriesCount := len(s.categories)
	parsingProcessing := models.ProcessingChannels{
		Done:  make(chan bool),
		Error: make(chan error),
	}

	if err := s.repository.ClearArticles(); err != nil {
		logger.WithFields(logger.Fields{
			MessageTemplate: "Error while parsing all categories: %v",
			Args: []interface{}{
				err,
			},
		}, logger.Error)

		return errors.ClearArticlesError
	}

	for _, category := range category_format.CombineCategoriesWithFilters(s.categories, s.filters) {
		go s.parseCategory(category, parsingProcessing)
	}

	for {
		select {
		case err := <-parsingProcessing.Error:
			return err
		case <-parsingProcessing.Done:
			parsedCategoriesCount++
		default:
			if parsedCategoriesCount >= categoriesCount {
				return nil
			}
		}
	}
}

func (s Service) parseCategory(category string, parsingProcessing models.ProcessingChannels) {
	err := s.ParseCategory(category)

	if err != nil {
		parsingProcessing.Error <- err
	} else {
		parsingProcessing.Done <- true
	}
}

func (s Service) ParseCategory(category string) error {
	articles := make(chan models.Article)
	parsingDone := make(chan bool)
	parsingProcessing := models.ProcessingChannels{
		Done:  parsingDone,
		Error: make(chan error),
	}
	savingProcessing := models.ProcessingChannels{
		Done:  make(chan bool),
		Error: make(chan error),
	}

	go func() {
		err := s.repository.ClearCategoryArticles(category)
		if err != nil {
			savingProcessing.Error <- err
		}
	}()
	go s.parser.ParseCategory(
		category,
		articles,
		parsingProcessing,
	)
	go s.repository.Save(
		articles,
		parsingDone,
		savingProcessing,
	)

	for {
		select {
		case err := <-parsingProcessing.Error:
			s.logParsingCategoryError(category, err)
			return errors.ParsingCategoryError
		case err := <-savingProcessing.Error:
			s.logParsingCategoryError(category, err)
			return errors.SavingArticlesError
		case <-savingProcessing.Done:
			return nil
		}
	}
}

func (s Service) logParsingCategoryError(category string, err error) {
	logger.WithFields(logger.Fields{
		MessageTemplate: "Error while parsing category: %v",
		Args: []interface{}{
			err,
		},
		Optional: map[string]interface{}{
			"category": category,
		},
	}, logger.Error)
}

func (s Service) GetCategories() []string {
	return s.categories
}

func (s Service) GetFilters() []string {
	return s.filters
}

func (s Service) GetArticles(category string) ([]models.Article, error) {
	articles, err := s.repository.GetByCategory(category)

	switch err {
	case nil:
		return articles, nil
	case internal_errors.NoArticlesFound:
		s.logGettingArticlesError(category, err)
		return nil, errors.NoArticlesForCategory
	default:
		return nil, errors.InternalError
	}
}

func (s Service) logGettingArticlesError(category string, err error) {
	logger.WithFields(logger.Fields{
		MessageTemplate: "Error while getting articles: %v",
		Args: []interface{}{
			err,
		},
		Optional: map[string]interface{}{
			"category": category,
		},
	}, logger.Error)
}
