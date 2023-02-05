package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// For reading from file
type DeliveryInfoRaw struct {
	Postcode     string `json:"postcode"`
	Recipe       string `json:"recipe"`
	DeliveryTime string `json:"delivery"`
}

type DeliveriesRaw []*DeliveryInfoRaw

// for writingJSON
type JSONoutput struct {
	UniqueRecipeCount int `json:"unique_recipe_count"`
	CountPerRecipe    []struct {
		Recipe string `json:"recipe"`
		Count  int    `json:"count"`
	} `json:"count_per_recipe"`
	BusiestPostcode struct {
		Postcode      string `json:"postcode"`
		DeliveryCount int    `json:"delivery_count"`
	} `json:"busiest_postcode"`
	CountPerPostcodeAndTime struct {
		Postcode      string `json:"postcode"`
		From          string `json:"from"`
		To            string `json:"to"`
		DeliveryCount int    `json:"delivery_count"`
	} `json:"count_per_postcode_and_time"`
	MatchByName []string `json:"match_by_name"`
}

const (
	fromENV     = "10AM"
	toENV       = "3PM"
	postcodeENV = "10120"
)

var wordsToFind = []string{"Potato", "Veggie", "Mushroom"}

func main() {
	t1 := time.Now()

	jsonFile, err := os.Open("hf_test_calculation_fixtures.json")
	if err != nil {
		log.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)

	fmt.Println("Open+Read time: ", time.Since(t1).Seconds())
	t2 := time.Now()

	var deliveries DeliveriesRaw
	json.Unmarshal(byteValue, &deliveries)

	fmt.Println("Unmarshal time: ", time.Since(t2).Seconds())
	t3 := time.Now()

	uniqueRecipes := mapRecipesToQuantity(&deliveries)

	fmt.Println("Mapping time: ", time.Since(t3).Seconds())

	sortedListOfRecipes := makeSortedRecipesList(uniqueRecipes) // need to use for loop on map to get quantity

	postcodesMap := mapPostcodesToQuantity(&deliveries)

	counter, err := countDeliveriesToPostcodeBetweenTimes(fromENV, toENV, postcodeENV, deliveries)
	if err != nil {
		log.Println(err)
	}
	matches := findMatches(*sortedListOfRecipes, wordsToFind)
	writeToJSON(*uniqueRecipes, *sortedListOfRecipes, *postcodesMap, counter, matches)
}

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

func makeSortedRecipesList(uniqueRecipesMap *map[string]int) *[]string {
	recipesList := make([]string, 0, len(*uniqueRecipesMap))
	for k := range *uniqueRecipesMap {
		recipesList = append(recipesList, k)
	}
	sort.Strings(recipesList)
	return &recipesList
}

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

func writeToJSON(uniqueRecipesMap map[string]int, sortedListOfRecipes []string, postcodesMap map[string]int, deliveriesForPostcode int, matches []string) {
	var out JSONoutput

	out.UniqueRecipeCount = len(uniqueRecipesMap)

	for i := 0; i < len(sortedListOfRecipes); i++ {
		out.CountPerRecipe = append(out.CountPerRecipe, struct {
			Recipe string `json:"recipe"`
			Count  int    `json:"count"`
		}{
			Recipe: sortedListOfRecipes[i],
			Count:  uniqueRecipesMap[sortedListOfRecipes[i]],
		})
	}

	postcode, counter := busiestPostcode(&postcodesMap)
	out.BusiestPostcode.Postcode = postcode
	out.BusiestPostcode.DeliveryCount = counter

	out.CountPerPostcodeAndTime.DeliveryCount = deliveriesForPostcode
	out.CountPerPostcodeAndTime.From = fromENV
	out.CountPerPostcodeAndTime.To = toENV
	out.CountPerPostcodeAndTime.Postcode = postcodeENV

	out.MatchByName = matches

	file, _ := json.MarshalIndent(out, "", "	")
	_ = os.WriteFile("out.json", file, 0644)
}

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
