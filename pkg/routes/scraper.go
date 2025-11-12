package routes

import (
	"net/http"
	"strings"

	util "github.com/tade3910/recipe_server/pkg"
	"github.com/tade3910/recipe_server/pkg/services"
	"gorm.io/gorm"
)

type scraperHandler struct {
	db      *gorm.DB
	scraper services.RecipeScraperInterface
}

func NewscraperHandler(db *gorm.DB, scraper services.RecipeScraperInterface) *scraperHandler {
	return &scraperHandler{
		db:      db,
		scraper: scraper,
	}
}

func (handler *scraperHandler) ParseRecipe(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if strings.TrimSpace(url) == "" {
		util.RespondWithError(w, http.StatusBadRequest, "No url received")
	}
	allIngredients, allInstructions, err := handler.scraper.ScrapeRecipe(url)
	if err != nil {
		util.RespondWithError(w, http.StatusNotFound, err)
		return
	}
	response := &map[string][][]string{
		"ingredients":  allIngredients,
		"instructions": allInstructions,
	}
	util.RespondWithJSON(w, http.StatusOK, response)
}
