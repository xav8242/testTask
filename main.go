package main

import (
	"crypto/md5"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const merchant_id = 6
const merchant_key = "KaTf5tZYHx4v7pgZ"
const app_id = 816
const app_key = "rTaasVHeteGbhwBx"

func main() {

	Handle()

}

func Handle() {

	http.HandleFunc("/callback_url/", callbackUrl)

	log.Println("Listen server")
	http.ListenAndServe(":8080", nil)
}

func callbackUrl(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}

	switch r.Header.Get("Content-Type") {

	case "application/json":

		data := umJson(body)

		inserSQL(data)

		str := cryptSHA256(sortAndJoin(data, ":", string(merchant_id)))

		fmt.Printf("Контрольная сумма %v\n", str)

	case "multipart/form-data":
		data := umJson(body)

		inserSQL(data)

		auth := r.Header.Get("Authorization")
		data["Authorization"] = auth

		str := cryptMD5(sortAndJoin(data, ".", app_key))
		fmt.Printf("Контрольная сумма %v\n", str)
	}

}

func umJson(data []byte) map[string]interface{} {

	pay := make(map[string]interface{})

	err := json.Unmarshal(data, &pay)
	if err != nil {
		log.Printf("%s\n", err)
	}
	return pay
}

func sortAndJoin(data map[string]interface{}, sep string, key string) string {
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	str := ""

	for _, k := range keys {
		if k == "sign" {
			continue
		}
		str = str + fmt.Sprintf("%v", data[k]) + sep
	}
	str = str + key

	return str

}

func cryptSHA256(str string) string {
	h := sha256.Sum256([]byte(str))

	str = fmt.Sprintf("%x", h)
	return str
}
func cryptMD5(str string) string {
	h := md5.Sum([]byte(str))
	return fmt.Sprintf("%x", h)
}

func inserSQL(data map[string]interface{}) {

	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	ins := "INSERT INTO pays( "
	val := ") VALUES ("

	for in, k := range keys {
		ins = ins + k

		val = val + "'" + fmt.Sprintf("%v", data[k]) + "'"
		if in != len(keys)-1 {
			ins = ins + ","
			val = val + ","
		}

	}
	query := ins + val + ")"
	query = strings.ReplaceAll(query, "invoice", "payment_id")
	query = strings.ReplaceAll(query, "project", "merchant_id")

	conn, err := sql.Open("mysql", "user:user@tcp(127.0.0.1:3306)/pay")
	if err != nil {
		fmt.Printf("Соединение с БД не установлено: %s\n", err)
	}
	defer conn.Close()

	insert, err := conn.Query(query)
	if err != nil {
		fmt.Printf("Ошибка выполнения запроса к БД: %s\n", err)
	}
	defer insert.Close()

}

func ParseDate(tempTime string) string {

	template := "2006-01-02 15:04:05"
	layout := "15:04:05 02-01-06"
	t, err := time.Parse(template, string(tempTime))
	if err != nil {
		fmt.Println(err)

	}
	return t.Format(layout)
}
