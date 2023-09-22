package main

import (
	"crypto/tls"
	json "encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	os "os"
	"sort"
	"strings"

	xlsx "github.com/tealeg/xlsx"
	"github.com/zviedris/portainerexport/models"
)

func exportExcel(dataPtr *map[string][]models.EnvVersion) {

	data := *dataPtr
	//fmt.Println(data)
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		fmt.Printf("Error creating sheet: %v\n", err)
		return
	}

	// Convert map keys into a slice
	var keys []string
	for key := range data {
		keys = append(keys, key)
	}

	// Sort the slice of keys in alphabetical order
	sort.Strings(keys)

	// Write data from the struct to the Excel file
	for _, key := range keys {
		dataRow := sheet.AddRow()
		dataRow.AddCell().SetString(key)
		containers := data[key]
		for _, version := range containers {
			dataRow1 := sheet.AddRow()
			dataRow1.AddCell().SetString(" ")
			dataRow1.AddCell().SetString(version.Enviornment)
			dataRow1.AddCell().SetString(version.Docker)
			dataRow1.AddCell().SetString(version.Stack)
			dataRow1.AddCell().SetString(version.DockerPath)
		}

	}

	// Save the Excel file
	err = file.Save("output.xlsx")
	if err != nil {
		fmt.Printf("Error saving Excel file: %v\n", err)
		return
	}

}

func main() {

	configFile, err := os.Open("config.json")
	if err != nil {
		// Handle error
	}
	defer configFile.Close()

	var config models.Config
	err = json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		// Handle error
	}
	results := map[string][]models.EnvVersion{}

	for _, env := range config.Enviornments {

		// Define the URL you want to send the GET request to
		url := env.Url

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // <--- Problem
		}

		// Create a new HTTP client
		client := &http.Client{Transport: tr}

		//for each stack
		//for _, stack := range config.Stacks {
		// Create a new GET request with query parameter
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}

		// Set the API key in the request header
		req.Header.Add("X-API-Key", env.ApiKey)
		//fmt.Println(stack.Name)
		// Add query parameter
		//q := req.URL.Query()

		//filterValue := "{\"label\":[\"com.docker.stack.namespace=" + stack.Name + "\"]}"

		//q.Add("filters", filterValue)
		//req.URL.RawQuery = q.Encode()

		// Send the GET request
		response, err := client.Do(req)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer response.Body.Close()

		// Read the response body
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}

		//fmt.Println("Response Body:")
		//fmt.Println(string(body))
		//fmt.Println()

		var portList []models.PortObject

		err = json.Unmarshal(body, &portList)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		for _, cont := range portList {
			//key := cont.Spec.Name
			excludeVal := false
			for _, excl := range config.Exclude {
				if strings.Contains(cont.Spec.Labels.Image, excl.Name) {
					excludeVal = true
				}
			}
			if !excludeVal {
				images := strings.Split(cont.Spec.Labels.Image, ":")
				path := strings.Split(images[0], "/")
				key := path[len(path)-1]
				var value models.EnvVersion
				value.Enviornment = env.Name
				value.Docker = images[1]
				value.Stack = cont.Spec.Labels.Namespace
				value.DockerPath = images[0]
				existingVal := results[key]
				existingVal = append(existingVal, value)
				results[key] = existingVal
			}

		}

		//}
	}

	exportExcel(&results)

}
