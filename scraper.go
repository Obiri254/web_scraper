package main

import (
	//colly import
	"fmt"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

func main() {
	args := os.Args
	url := args[1]
	collector := colly.NewCollector()

	type Dictionary map[string]string

	type RecipeSpecs struct {
		difficulty, prepTime, cookingTime, servingSize, priceTier string
	}

	type Recipe struct {
		url, name	string
		ingredients	[]Dictionary
		specifications	RecipeSpecs
	}

	//Whenever the collector is about to make a new request
	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	collector.OnResponse(func(r *colly.Response) {
		fmt.Println("Got a response from", r.Request.URL)
	})
	collector.OnError(func(r *colly.Response, err error) {
		fmt.Println("Ooh no, an error occurred!:", err)
	})
	var recipes []Recipe
	collector.OnHTML("main", func(main *colly.HTMLElement) {
		recipe := Recipe{}
		ingredients_dictionary := Dictionary{}
		recipe.url = url

		recipe.name = main.ChildText(".gz-title-recipe")
		println("Scraping recipe for:", recipe.name)

		main.ForEach(".gz-name-featured-data", func(i int, specListElement *colly.HTMLElement)  {
			if strings.Contains(specListElement.Text, "Difficolta: ") {
				recipe.specifications.difficulty = specListElement.ChildText("strong")
			}
			if strings.Contains(specListElement.Text, "Preparazione: ") {
				recipe.specifications.prepTime = specListElement.ChildText("strong")
			}
			if strings.Contains(specListElement.Text, "Cottura: ") {
				recipe.specifications.cookingTime = specListElement.ChildText("strong")
			}
			if strings.Contains(specListElement.Text, "Dosi per: ") {
				recipe.specifications.servingSize = specListElement.ChildText("strong")
			}
			if strings.Contains(specListElement.Text, "Costo: ") {
				recipe.specifications.priceTier = specListElement.ChildText("strong")
			}
		})
		main.ForEach(".gz-ingredient", func (i int, ingredient *colly.HTMLElement)  {
			ingredients_dictionary[ingredient.ChildText("a")] = ingredient.ChildText("span")
		})
		recipe.ingredients = append(recipe.ingredients, ingredients_dictionary)
		recipes = append(recipes, recipe)
	})
	collector.Visit(url)

}