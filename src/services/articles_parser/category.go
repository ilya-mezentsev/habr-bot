package articles_parser

import (
	"fmt"
	"strings"
)

type Category struct {
	name   string
	filter string
}

const categoryFilterSplitter = ":"

func NewCategory(data string) Category {
	parts := strings.Split(data, categoryFilterSplitter)
	var name, filter string
	name = parts[0]
	if len(parts) > 1 {
		filter = parts[1]
	}

	return Category{name, filter}
}

func (c Category) GetName() string {
	return c.name
}

func (c Category) GetFilter() string {
	return c.filter
}

func CombineCategoriesWithFilters(categories, filters []string) []string {
	var combinedCategories []string
	for _, category := range categories {
		combinedCategories = append(combinedCategories, category)

		for _, filter := range filters {
			combinedCategories = append(
				combinedCategories,
				fmt.Sprintf("%s%s%s", category, categoryFilterSplitter, filter),
			)
		}
	}

	return combinedCategories
}
