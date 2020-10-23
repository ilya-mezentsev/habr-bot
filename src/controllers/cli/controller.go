package cli

import (
	"bufio"
	"controllers/commands"
	"fmt"
	"interfaces"
	"models"
	"os"
	"services/errors"
	"strings"
)

var currentCategory = ""

type Controller struct {
	service   interfaces.ArticlesService
	presenter interfaces.ArticlesPresenter
}

func New(
	service interfaces.ArticlesService,
	presenter interfaces.ArticlesPresenter,
) Controller {
	return Controller{service, presenter}
}

func (c Controller) Run(processing models.ProcessingChannels) {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter command: ")
		command, err := reader.ReadString('\n')
		if err != nil {
			processing.Error <- err
			return
		}

		err = c.execute(strings.TrimSuffix(command, "\n"))
		if err != nil {
			processing.Error <- err
			return
		}
	}
}

func (c Controller) execute(command string) error {
	switch command {
	case commands.ParseAll:
		return c.parseAll()
	case commands.GetCategories:
		return c.showCategories()
	case commands.GetArticles:
		return c.getArticles()
	case commands.ParseArticles:
		return c.parseArticles()
	case commands.Clean:
		return c.clean()
	default:
		return c.unknown(command)
	}
}

func (c Controller) parseAll() error {
	err := c.service.ParseAll()
	if err != nil {
		return err
	}

	fmt.Println("All categories is parsed")
	return nil
}

func (c Controller) showCategories() error {
	_ = c.presenter.ShowAsButtons(c.service.GetCategories())

	return nil
}

func (c Controller) getArticles() error {
	if !c.isValidCategory(currentCategory) {
		fmt.Println("Category is not set or invalid")
		return nil
	}

	articles, err := c.service.GetArticles(currentCategory)
	switch {
	case err == nil:
		break
	case err == errors.NoArticlesForCategory:
		fmt.Println("No parsed articles for category")
		err = c.parseArticles()
		if err != nil {
			return err
		}
		return c.getArticles()
	default:
		return err
	}

	_ = c.presenter.ShowArticles(articles)
	return nil
}

func (c Controller) parseArticles() error {
	if !c.isValidCategory(currentCategory) {
		fmt.Println("Category is not set or invalid")
		return nil
	}

	err := c.service.ParseCategory(currentCategory)
	if err != nil {
		return err
	}

	fmt.Println("Articles is parsed")
	return nil
}

func (c Controller) clean() error {
	currentCategory = ""

	return nil
}

func (c Controller) unknown(command string) error {
	if c.isValidCategory(command) {
		currentCategory = command
		fmt.Printf("Category %s set as current\n", command)
	} else {
		fmt.Println("Unknown command")
		fmt.Printf("Available commands:\n%s\n", commands.All())
	}

	return nil
}

func (c Controller) isValidCategory(category string) bool {
	for _, c := range c.service.GetCategories() {
		if category == c {
			return true
		}
	}

	return false
}
