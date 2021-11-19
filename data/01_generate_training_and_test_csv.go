package main

import (
	"encoding/csv"
	"log"
	"math/rand"
	"os"
	"strconv"
)

//Definir arrays por cada region

//Generar dataset de entrenamiento y dataset de test
func main() {
	// Abrir dataset
	f, err := os.Open("dataset.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	//Crear csvReader y establecer numero de columnas
	covidData := csv.NewReader(f)
	covidData.Comma = '|'
	covidData.FieldsPerRecord = 15

	//Lee todos los registros
	records, err := covidData.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	
	// filtramos y agrupamos la data
	records,_ = filterData(records)

	//ALBERT: REPARTIR LA DATA EN LOS ARRAYS------------------------------------------------->

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