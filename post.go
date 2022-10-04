package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Pays struct {
	MerchantId int    `json:"merchant_id"`
	PaymentId  int    `json:"payment_id"`
	Status     string `json:"status"`
	Amount     int    `json:"amount"`
	AmountPaid int    `json:"amount_paid"`
	Timestamp  int    `json:"timestamp"`
	Sign       string `json:"sign"`
}

func main() {
	getToCall()
}

func tojson() []byte {
	pay := Pays{
		MerchantId: 6,
		PaymentId:  13,
		Status:     "completed",
		Amount:     500,
		AmountPaid: 500,
		Timestamp:  1654103837,
		Sign:       "f027612e0e6cb321ca161de060237eeb97e46000da39d3add08d09074f931728",
	}

	data, err := json.Marshal(pay)
	if err != nil {
		fmt.Println(err)

	}
	// fmt.Println(string(data))

	return data
}

func getToCall() {

	data := tojson()
	fmt.Printf("%s\n", data)
	req, err := http.NewRequest("POST", "http://127.0.0.1:8080/callback_url/", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	readData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(readData))

	defer resp.Body.Close()

	// url := "http://127.0.0.1:8080/callback_url/?" + data
	// g, err := http.Get(url)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// body, err := ioutil.ReadAll(g.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(string(body))
	// g.Body.Close()
}
