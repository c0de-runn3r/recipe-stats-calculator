package main

import (
	"fmt"
	"log"
	"os"
	"recipe-calc/processor"

	"github.com/joho/godotenv"
)

func main() {
	processENV()
	inputFile := os.Getenv("INPUT_FILE")

	usrInput := processor.GetSettings()

	fmt.Println("Reading input file...")
	deliveries := processor.ReadFile(inputFile)

	fmt.Println("Calculating data...")
	calculatedData := processor.MakeCalculations(*deliveries, usrInput)

	processor.WriteToJSON(*calculatedData, usrInput)
	fmt.Println("Done. You can find JSON file with processed data in output directory.")
}

// to process enviromental variables from .env file
func processENV() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
