package main

import (
	//"bufio"
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
)
var direccion string = "localhost:9102"
var Node3direccion string = "localhost:9103"


type CoefsSend struct {
	Bcoef float64 `json:"b"`
	Mcoef float64 `json:"m"`
	Region string `json:"region"`
	Date   string `json:"date"`
}

type Data struct {
	Region string `json:"region"`
	Date   string `json:"date"`
	NCases int    `json:"nCases"`
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

	var info CoefsSend 
	json.Unmarshal([]byte(bytesInfo),&info)
	fmt.Println(info)

	data := calculatePrediction(info.Region,info.Date,info.Bcoef,info.Mcoef)

	enviar(Node3direccion,data)
}

func calculatePrediction(region string, date string, b float64, m float64) (prediction Data){
	x,_ := strconv.Atoi(date)
	nCases := b + m*float64(x)
	return Data{region, date,int(nCases)}
}

func enviar(direccion string, info Data){
	con,_ := net.Dial("tcp",direccion)
	defer con.Close()
	fmt.Println(info)
	jsonByte,_ := json.Marshal(info)
	fmt.Fprintln(con,string(jsonByte))
}