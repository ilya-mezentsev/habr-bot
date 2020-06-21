package category_format

import (
	"fmt"
	"testing"
	"utils"
)

func TestCategory_GetName(t *testing.T) {
	name := "hey"

	utils.AssertEqual(name, New(name).GetName(), t)
}

func TestCategory_GetFilter(t *testing.T) {
	name := "name:filter"

	utils.AssertEqual("filter", New(name).GetFilter(), t)
}

func TestCategory_GetFilterEmpty(t *testing.T) {
	name := "hey"

	utils.AssertEqual("", New(name).GetFilter(), t)
}

func TestGetFormattedCategory(t *testing.T) {
	name, filter := "hey", "filter"

	utils.AssertEqual(
		fmt.Sprintf("%s%s%s", name, categoryFilterSplitter, filter),
		GetFormattedCategory(name, filter),
		t,
	)
}

func TestGetFormattedCategory_EmptyFilter(t *testing.T) {
	name, filter := "hey", ""

	utils.AssertEqual(
		name,
		GetFormattedCategory(name, filter),
		t,
	)
}

func TestCombineCategoriesWithFilters(t *testing.T) {
	names := []string{"n1", "n2"}
	filters := []string{"f1", "f2"}
	expected := []string{"n1", "n1:f1", "n1:f2", "n2", "n2:f1", "n2:f2"}

	for index, categoryName := range CombineCategoriesWithFilters(names, filters) {
		utils.AssertEqual(
			expected[index],
			categoryName,
			t,
		)
	}
}
