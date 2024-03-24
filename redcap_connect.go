package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"strings"
	"localhost/redcap_connect/redcap"
)

func sendREDCapRequest(params map[string]string) (interface{},error) {
	requestBody := ""

	for param,value := range params {
		if requestBody != "" {
			requestBody += "&"
		}
		requestBody += param + "=" + value
	}
	requestUrl := "https://redcap.vanderbilt.edu/api/"

	bodyReader := strings.NewReader(requestBody)

	resp, err := http.Post(requestUrl, "application/x-www-form-urlencoded", bodyReader)
	if err != nil {
		fmt.Println(err)
		return nil,err
	}

	fmt.Println("Received a response!")
	fmt.Println(resp.StatusCode)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	
	var responseString interface{}

	json.Unmarshal([]byte(body), &responseString)

	return responseString,nil
}

func fetchApiToken() (string,error) {
	jsonFile, err := os.Open("../credentials/REDCAP_ZIP_IMPORT_API.json")
	if err != nil {
		fmt.Println(err)
		return "",err
	}
	fmt.Println("successfully Opened json file")

	byteValue, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		fmt.Println(err)
		return "",err
	}

	var result map[string]string
	json.Unmarshal([]byte(byteValue), &result)
	token := result["api_token"]

	defer jsonFile.Close()

	fmt.Println(token)

	return token,nil
}

func main() {
	db, err := redcap.ConnectToDb() 
	pingErr := db.Ping()
	if pingErr != nil {
		fmt.Println(pingErr)
	}

	token, err := fetchApiToken()
	if err != nil {
		return
	}

	requestParams := make(map[string]string)

	requestParams["format"] = "json"
	requestParams["content"] = "record"
	requestParams["token"] = token
	requestParams["returnFormat"] = "json"

//	body := make(map[string]string,4)
//
//	body["returnFormat"] = "json"
//	body["action"] = "records"
//	body["api_token"] = token
//
	responseString,err := sendREDCapRequest(requestParams)
	
	recordList := responseString.([]interface{})
	for _,recordObj := range recordList {
		record := recordObj.(map[string]interface{})
		fmt.Println("Record is",record["record_id"])
	}
}
