package redcap

import (
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
	"database/sql"
	"time"
	_ "github.com/go-sql-driver/mysql"
)

type Credentials struct {
	Host string
	Port string
	Db string
	User string
	Password string
}

func getCredentials() (Credentials, error) {
	creds := Credentials{}
	jsonFile, err := os.Open("../credentials/redcap_db.json")
	if err != nil {
		fmt.Println(err)
		return creds,err
	}

	byteValue, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		fmt.Println(err)
		return creds,err
	}
	fmt.Println(string(byteValue))
	json.Unmarshal([]byte(byteValue), &creds)

	defer jsonFile.Close()

	fmt.Println(creds)

	return creds,nil
}

func ConnectToDb() (*sql.DB, error) {
	params,err := getCredentials()
	if err != nil {
		return nil,err
	}
	fmt.Println(params) 
	connString := params.User + ":" + params.Password + "@" + params.Host + ":" + params.Port + "/" + params.Db
	fmt.Println(connString)

	db, err := sql.Open("mysql", connString)
	if err != nil {
		fmt.Println(err)
		return db,err
	}
	db.SetConnMaxLifetime(time.Minute * 30)

	return db,nil
}
