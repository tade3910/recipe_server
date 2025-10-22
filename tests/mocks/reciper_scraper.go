package mocks

import (
	"errors"
)

type MockRecipeScraper struct{}

var MockIngredients = [][]string{
	{"egg", "flour"},
	{"milk", "sugar"},
}

var MockInstructions = [][]string{
	{"mix ingredients"},
	{"bake at 350"},
}

const ScrapingError = "scraping error"

const ErrorUrl = "https://error"

func (scraper *MockRecipeScraper) ScrapeRecipe(url string) (
	allIngredients [][]string,
	allInstructions [][]string,
	err error) {

	if url == ErrorUrl {
		err = errors.New(ScrapingError)
		return
	}
	allIngredients = MockIngredients
	allInstructions = MockInstructions
	return
}
