package main

import (
	//"bufio"
	"fmt"
	"net"
	"encoding/json"
)
var direccion string = "localhost:9101"

type PredictionRequest struct {
	Region string `json:"region"`
	Date   string `json:"date"`
}

type Coefs struct {
	b float64 `json:"b"`
	m float64 `json:"m"`
}



func main(){
	enviar(direccion,PredictionRequest{"LIMA","20211225"})
}




func enviar(direccion string, info PredictionRequest){
	con,_ := net.Dial("tcp",direccion)
	defer con.Close()

	jsonByte,_ := json.Marshal(info)
	fmt.Fprintln(con,string(jsonByte))
}