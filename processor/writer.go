package processor

import (
	"encoding/json"
	"log"
	"os"
	"path"
)

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

// Function takes calculated data, structs it and writes into output JSON file.
func WriteToJSON(calculatedData CalculatedData, usrInput UserInput) {
	var out JSONoutput

	out.UniqueRecipeCount = len(calculatedData.UniqueRecipes)

	for i := 0; i < len(calculatedData.SortedListOfRecipes); i++ {
		out.CountPerRecipe = append(out.CountPerRecipe, struct {
			Recipe string `json:"recipe"`
			Count  int    `json:"count"`
		}{
			Recipe: calculatedData.SortedListOfRecipes[i],
			Count:  calculatedData.UniqueRecipes[calculatedData.SortedListOfRecipes[i]],
		})
	}

	postcode, counter := busiestPostcode(&calculatedData.PostcodesMap)
	out.BusiestPostcode.Postcode = postcode
	out.BusiestPostcode.DeliveryCount = counter

	out.CountPerPostcodeAndTime.DeliveryCount = calculatedData.Counter
	out.CountPerPostcodeAndTime.From = usrInput.TimeFrom
	out.CountPerPostcodeAndTime.To = usrInput.TimeTo
	out.CountPerPostcodeAndTime.Postcode = usrInput.Postcode

	out.MatchByName = calculatedData.Matches

	err := writeToFile(out)
	if err != nil {
		log.Println(err)
	}
}

// to write into file
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
