package processor

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

// Struct for raw data in JSON file.
type DeliveryInfoRaw struct {
	Postcode     string `json:"postcode"`
	Recipe       string `json:"recipe"`
	DeliveryTime string `json:"delivery"`
}

// Contains all the raw data readed from JSON file and processed to Go struct format.
type DeliveriesRaw []*DeliveryInfoRaw

// Function reads the JSON file and processes raw JSON data to Go struct format.
func ReadFile(inputFile string) *DeliveriesRaw {
	jsonFile, err := os.Open(inputFile)
	if err != nil {
		log.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)

	var deliveries DeliveriesRaw
	json.Unmarshal(byteValue, &deliveries)
	return &deliveries
}
