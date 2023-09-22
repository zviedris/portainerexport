package processing

import (
	"fmt"
	"sort"

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
