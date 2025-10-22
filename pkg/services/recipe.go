package services

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"unicode"

	"golang.org/x/net/html"
)

type RecipeScraperInterface interface {
	ScrapeRecipe(url string) ([][]string, [][]string, error)
}

type RecipeScraper struct{}

func (scraper *RecipeScraper) ScrapeRecipe(url string) (
	allIngredients [][]string,
	allInstructions [][]string,
	err error) {
	// Make the HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("failed to fetch the URL: %s, Status Code: %d", url, resp.StatusCode)
		return
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return
	}
	// scraper.debugNode(doc, 0)
	allIngredients = parseIngredients(doc)
	allInstructions = parseInstructions(doc)

	return
}

func getListNodes(n *html.Node, listNodes *[]*html.Node) {
	if n.Type == html.ElementNode && (n.Data == "ul" || n.Data == "ol") {
		*listNodes = append(*listNodes, n)
		return
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		getListNodes(c, listNodes)
	}
}

func getEnglishString(s string) string {
	englishString := ""
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsSpace(r) || isPunctuation(r) || isNumber(r) {
			englishString += string(r)
		}
	}
	return englishString
}

func isNumber(r rune) bool {
	return unicode.IsDigit(r) || unicode.Is(unicode.No, r)
}

func isPunctuation(r rune) bool {
	switch r {
	case '.', ',', ';', ':', '\'', '"', '!', '?', '-', '(', ')', '[', ']', '{', '}', '/', '\\', '&', '%', '$', '#', '@', '*', '+', '=', '<', '>', '|', '~', '`':
		return true
	default:
		return false
	}
}

func getListChild(n *html.Node) string {
	child := ""
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			child += getEnglishString(c.Data)
		} else if c.Type == html.ElementNode && (c.Data == "noscript" || c.Data == "figcaption") {
			continue
		}
		child += getListChild(c)
	}
	return strings.TrimSpace(child)
}

func getLists(n *html.Node) []string {
	lists := []string{}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "li" {
			child := getListChild(c)
			if len(child) > 0 {
				lists = append(lists, child)
			}
		}
	}
	if len(lists) == 0 {
		return nil
	}
	return lists
}

func matchesTargets(targets []string, check string) bool {
	lowerCaseCheck := strings.TrimSpace((strings.ToLower(check)))
	for _, target := range targets {
		if strings.Compare(strings.ToLower(target), lowerCaseCheck) == 0 {
			return true
		}
	}
	return false
}

func getTargetListNodes(n *html.Node, targetListNodes *[]*html.Node, targets []string) bool {
	if n.Type == html.TextNode && matchesTargets(targets, n.Data) {
		log.Print("Found potential target")
		return true
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if getTargetListNodes(c, targetListNodes, targets) {
			// check if c contains list
			listNodes := []*html.Node{}
			getListNodes(c, &listNodes)
			if len(listNodes) > 0 {
				*targetListNodes = append(*targetListNodes, listNodes...)
				return false
			}
			return true
		}
	}
	return false
}

func parseInstructions(n *html.Node) [][]string {
	instructionListNodes := []*html.Node{}
	getTargetListNodes(n, &instructionListNodes, []string{"Instructions", "Directions"})
	if len(instructionListNodes) == 0 {
		log.Fatal("Could not find instructions")
	}
	instructionLists := [][]string{}
	for _, node := range instructionListNodes {
		instructionList := getLists(node)
		instructionLists = append(instructionLists, instructionList)
	}
	return instructionLists
}

func parseIngredients(n *html.Node) [][]string {
	ingredientListNodes := []*html.Node{}
	getTargetListNodes(n, &ingredientListNodes, []string{"ingredients"})
	if len(ingredientListNodes) == 0 {
		log.Fatal("Could not find ingredients")
	}
	ingredientLists := [][]string{}
	for _, node := range ingredientListNodes {
		ingredientList := getLists(node)
		if isIngredientList(ingredientList) {
			ingredientLists = append(ingredientLists, ingredientList)
		}
	}
	return ingredientLists
}

func isIngredientList(ingredientList []string) bool {
	for _, ingredient := range ingredientList {
		if isNumber([]rune(ingredient)[0]) {
			return true
		}
	}
	return false
}

// Recursive function to print HTML nodes with indentation for better readability
func (scraper *RecipeScraper) debugNode(n *html.Node, indentLevel int) {
	// Add indentation for pretty printing
	indent := fmt.Sprintf("%*s", indentLevel*2, "")

	// Print the opening tag and its attributes (if any)
	if n.Type == html.ElementNode {
		fmt.Printf("%s<%s", indent, n.Data)
		for _, attr := range n.Attr {
			fmt.Printf(" %s=\"%s\"", attr.Key, attr.Val)
		}
		fmt.Println(">") // Opening tag
	}

	// Print text node content (if any)
	if n.Type == html.TextNode {
		fmt.Printf("%s%s\n", indent, n.Data)
	}

	// Recursively print child nodes
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		scraper.debugNode(c, indentLevel+1) // Increase indentation for nested elements
	}

	// Print the closing tag (if it is an element)
	if n.Type == html.ElementNode {
		fmt.Printf("%s</%s>\n", indent, n.Data)
	}
}
