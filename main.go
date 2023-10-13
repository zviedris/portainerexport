package main

import (
	json "encoding/json"
	"fmt"
	os "os"

	flag "github.com/spf13/pflag"
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

	var outputFormat *string = flag.String("format", "excel", "Format in which to render output. Supported formats: excel, markdown")
	flag.Parse()

	//do request and get info
	results := processing.ProcessPortainer(&config)

	if outputFormat == nil || *outputFormat == "excel" {
		processing.ExportExcel(&results)
	} else if *outputFormat == "markdown" {
		processing.ExportMarkdown(&results)
	}
}
