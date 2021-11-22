package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"encoding/csv"
	"os"
	"fmt"
	"strconv"
)

var formulas [][]string

type Data struct {
	Region string `json:"region"`
	Date   string `json:"date"`
	NCases int    `json:"nCases"`
}

func handleRequest() {
	http.HandleFunc("/make-prediction", make_prediction)

	log.Fatal(http.ListenAndServe(":9001", nil))
}

func init_mock_data() Data {
	pred := Data{"Lima", "20221102", 772215}
	return pred
}

func make_prediction(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	response.Header().Set("Access-Control-Allow-Origin", "*")

	region_param := request.URL.Query().Get("region")
	if region_param == "" {
		http.Error(response, "missing region param", http.StatusBadRequest)
		return
	}

	date_param := request.URL.Query().Get("date")
	if date_param == "" {
		http.Error(response, "missing date param", http.StatusBadRequest)
		return
	}

	//get data from prediction
	//pass region param and date param
	//prediction := init_mock_data()
	prediction := calculatePrediction(region_param,date_param)

	jsonBytes, _ := json.MarshalIndent(prediction, "", " ")
	io.WriteString(response, string(jsonBytes))

}

func calculatePrediction(region string, date string) (prediction Data){
	region = strings.ToUpper(region)
	b,m := getCoeff(region);
	x,_ := strconv.Atoi(date)
	nCases := b + m*float64(x)
	return Data{region, date,int(nCases)}
}
func getCoeff(region string) (b float64,m float64){
	
	for _,a := range formulas {
		if a[0]==region {
			b, _ := strconv.ParseFloat(a[1], 64)
			m, _ := strconv.ParseFloat(a[2], 64)
			return b,m 
		}
	}
	return float64(0), float64(0)
}

func loadFormulas(){
	dataset, err := os.Open("regression_formulas.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer dataset.Close()

	// Crear csvReader y establecer numero de columnas
	covidDataReader := csv.NewReader(dataset)
	covidDataReader.Comma = ','

	// Lee todos los registros
	records, err := covidDataReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	formulas = records

	fmt.Println(formulas)
}

func main() {
	loadFormulas();
	handleRequest()
}
