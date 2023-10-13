package processing

import (
	"fmt"
	"os"
	"sort"
	"strings"

	xlsx "github.com/tealeg/xlsx"
	model "github.com/zviedris/portainerexport/model"
)

func ExportExcel(dataPtr *map[string][]model.EnvVersion) {

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

	headerRow := sheet.AddRow()
	headerRow.AddCell().SetString("Image")
	headerRow.AddCell().SetString("Environment")
	headerRow.AddCell().SetString("Version")
	headerRow.AddCell().SetString("Stack")
	headerRow.AddCell().SetString("Image url")

	// Write data from the struct to the Excel file
	for _, key := range keys {
		dataRow := sheet.AddRow()
		dataRow.AddCell().SetString(key)
		containers := data[key]
		for _, version := range containers {
			dataRow1 := sheet.AddRow()
			dataRow1.AddCell().SetString(" ")
			dataRow1.AddCell().SetString(version.Environment)
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

func ExportMarkdown(dataPtr *map[string][]model.EnvVersion) {
	file, err := os.Create("output.md")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()
	data := *dataPtr
	// Convert map keys into a slice
	var keys []string
	for key := range data {
		keys = append(keys, key)
	}

	// Sort the slice of keys in alphabetical order
	sort.Strings(keys)

	var builder strings.Builder

	builder.WriteString("| Image   | Environment | Version | Stack |Image url     |\n")
	builder.WriteString("|---------|-------------|---------|-------|--------------|\n")

	// Write data from the struct to the Excel file
	for _, key := range keys {
		builder.WriteString(fmt.Sprintf("|%s|||||\n", key))
		containers := data[key]
		for _, version := range containers {
			builder.WriteString(fmt.Sprintf("| |%s|%s|%s|%s|\n", version.Environment, version.Docker, version.Stack, version.DockerPath))
		}

	}

	// Write the content to the file
	_, err = file.WriteString(builder.String())
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("File written successfully")
}
