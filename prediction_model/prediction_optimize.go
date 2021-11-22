package predictor

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"os"
	"log"

	"github.com/sajari/regression"
)

// Predicted = -47873599.4334 + Date*2.3699

// Definir arrays por cada region
var listRegion = []string{} // Contiene todas las regiones
var dataPerRegion = [][]string{} //Contiene los datos de contagios por fecha separado por region
var formulas = [][]string{}

// Generar dataset de entrenamiento y dataset de test
func predictor() {

	//url := "https://media.githubusercontent.com/media/Gonzarod/dataset/master/pm21Septiembre2021.csv"
	url := "https://raw.githubusercontent.com/Lutimi/dataset/master/pm21Septiembre2021.csv"

	covidData, _ := readCSVFromUrl(url)

	filterData, _ := filterData(covidData)

	// Se asigna valor a dataPerRegion
	makeDataSetPerRegion(filterData)

	// Se entrena el modelo por region
	for n, region := range listRegion {
		fmt.Printf("Training for %s ----------\n", region)
		dataSetRegion := getDatasetbyRegion(n)
		trainModel(dataSetRegion, region)
	}

	writeFormulas();

	//fmt.Println(filterData)

	//for _, view := range filterData { fmt.Println(view) }
}

func getDatasetbyRegion(region int) [][]string {
	regionData := [][]string{}
	regionData = append(regionData, []string{"Date", "Region", "CovidCases"})
	join := []string{}
	for x := range dataPerRegion[region] {
		if x == len(dataPerRegion[region]) - 3 { break }
		if x % 3 == 0 {
			join = append(join, dataPerRegion[region][x])
			join = append(join, dataPerRegion[region][x + 1])
			join = append(join, dataPerRegion[region][x + 2])
			regionData = append(regionData, join)
			join = nil
		}
	}

	return regionData
}

func makeDataSetPerRegion(records [][]string) {
	regions := []string{}
	for dx := range listRegion { regions = append(regions, listRegion[dx]) }
	listRegion = nil

	allKeys := make(map[string]bool)
	for _, item := range regions {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			listRegion = append(listRegion, item)
		}
	}

	datsetRegion := []string{}
	for dx := range listRegion {
		for _, training := range records {
			if listRegion[dx] == training[1] { 
				datsetRegion = append(datsetRegion, training...) 
			} else if listRegion[dx] != training[1] { continue }
		}
		dataPerRegion = append(dataPerRegion, datsetRegion)
		datsetRegion = nil
	}
}

func filterData(data [][]string) ([][]string, error) {
	newData := [][]string{}
	column := []string{}

	for col, row := range data {
		if row[14] == "POSITIVO" {
			column = append(column, data[col][2])
			column = append(column, data[col][10])
			column = append(column, "0")

			newData = append(newData, column) // filter data
			listRegion = append(listRegion, data[col][10]) // create for model region
			column = nil
		}
 	}
	
	// fmt.Println("Complete First Step filter")

	// Se crea un modelo de fecha by region
	mainData := [][]string{}
	mainData = append(mainData, newData[0])
	complete:= true
	for complete {
		diferent := false
		content := len(mainData)
		for _, data:= range newData {
			for _, filter := range mainData {
				if filter[0] != data[0] && filter[1] != data[1] { diferent = true 
				} else if filter[0] == data[0] && filter[1] == data[1]{ 
					diferent = false
					break
				}
			}
			if diferent {
				mainData = append(mainData, data)
				diferent = false
				break
			}
		}
		if len(mainData) != content { complete = true
		} else { complete = false }
	}

	// fmt.Println("Complete Second Step filter")

	count := 0
	for _, cases := range mainData {
		for _, old := range newData {
			if cases[0] == old[0] && cases[1] == old[1] { count++ }
		}
		cases[2] = strconv.Itoa(count)
		count = 0
	}

	// fmt.Println("Finish filter data")

	return mainData, nil
}

/*
func convertStringbyDate(data string) string {
	date := []rune{}
	for px, letter := range data {
		if px < 3 || px == 4 || px > 5 {
			date = append(date, letter)
		} else if px == 3 || px == 5 {
			date = append(date, letter)
			date = append(date, '-')
		}
	}

	return string(date)
} */

func trainModel(records [][]string,region string){
	
	var r regression.Regression
	r.SetObserved("CovidCases")
	r.SetVar(0, "Date")

	// Loop of records in the CSV, adding the training data to the regressionvalue.
	for i, record := range records {
		// Skip the header.
		if i == 0 {
			continue
		}

		// Parse the house price, "y".
		price, err := strconv.ParseFloat(records[i][2], 64)
		if err != nil {
			log.Fatal(err)
		}

		// Parse the grade value.
		grade, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			log.Fatal(err)
		}

		// Add these points to the regression value.
		r.Train(regression.DataPoint(price, []float64{grade}))
	}
	
	// Train/fit the regression model.
	r.Run()
	// Output the trained model parameters.
	
	fmt.Printf("\nRegression Formula:\n%v\n\n", r.Formula)
	fmt.Println(r.GetCoeffs())
	if len(r.GetCoeffs())>0 {
		coef := r.GetCoeffs()
		b := fmt.Sprintf("%f",coef[0])
		m := fmt.Sprintf("%f",coef[1])
		formulaRegion := []string{}
		formulaRegion = append(formulaRegion,region)
		formulaRegion = append(formulaRegion, b)
		formulaRegion = append(formulaRegion, m)
		formulas = append(formulas, formulaRegion)	
	}
}

func readCSVFromUrl(url string) ([][]string, error) {
	data := [][]string{}

	resp, err := http.Get(url)
	if err != nil { return nil, err }

	reader := csv.NewReader(resp.Body)
	reader.Comma = '|'
	reader.FieldsPerRecord = 15

	for {
		rec, err := reader.Read()
		if err == io.EOF { break 
		} else { data = append(data, rec) }
	}

	//if err != nil { fmt.Print(err) }

	//fmt.Print(data[0])
	//fmt.Print(data[1])

	fmt.Println("Finish Read CSV")

	return data, nil
}

func writeFormulas() {
	csvFile, err := os.Create("regressionFormulas.csv")

	if err != nil { log.Fatalf("failed creating file: %s", err)}

	csvWriter := csv.NewWriter(csvFile)

	for _, r := range formulas {
		if len(r) > 0 {
			_ = csvWriter.Write(r)
		}
	}

	csvWriter.Flush()
	csvFile.Close()
} 