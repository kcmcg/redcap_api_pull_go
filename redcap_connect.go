package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"bytes"
)

func sendREDCapRequest(params map[string]string) (map[string]interface{},error) {
	requestBody := ""

	for param,value := range params {
		if requestBody != "" {
			requestBody += "&"
		}
		requestBody += param + "=" + value
	}
	requestUrl := "https://redcap.vanderbilt.edu/api/"

	bodyReader := bytes.NewReader([]byte(requestBody))

	resp, err := http.Post(requestUrl, "application/x-www-form-urlencoded", bodyReader)
	if err != nil {
		fmt.Println(err)
		return nil,err
	}

	fmt.Println("Received a response!")
	
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	
	fmt.Println(string(body))

	var responseString map[string]interface{}

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
	fmt.Println(responseString)
}
