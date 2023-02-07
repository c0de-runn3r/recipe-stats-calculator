package processor

import (
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type CalculatedData struct {
	UniqueRecipes       map[string]int
	SortedListOfRecipes []string
	PostcodesMap        map[string]int
	Counter             int
	Matches             []string
}

// Function takes data from JSON file and user CLI input and returns processed data.
func MakeCalculations(deliveries DeliveriesRaw, usrInput UserInput) *CalculatedData {
	uniqueRecipes := mapRecipesToQuantity(&deliveries)
	sortedListOfRecipes := makeSortedRecipesList(uniqueRecipes) // need to use for loop on map to get quantity
	postcodesMap := mapPostcodesToQuantity(&deliveries)
	counter, err := countDeliveriesToPostcodeBetweenTimes(usrInput.TimeFrom, usrInput.TimeTo, usrInput.Postcode, deliveries)
	if err != nil {
		log.Println(err)
	}
	matches := findMatches(*sortedListOfRecipes, usrInput.WordsToFind)

	return &CalculatedData{
		UniqueRecipes:       *uniqueRecipes,
		SortedListOfRecipes: *sortedListOfRecipes,
		PostcodesMap:        *postcodesMap,
		Counter:             counter,
		Matches:             matches,
	}
}

// to count quantity of every unique recipe
func mapRecipesToQuantity(deliveries *DeliveriesRaw) *map[string]int {
	m := make(map[string]int)
	for _, v := range *deliveries {
		_, ok := m[v.Recipe]
		if !ok {
			m[v.Recipe] = 0
		}
		m[v.Recipe]++
	}
	return &m
}

// to count quantity of every unique postcode
func mapPostcodesToQuantity(deliveries *DeliveriesRaw) *map[string]int {
	m := make(map[string]int)
	for _, v := range *deliveries {
		_, ok := m[v.Postcode]
		if !ok {
			m[v.Postcode] = 0
		}
		m[v.Postcode]++
	}
	return &m
}

// to sort recipes alphabetically
func makeSortedRecipesList(uniqueRecipesMap *map[string]int) *[]string {
	recipesList := make([]string, 0, len(*uniqueRecipesMap))
	for k := range *uniqueRecipesMap {
		recipesList = append(recipesList, k)
	}
	sort.Strings(recipesList)
	return &recipesList
}

// to find postcode with the most orders
func busiestPostcode(postcodeMap *map[string]int) (string, int) {
	var postcode string
	var counter int
	for k, v := range *postcodeMap {
		if v > counter {
			postcode = k
			counter = v
		}
	}
	return postcode, counter
}

// to count deliveries for given postcode during given time frames
func countDeliveriesToPostcodeBetweenTimes(fromStr string, toStr string, postcode string, data DeliveriesRaw) (int, error) {
	from, err := strconv.Atoi(fromStr[0 : len(fromStr)-2])
	if err != nil {
		return 0, err
	}

	to, err := strconv.Atoi(toStr[0 : len(toStr)-2])
	if err != nil {
		return 0, err
	}

	counter := 0
	for i := 0; i < len(data); i++ {
		if data[i].Postcode == postcode {
			open, close, err := deliveryTime(data[i].DeliveryTime)
			if err != nil {
				return 0, err
			}
			if open <= from && close >= to {
				counter++
			}
		}
	}
	return counter, nil
}

// to parse delivery times from file
func deliveryTime(data string) (int, int, error) {
	patternOpen := regexp.MustCompile("[a-z] (.?[0-9])AM")
	patternClose := regexp.MustCompile("- (.?[0-9])PM")

	openStr := patternOpen.Find([]byte(data))
	closeStr := patternClose.Find([]byte(data))

	open, err := strconv.Atoi(string(openStr)[2 : len(openStr)-2])
	if err != nil {
		return 0, 0, err
	}

	close, err := strconv.Atoi(string(closeStr)[2 : len(closeStr)-2])
	if err != nil {
		return 0, 0, err
	}
	return open, close, nil
}

// to find word matches in list of recipes
func findMatches(sortedListOfRecipes []string, words []string) []string {
	matches := make([]string, 0)

	for i := 0; i < len(sortedListOfRecipes); i++ {
		for _, word := range words {
			if strings.Contains(sortedListOfRecipes[i], word) {
				matches = append(matches, sortedListOfRecipes[i])
			}
		}
	}
	return matches
}
