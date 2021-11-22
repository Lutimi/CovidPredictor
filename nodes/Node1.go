package main

import (
	//"bufio"
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)
var direccion string = "localhost:9101"
var Node2direccion string = "localhost:9102"

type PredictionRequest struct {
	Region string `json:"region"`
	Date   string `json:"date"`
}

type CoefsSend struct {
	Bcoef float64 `json:"b"`
	Mcoef float64 `json:"m"`
	Region string `json:"region"`
	Date   string `json:"date"`
}

func main(){

	ln, _ := net.Listen("tcp",direccion)

	defer ln.Close()

	for {
		con, _ := ln.Accept()
		go manejadorConeXion(con)
	}

	//enviar(Node1direccion,PredictionResponse{"Lima","20210215",12})
}

func manejadorConeXion(con net.Conn) {
	defer con.Close()

	bufferIn := bufio.NewReader(con)
	bytesInfo, _ := bufferIn.ReadString('\n')

	var info PredictionRequest 
	json.Unmarshal([]byte(bytesInfo),&info)
	fmt.Println(info)

	b,m := getCoeff(info.Region)
	fmt.Println(b)
	fmt.Println(m)
	//Enviar siguiente node : Nodo 2
	enviar(Node2direccion,CoefsSend{b,m,info.Region,info.Date})
}


func enviar(direccion string, info CoefsSend){
	con,_ := net.Dial("tcp",direccion)
	defer con.Close()
	fmt.Println(info)
	jsonByte,_ := json.Marshal(info)
	fmt.Fprintln(con,string(jsonByte))
}


func getCoeff(region string) (b float64,m float64){
	region = strings.ToUpper(region)
	formulas := loadFormulas()

	for _,a := range formulas {
		if a[0]==region {
			b, _ := strconv.ParseFloat(a[1], 64)
			m, _ := strconv.ParseFloat(a[2], 64)
			return b,m 
		}
	}
	return float64(0), float64(0)
}

func loadFormulas() ([][]string){
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

	return records
}