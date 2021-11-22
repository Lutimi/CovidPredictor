package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	//"os"
	//"time"
)

type PredictionResponse struct {
	Region string `json:"region"`
	Date   string `json:"date"`
	NCases int    `json:"nCases"`
}
type Data struct {
	Region string `json:"region"`
	Date   string `json:"date"`
	NCases int    `json:"nCases"`
}

var direccion string = "localhost:9103"
var apiDireccion string = "localhost:9066"
func main(){
	//rol SERVIDOR------------------------------------------------------- A
	ln, _ := net.Listen("tcp",direccion)

	defer ln.Close()

	for {
		con, _ := ln.Accept()
		go manejadorConeXion(con)
	}

}

func manejadorConeXion(con net.Conn)  {

	defer con.Close()
	fmt.Println("Manejo")
	bufferIn := bufio.NewReader(con)
	bytesInfo, _ := bufferIn.ReadString('\n')

	var data Data 
	json.Unmarshal([]byte(bytesInfo),&data)
	fmt.Println(data)

	//Commentar para prueba local:
	enviarMicroML(apiDireccion,PredictionResponse{data.Region,data.Date,data.NCases})

}

func enviarMicroML(direccion string, info PredictionResponse){
	con,_ := net.Dial("tcp",direccion)
	defer con.Close()
	fmt.Println(info)
	jsonByte,_ := json.Marshal(info)
	fmt.Fprintln(con,string(jsonByte))
}