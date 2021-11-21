package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
)

// Definir arrays por cada region
var sampleData = []string{}
var model = [][]string{}

// Generar dataset de entrenamiento y dataset de test
func main() {

	// Abrir dataset
	dataset, err := os.Open("dataset.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer dataset.Close()

	// Crear csvReader y establecer numero de columnas
	covidData := csv.NewReader(dataset)
	covidData.Comma = '|'
	covidData.FieldsPerRecord = 15

	// Lee todos los registros
	records, err := covidData.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	
	// filtramos y agrupamos la data
	records, _ = filterData(records)

	// ALBERT: REPARTIR LA DATA EN LOS ARRAYS------------------------------------------------->
	datsetRegion := []string{}
	
	for dx := range sampleData {
		for _, training := range records {
			if sampleData[dx] == training[1] { 
				datsetRegion = append(datsetRegion, training...) 
			} else if sampleData[dx] != training[1] { continue }
		}
		model = append(model, datsetRegion)
		datsetRegion = nil
	}

	// for _, view := range sampleData { fmt.Println(view) }

	// Example get dataset by region
	regionData := [][]string{}
	regionData = getDatasetRegion(3)

	fmt.Print(regionData)
}

func getDatasetRegion(position int) ([][] string) {
	regionData := [][]string{}
	headers := []string{"Date", "Region", "CovidCases"}
	regionData = append(regionData, headers)
	join := []string{}
	for x := range model[position] {
		if x == len(model[position]) - 3 { break }
		if x % 3 == 0 {
			join = append(join, model[position][x])
			join = append(join, model[position][x + 1])
			join = append(join, model[position][x + 2])
			regionData = append(regionData, join)
			join = nil
		}
	}
	return regionData
}

func filterData(data [][]string) ([][]string, error) {

	// select data of Dataset
	newData := [][]string{}
	line := []string{}
	for column, row := range data {
		if row[14] == "POSITIVO" { 
			line = append(line, data[column][2])
			line = append(line, data[column][10])
			newData = append(newData, line)
			line = nil
		}
	}

	// create reference model region
	diverse := false
	finish := true
	for finish {
		content := len(sampleData)
		if len(sampleData) == 0 { sampleData = append(sampleData, newData[0][1]) 
		} else {
			for _, data := range newData {
				for dx := range sampleData {
					if data[1] != sampleData[dx] { diverse = true 
					} else if data[1] == sampleData[dx] { 
						diverse = false
						break 
					} 
				}
				if diverse { 
					sampleData = append(sampleData, data[1])
					diverse = false
					break
				}
			}
		} 
		if content != len(sampleData) { finish = true 
		} else { finish = false }
	}

	// delete duplicate data
	modelData := [][]string{}
	diferent := false
	complete := true
	for complete {
		content := len(modelData)
		if len(modelData) == 0 { modelData = append(modelData, newData[0]) 
		} else {
			for _, data := range newData {
				for _, model := range modelData {
					if data[0] != model[0] && data[1] != model[1] { diferent = true 
					} else if data[0] == model[0] && data[1] == model[1] { 
						diferent = false
						break 
					} 
				}
				if diferent { 
					modelData = append(modelData, data)
					diferent = false
					break
				}
			}
		} 
		if content != len(modelData) { complete = true 
		} else { complete = false }
	}

	// add COVID-19 cases
	mainData := [][]string{}
	join := []string{}
	for _, row := range modelData {
		join = append(join, row[0])
		join = append(join, row[1])
		join = append(join, "0")
		mainData = append(mainData, join)
		join = nil
	}
	count := 0
	for _, cases := range mainData {
		for _, oldData := range newData {
			if cases[0] == oldData[0] && cases[1] == oldData[1] { count++ }
		}
		cases[2] = strconv.Itoa(count)
		count = 0
	}

	// view
	//for _, view := range mainData { f.Println(view) }

	return mainData, nil
}