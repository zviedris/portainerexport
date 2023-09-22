package main

import (
	json "encoding/json"
	"fmt"
	os "os"

	model "github.com/zviedris/portainerexport/model"
	processing "github.com/zviedris/portainerexport/processing"
)

func main() {

	configFile, err := os.Open("config.json")
	if err != nil {
		// Handle error
		fmt.Println("Config file not found")
		return
	}
	defer configFile.Close()

	var config model.Config
	err = json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		fmt.Println("Could not parse config file!")
		return

	}

	//do request and get info
	results := processing.ProcessPortainer(&config)

	//export results to Excel
	processing.ExportExcel(&results)

}
