package articles_parser

import (
	"testing"
	"utils"
)

func TestCategory_GetName(t *testing.T) {
	name := "hey"

	utils.AssertEqual(name, NewCategory(name).GetName(), t)
}

func TestCategory_GetFilter(t *testing.T) {
	name := "name:filter"

	utils.AssertEqual("filter", NewCategory(name).GetFilter(), t)
}

func TestCategory_GetFilterEmpty(t *testing.T) {
	name := "hey"

	utils.AssertEqual("", NewCategory(name).GetFilter(), t)
}
