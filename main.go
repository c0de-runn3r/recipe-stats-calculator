package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
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

type UserInput struct {
	TimeFrom    string
	TimeTo      string
	Postcode    string
	WordsToFind []string
}

func processENV() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// to get data processing settings from user
func getSettings() UserInput {
	var input UserInput

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter the postcode to look for (example: 10120): ")
		postcode, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("An error occured while reading input. Please try again", err)
			continue
		}
		postcode = strings.TrimSuffix(postcode, "\n")
		ok, _ := regexp.Match("[0-9]{5}", []byte(postcode))
		if !ok {
			fmt.Println("Wrong postcode format. Please try again")
			continue
		}
		input.Postcode = postcode
		break
	}
	for {
		fmt.Print("Enter the time frame in format 'xxAM xxPM (example: 10AM 3PM): ")
		timeFrame, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("An error occured while reading input. Please try again", err)
			continue
		}
		timeFrame = strings.TrimSuffix(timeFrame, "\n")
		ok, _ := regexp.Match("(.?[0-9])AM (.?[0-9])PM", []byte(timeFrame))
		if !ok {
			fmt.Println("Wrong postcode format. Please try again")
			continue
		}

		times := strings.Split(timeFrame, " ")

		input.TimeFrom = times[0]
		input.TimeTo = times[1]
		break
	}
	for {
		fmt.Print("Enter the words to look for separated by coma and whitespace (example: Potato, Veggie, Mushroom): ")
		wordList, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("An error occured while reading input. Please try again", err)
			continue
		}
		wordList = strings.TrimSuffix(wordList, "\n")
		ok, _ := regexp.Match("(([a-zA-Z]+)(, )?)+", []byte(wordList))
		if !ok {
			fmt.Println("Wrong wordlist format. Please try again")
			continue
		}
		wordArr := strings.Split(wordList, ", ")
		input.WordsToFind = wordArr
		break
	}
	return input
}

func main() {
	processENV()
	inputFile := os.Getenv("INPUT_FILE")

	usrInput := getSettings()
	fmt.Println("Opening input file...")
	jsonFile, err := os.Open(inputFile)
	if err != nil {
		log.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)

	fmt.Println("Reading input file...")
	var deliveries DeliveriesRaw
	json.Unmarshal(byteValue, &deliveries)

	fmt.Println("Calculating data...")
	uniqueRecipes := mapRecipesToQuantity(&deliveries)
	sortedListOfRecipes := makeSortedRecipesList(uniqueRecipes) // need to use for loop on map to get quantity
	postcodesMap := mapPostcodesToQuantity(&deliveries)
	counter, err := countDeliveriesToPostcodeBetweenTimes(usrInput.TimeFrom, usrInput.TimeTo, usrInput.Postcode, deliveries)
	if err != nil {
		log.Println(err)
	}
	matches := findMatches(*sortedListOfRecipes, usrInput.WordsToFind)
	writeToJSON(*uniqueRecipes, *sortedListOfRecipes, *postcodesMap, counter, matches, usrInput)
	fmt.Println("Done. You can find JSON file with processed data in output directory.")
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

// to sort it alphabetically
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

func writeToJSON(uniqueRecipesMap map[string]int, sortedListOfRecipes []string, postcodesMap map[string]int, deliveriesForPostcode int, matches []string, usrInput UserInput) {
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
	out.CountPerPostcodeAndTime.From = usrInput.TimeFrom
	out.CountPerPostcodeAndTime.To = usrInput.TimeTo
	out.CountPerPostcodeAndTime.Postcode = usrInput.Postcode

	out.MatchByName = matches

	err := writeToFile(out)
	if err != nil {
		log.Println(err)
	}
}

func writeToFile(data JSONoutput) error {
	file, err := json.MarshalIndent(data, "", "	")
	if err != nil {
		return err
	}

	outputFile := os.Getenv("OUTPUT_FILE")
	outputDir := os.Getenv("OUTPUT_DIR")
	_, err = os.ReadDir(outputDir)
	if err != nil {
		os.Mkdir(outputDir, 0644)
	}
	os.Chmod(outputDir, 0777)
	path := path.Join(outputDir, outputFile)
	err = os.WriteFile(path, file, 0777)
	if err != nil {
		return err
	}
	return nil
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
