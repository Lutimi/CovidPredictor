package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"github.com/sajari/regression"
)

// Definir arrays por cada region
var regions = []string{"AMAZONAS","ANCASH","APURIMAC","AREQUIPA","AYACUCHO","CAJAMARCA","CUSCO","HUANCAVELICA","HUANUCO","ICA","JUNIN","LA LIBERTAD","LAMBAYEQUE","LIMA","LORETO","MADRE DE DIOS","MOQUEGUA","PASCO","PIURA","PUNO","SAN MARTIN","TACNA","TUMBES","UCAYALI"} // Contiene todas las regiones
var dataPerRegion = [][]string{} //Contiene los datos de contagios por fecha separado por region
var formulas = [][]string{}
// Generar dataset de entrenamiento y dataset de test
func main() {

	// Obtener todos los datos del dataset
	
	//uri := "https://media.githubusercontent.com/media/Gonzarod/dataset/master/pm21Septiembre2021.csv"
	//uri := "https://raw.githubusercontent.com/Lutimi/dataset/master/pm21Septiembre2021.csv"
	//covidData,_ := readCSVFromUrl(uri)	
	covidData := getCovidDataFromCSV()
	// filtramos y agrupamos por dia y region la data 
	//filterData, _ := filterData(covidData)

	// Se asigna valor a dataPerRegion
	//makeDataSetPerRegion(filterData)
	makeDataSetPerRegion(covidData)
	
	//regionData := getDatasetRegion(3)
	//fmt.Print(dataPerRegion)

	setDataSetPerRegion(covidData)

	

	for _, region := range regions { 
		fmt.Printf("Training for %s ----------\n",region)
		n := getIndexOfRegion(region)
		dataSetRegion := getDatasetRegion(n)
		fmt.Printf("Size: %d\n",len(dataSetRegion))
		trainModel(dataSetRegion,region)
	}

    //fmt.Println(regionData);

	//---write formulas
	writeFormulas();

}

func getCovidDataFromCSV() ([][] string){
	dataset, err := os.Open("clean_dataset.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer dataset.Close()

	// Crear csvReader y establecer numero de columnas
	covidDataReader := csv.NewReader(dataset)
	covidDataReader.Comma = ';'

	// Lee todos los registros
	records, err := covidDataReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	return records
}

func getDatasetRegion(position int) ([][] string) {
	regionData := [][]string{}
	headers := []string{"Date", "Region", "CovidCases"}
	regionData = append(regionData, headers)
	join := []string{}
	for x := range dataPerRegion[position] {
		if x == len(dataPerRegion[position]) - 3 { break }
		if x % 3 == 0 {
			join = append(join, dataPerRegion[position][x])
			join = append(join, dataPerRegion[position][x + 1])
			join = append(join, dataPerRegion[position][x + 2])
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
		content := len(regions)
		if len(regions) == 0 { regions = append(regions, newData[0][1]) 
		} else {
			for _, data := range newData {
				for dx := range regions {
					if data[1] != regions[dx] { diverse = true 
					} else if data[1] == regions[dx] { 
						diverse = false
						break 
					} 
				}
				if diverse { 
					regions = append(regions, data[1])
					diverse = false
					break
				}
			}
		} 
		if content != len(regions) { finish = true 
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
				for _, dataPerRegion := range modelData {
					if data[0] != dataPerRegion[0] && data[1] != dataPerRegion[1] { diferent = true 
					} else if data[0] == dataPerRegion[0] && data[1] == dataPerRegion[1] { 
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

func makeDataSetPerRegion(records [][]string) {
	datsetRegion := []string{}
	
	for dx := range regions {
		for _, training := range records {
			if regions[dx] == training[1] { 
				datsetRegion = append(datsetRegion, training...) 
			} else if regions[dx] != training[1] { continue }
		}
		dataPerRegion = append(dataPerRegion, datsetRegion)
		datsetRegion = nil
	}
}

func getIndexOfRegion(value string) int {
    for p, v := range regions {
        if (v == value) {
            return p
        }
    }
    return -1
}

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


	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(resp.Body)
	reader.Comma = '|'
	reader.FieldsPerRecord = 15
	data := [][]string{}

	for {
        rec, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            //log.Fatal(err)
        } else {
			data = append(data, rec)
		}
	}

	if err != nil {
		fmt.Println(err)
	}



	fmt.Print(data[0])
	fmt.Print(data[1])

	return data, nil
}

func writeFormulas() {

	csvFile, err := os.Create("regression_formulas.csv")

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	csvwriter := csv.NewWriter(csvFile)
	
	for _, r := range formulas {
		if len(r) > 0 {
			_ = csvwriter.Write(r)
		}
	}

	csvwriter.Flush()
	csvFile.Close()

}

func setDataSetPerRegion(data [][]string){
	datsetRegion := []string{}
	
	for dx := range regions {
		for _, training := range data {
			if regions[dx] == training[1] { 
				datsetRegion = append(datsetRegion, training...) 
			} else if regions[dx] != training[1] { continue }
		}
		dataPerRegion = append(dataPerRegion, datsetRegion)
		datsetRegion = nil
	}
}