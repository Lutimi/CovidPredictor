package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

const filename = "database.csv"

type Data struct {
	Region string `json:"region"`
	Date   string `json:"date"`
	NCases int    `json:"nCases"`
}

func handleRequest() {
	http.HandleFunc("/save-prediction", save_prediction)
	http.HandleFunc("/list-predictions", list_predictions)

	log.Fatal(http.ListenAndServe(":9000", nil))
}

func main() {

	//only do this when running the app for the first time

	init_csv_file(filename)

	//if running program for the first time
	//file = init_csv_file(filename)
	//write_csv_headers(file)

	handleRequest()

}

func save_prediction(response http.ResponseWriter, request *http.Request) {
	//response.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE")
	
	if request.Method == "POST" {
		response.Header().Set("Content-Type", "application/json")
		response.Header().Set("Access-Control-Allow-Origin", "*")
	
		if request.Header.Get("Content-Type") == "application/json" {
			jsonBytes, err := ioutil.ReadAll(request.Body)
			if err != nil {
				http.Error(response, "Problem reading JSON", http.StatusInternalServerError)
			} else {
				//unmarshal incoming data
				var data Data
				json.Unmarshal(jsonBytes, &data)
				//open the csv file
				outputFile, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
				if err != nil {
					log.Fatal(err)
				}
				//write the prediction to csv
				write_to_csv(outputFile, &data)
				defer outputFile.Close()

				//send success msg
				io.WriteString(response, `
					{
						"message":"Operation Successful"
					}
				`)
			}
		}
	} else if request.Method == "OPTIONS" {
		//response.Header().Set("Content-Type", "application/json")
		//: 
		//Access-Control-Request-Method: POST	
		response.Header().Set("Access-Control-Allow-Headers","content-type")
		response.Header().Set("Access-Control-Allow-Origin", "*")
		response.Header().Set("Access-Control-Allow-Methods", "POST")
		//response.Header().Set("Access-Control-Request-Headers", "content-type")
		//response.Header().Del("Content-Type")
		//response.Header().Set("Access-Control-Allow-Origin", "*")
		//io.WriteString(response, "Wait")
	} else {
		http.Error(response, "invalid method", http.StatusMethodNotAllowed)
	}
}

//GET
func list_predictions(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	response.Header().Set("Access-Control-Allow-Origin", "*")

	//read the data from csv
	predictions := read_data_from_csv(filename)

	//if we want to filter by param
	var filtered_predictions []Data

	region_param := request.URL.Query().Get("region")
	if region_param != "" {

		//save entries with matching regions in the new slice
		for _, pred := range predictions {
			if pred.Region == region_param {
				filtered_predictions = append(filtered_predictions, pred)
			}
		}
		//assign the new slice to return
		predictions = filtered_predictions
	}
	//otherwise send the unfiltered one
	jsonBytes, _ := json.MarshalIndent(predictions, "", " ")
	io.WriteString(response, string(jsonBytes))
}

func init_csv_file(filename string) (file *os.File) {
	//check if file exists if not creeate
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	return file
}

//only call this once to write the headers
func write_csv_headers(file *os.File) {
	writer := csv.NewWriter(file)
	var row []string
	row = append(row, "region")
	row = append(row, "date")
	row = append(row, "ncases")
	writer.Write(row)

	writer.Flush()
}

//loop over the data struct and write as csv row
func write_to_csv(file *os.File, data *Data) {

	writer := csv.NewWriter(file)
	var row []string
	row = append(row, data.Region)
	row = append(row, data.Date)
	row = append(row, strconv.Itoa(data.NCases))
	fmt.Println(row)

	err := writer.Write(row)
	if err != nil {
		log.Fatal(err)
	}

	writer.Flush()
	err = writer.Error()
	if err != nil {
		log.Fatal(err)
	}
}

func read_data_from_csv(filename string) []Data {
	//open file
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	//create csv reader
	csvReader := csv.NewReader(f)
	lines, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	// convert csv lines to array of structs
	var data_list []Data

	for i, line := range lines {
		if i > 0 { // omit header line
			var pred Data
			for j, field := range line {
				if j == 0 {
					//first column in csv is region
					pred.Region = field
				} else if j == 1 {
					//second column in csv is date
					pred.Date = field
				} else if j == 2 {
					var err error
					//third column in csv is number of cases
					pred.NCases, err = strconv.Atoi(field)
					if err != nil {
						log.Fatal(err)
						continue
					}
				}
			}
			//append to list and return
			data_list = append(data_list, pred)
		}
	}
	return data_list
}
