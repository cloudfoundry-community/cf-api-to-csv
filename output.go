package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aanelli/cf-metrics/flatten"
)

func printAsJSON(fileName string, data interface{}) error {
	output, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("error creating file", err)
		return err
	}
	defer file.Close()

	bytesWritten, err := file.Write(output)
	if err != nil {
		fmt.Println("error writing to file", err)
		return err
	}
	fmt.Printf("Wrote %d bytes.\n", bytesWritten)
	return nil
}

//https://github.com/360EntSecGroup-Skylar/excelize
func printAsCSV(fileName string, data []cfData) error {
	outputCSV := [][]string{}

	for _, datapoint := range data {
		outputCSV = append(outputCSV, []string{datapoint.Name, datapoint.GUID, datapoint.OrganizationGUID})

		outputCSV = append(outputCSV, []string{"\n"})
		outputCSV = append(outputCSV, []string{"APPS"})
		for _, app := range datapoint.Apps {
			temp, err := convertCFAPIResourceToCSVString(app)
			if err != nil {
				return err
			}
			outputCSV = append(outputCSV, temp)
		}

		outputCSV = append(outputCSV, []string{"\n"})
		outputCSV = append(outputCSV, []string{"APP CREATES"})
		for _, appCreate := range datapoint.AppCreates {
			temp, err := convertCFAPIResourceToCSVString(appCreate)
			if err != nil {
				return err
			}
			outputCSV = append(outputCSV, temp)
		}

		outputCSV = append(outputCSV, []string{"\n"})
		outputCSV = append(outputCSV, []string{"APP STARTS"})
		for _, appStart := range datapoint.AppStarts {
			temp, err := convertCFAPIResourceToCSVString(appStart)
			if err != nil {
				return err
			}
			outputCSV = append(outputCSV, temp)
		}

		outputCSV = append(outputCSV, []string{"\n"})
		outputCSV = append(outputCSV, []string{"APP UPDATES"})
		for _, appUpdate := range datapoint.AppUpdates {
			temp, err := convertCFAPIResourceToCSVString(appUpdate)
			if err != nil {
				return err
			}
			outputCSV = append(outputCSV, temp)
		}

		outputCSV = append(outputCSV, []string{"\n"})
		outputCSV = append(outputCSV, []string{"SPACE CREATES"})
		for _, spaceCreate := range datapoint.SpaceCreates {
			temp, err := convertCFAPIResourceToCSVString(spaceCreate)
			if err != nil {
				return err
			}
			outputCSV = append(outputCSV, temp)
		}

		outputCSV = append(outputCSV, []string{"\n"})
		outputCSV = append(outputCSV, []string{"SERVICE BINDINGS"})
		for _, serviceBinding := range datapoint.ServiceBindings {
			temp, err := convertCFAPIResourceToCSVString(serviceBinding)
			if err != nil {
				return err
			}
			outputCSV = append(outputCSV, temp)
		}
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	for _, value := range outputCSV {
		//fmt.Println("writing value: ", value, "to file"+"\n\n\n\n\n\n\n\n\n")
		err := writer.Write(value)
		if err != nil {
			return err
		}
	}
	return nil
}

func convertCFAPIResourceToCSVString(resource cfAPIResource) ([]string, error) {
	//turn the interface back into json
	jsonBytes, err := json.Marshal(resource.Entity)
	if err != nil {
		return nil, err
	}

	//go from json to string
	stringyJSON := string(jsonBytes[:])

	//flatten the json, and separate by comma
	flatData, err := flatten.FlattenString(stringyJSON, "")
	if err != nil {
		return nil, nil
	}
	flatSlice := strings.Split(flatData, ",")
	flatSlice = append(flatSlice, resource.Metadata.CreatedAt.String(), resource.Metadata.GUID, resource.Metadata.UpdatedAt.String(), resource.Metadata.URL)

	return flatSlice, nil
}

func printProgressBar(iteration int, total int, prefix string, suffix string, decimals int, length int) {

}
