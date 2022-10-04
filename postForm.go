package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Pays2 struct {
	Project     int    `json:"project"`
	Invoice     int    `json:"invoice"`
	Status      string `json:"status"`
	Amount      int    `json:"amount"`
	Amount_paid int    `json:"amount_paid"`
	Rand        string `json:"rand"`
}

func main() {
	getToCall()
}

func tojson() []byte {
	pay := Pays2{
		Project:     816,
		Invoice:     73,
		Status:      "completed",
		Amount:      700,
		Amount_paid: 700,
		Rand:        "SNuHufEJ",
	}

	data, err := json.Marshal(pay)
	if err != nil {
		fmt.Println(err)

	}
	fmt.Println(string(data))

	return data
}

func getToCall() {

	data := tojson()
	req, err := http.NewRequest("POST", "http://127.0.0.1:8080/callback_url/", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "multipart/form-data")
	req.Header.Set("Authorization", "d84eb9036bfc2fa7f46727f101c73c73")

	if err != nil {
		fmt.Println(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
	}

	readData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(readData))

	defer resp.Body.Close()

}
