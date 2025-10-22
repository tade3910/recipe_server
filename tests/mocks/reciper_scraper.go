package mocks

import "fmt"

type MockRecipeScraper struct{}

var MockIngredients = [][]string{
	{"egg", "flour"},
	{"milk", "sugar"},
}

var MockInstructions = [][]string{
	{"mix ingredients"},
	{"bake at 350"},
}

var ErrorMock = fmt.Errorf("scraping error")

func (scraper *MockRecipeScraper) ScrapeRecipe(url string) (
	allIngredients [][]string,
	allInstructions [][]string,
	err error) {

	if url == "error" {
		err = ErrorMock
		return
	}
	allIngredients = MockIngredients
	allInstructions = MockInstructions
	return
}
