package processor

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type UserInput struct {
	TimeFrom    string
	TimeTo      string
	Postcode    string
	WordsToFind []string
}

// Function asks user for data in CLI needed to make calculations and returns it.
func GetSettings() UserInput {
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
