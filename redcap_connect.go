package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"bytes"
)

func main() {
	jsonFile, err := os.Open("../credentials/REDCAP_ZIP_IMPORT_API.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("successfully Opened json file")

	var token string

	byteValue, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		fmt.Println(err)
		return
	}

	var result map[string]string
	json.Unmarshal([]byte(byteValue), &result)
	token = result["api_token"]

	defer jsonFile.Close()

	fmt.Println(token)

	requestUrl := "https://redcap.vanderbilt.edu/api/"

	jsonBody := []byte(`format=json&content=record&token=` + token + `&returnFormat=json`)
	bodyReader := bytes.NewReader(jsonBody)
//	body := make(map[string]string,4)
//
//	body["returnFormat"] = "json"
//	body["action"] = "records"
//	body["api_token"] = token
//
	resp, err := http.Post(requestUrl, "application/x-www-form-urlencoded", bodyReader)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Received a response!")
	
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	
	fmt.Println(string(body))

	var responseString map[string]interface{}

	json.Unmarshal([]byte(body), &responseString)

	fmt.Println(responseString)
}
