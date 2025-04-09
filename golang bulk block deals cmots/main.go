package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var api_key = os.Getenv("api_key_cmots")

type bulk_response struct {
	Success bool `json:"success"`
	Data    []struct {
		Date       string  `json:"Date"`
		COCODE     float64 `json:"CO_CODE"`
		Scripcode  string  `json:"scripcode"`
		Serial     float64 `json:"Serial"`
		ScripName  string  `json:"ScripName"`
		ClientName string  `json:"ClientName"`
		Buysell    string  `json:"buysell"`
		QTYSHARES  float64 `json:"QTYSHARES"`
		AVGPRICE   float64 `json:"AVG_PRICE"`
	} `json:"data"`
	Message string `json:"message"`
}

var db *sql.DB
var dsn string

func main() {
	var err error
	dsn = os.Getenv("dsn")
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Println(err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(time.Minute * 5)
	defer db.Close()
	log_file, err := os.Create("bulk_bluck_log.log")
	if err != nil {
		fmt.Println(err)
	}
	i := 0
	for i == 0 {
		get_bulk_deals(log_file)
		get_block_deals(log_file)
		get_mutualfunds_holdings(log_file)
		time.Sleep(time.Hour * 24)
	}

}

func get_bulk_deals(log_file *os.File) {
	var unfolded_bulk bulk_response
	url := "https://insbaapis.cmots.com/api/BulkBlockDeal/NSE/Bulk"

	log.SetOutput(log_file)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("Authorization", "Bearer "+api_key)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	res_body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}
	// fmt.Println(string(res_body))
	json.Unmarshal(res_body, &unfolded_bulk)
	// fmt.Println(unfolded_bulk)

	for i := 0; i < len(unfolded_bulk.Data); i++ {
		data := unfolded_bulk.Data[i]
		_, err := db.Exec("call stp_bulk_deals(?,?,?,?,?,?,?,?,?)", data.Date, data.COCODE, data.Scripcode, data.Serial, data.ScripName, data.ClientName, data.Buysell, data.QTYSHARES, data.AVGPRICE)
		if err != nil {
			log.Println(err)
		}
	}

}

func get_block_deals(log_file *os.File) {
	var unfolded_bulk bulk_response
	url := "https://insbaapis.cmots.com/api/BulkBlockDeal/NSE/Block"

	log.SetOutput(log_file)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("Authorization", "Bearer "+api_key)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	res_body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}
	// fmt.Println(string(res_body))
	json.Unmarshal(res_body, &unfolded_bulk)
	// fmt.Println(unfolded_bulk)

	for i := 0; i < len(unfolded_bulk.Data); i++ {
		data := unfolded_bulk.Data[i]
		_, err := db.Exec("call stp_block_deals(?,?,?,?,?,?,?,?,?)", data.Date, data.COCODE, data.Scripcode, data.Serial, data.ScripName, data.ClientName, data.Buysell, data.QTYSHARES, data.AVGPRICE)
		if err != nil {
			log.Println(err)
		}
	}

}

func get_mutualfunds_holdings(log_file *os.File) {
	url := "https://insbaapis.cmots.com/api/CompanyWiseMFHolding/92/10"

	log.SetOutput(log_file)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("Authorization", "Bearer "+api_key)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	res_body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(string(res_body))
}
