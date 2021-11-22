package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var formulas [][]string
var NodesSwitch bool = false
var chanData chan Data

var direccion string = "localhost:9066"
var direccionNode string = "localhost:9101"

type Data struct {
	Region string `json:"region"`
	Date   string `json:"date"`
	NCases int    `json:"nCases"`
}

type PredictionRequest struct {
	Region string `json:"region"`
	Date   string `json:"date"`
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
	var prediction Data
	if NodesSwitch == true {
		enviar(direccionNode,PredictionRequest{region_param,date_param})
		predictionChan := <- chanData;
		prediction = predictionChan
	} else {
		prediction = calculatePrediction(region_param,date_param)
	}
	

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


func listenToNode() {
	ln, _ := net.Listen("tcp",direccion)

	defer ln.Close()

	for {
		con, _ := ln.Accept()
		go manejadorConeXion(con)
	}
}

func manejadorConeXion(con net.Conn)  {

	defer con.Close()
	fmt.Println("Recibo de nodo")
	bufferIn := bufio.NewReader(con)
	bytesInfo, _ := bufferIn.ReadString('\n')

	var data Data 
	json.Unmarshal([]byte(bytesInfo),&data)
	fmt.Println(data)
	chanData <- data
	//Commentar para prueba local:
	

}
func enviar(direccion string, info PredictionRequest){
	con,_ := net.Dial("tcp",direccion)
	defer con.Close()
	fmt.Println(info)
	jsonByte,_ := json.Marshal(info)
	fmt.Fprintln(con,string(jsonByte))
}



func main() {
	chanData=make(chan Data)
	go listenToNode();
	loadFormulas();
	handleRequest()
}
