package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aanelli/cf-metrics/flatten"
)

func printAsJSON(fileName string, data interface{}) error {
	output, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("error creating file")
		return err
	}
	defer file.Close()

	bytesWritten, err := file.Write(output)
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Printf("Wrote %d bytes.\n", bytesWritten)
	return nil
}

//https://github.com/360EntSecGroup-Skylar/excelize
func printAsCSV(fileName string, data []cfData) error {
	outputCSV := [][]string{}

	for _, datapoint := range data {
		jsonBytes, err := json.Marshal(datapoint)
		if err != nil {
			return err
		}

		stringyJSON := string(jsonBytes[:])
		flatData, err := flatten.FlattenString(stringyJSON, "")
		if err != nil {
			return nil
		}
		outputCSV = append(outputCSV, []string{flatData})
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
