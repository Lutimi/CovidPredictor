package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func main() {
	preprocessData()
}

func preprocessData() (err error) {
	err = readCSV()
	if err != nil { return err }

	filterData()
	return nil
}

func readCSV() error {
	// get URL form RawLink
	urlDataset := "https://raw.githubusercontent.com/Lutimi/dataset/master/pm21Septiembre2021.csv"

	// name of new file 
	filename := "pm21Setiembre2021.csv"

	// get Data from RawLink
	resp, err := http.Get(urlDataset)
	if err != nil { return err }
	defer resp.Body.Close()

	// validation status 
	if resp.StatusCode != 200 { return errors.New("recieved non 200 response code fetching dataset")}

	// create new file
	file, err := os.Create(filename)
	if err != nil { return err }

	// move data in new file
	_, err = io.Copy(file, resp.Body)
	if err != nil { return err }

	defer file.Close()
	return nil
}

func filterData() ([][]string, error) {
	// open file
	csvData, err := os.Open("pm21Septiembre2021.csv")
	if err != nil { return nil, err }
	defer csvData.Close()

	reader := csv.NewReader(csvData)
	reader.Comma = '|'
	reader.LazyQuotes = true
	data, err := reader.ReadAll()
	if err != nil { return nil, err }

	// select data of Dataset
	newData := [][]string{}
	line := []string{}
	for column, row := range data {
		if row[14] == "POSITIVO" { 
			line = append(line, data[column][2])
			line = append(line, data[column][10])
			line = append(line, data[column][14])
			newData = append(newData, line)
			line = nil
		}
	}

	// delete duplicate data
	sampleData := [][]string{}
	diferent := false
	complete := true
	for complete {
		content := len(sampleData)
		if len(sampleData) == 0 { sampleData = append(sampleData, newData[0]) 
		} else {
			for _, data := range newData {
				for _, sample := range sampleData {
					if data[0] != sample[0] && data[1] != sample[1] { diferent = true 
					} else if data[0] == sample[0] && data[1] == sample[1] { 
						diferent = false
						break 
					} 
				}
				if diferent { 
					sampleData = append(sampleData, data)
					diferent = false
					break
				}
			}
		} 
		if content != len(sampleData) { complete = true 
		} else { complete = false }
	}

	// add COVID-19 cases
	mainData := [][]string{}
	join := []string{}
	for _, row := range sampleData {
		join = append(join, row[0])
		join = append(join, row[1])
		join = append(join, row[2])
		join = append(join, "0")
		mainData = append(mainData, join)
		join = nil
	}
	count := 0
	for _, cases := range mainData {
		for _, oldData := range newData {
			if cases[0] == oldData[0] && cases[1] == oldData[1] { count++ }
		}
		cases[3] = strconv.Itoa(count)
		count = 0
	}

	// view
	for _, view := range mainData { fmt.Println(view) }

	return mainData, nil
}