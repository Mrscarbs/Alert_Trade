package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type stuct_quaterly_ratios struct {
	Success bool `json:"success"`
	Data    []struct {
		CoCode            int     `json:"co_code"`
		Qtrend            int     `json:"qtrend"`
		Mcap              float64 `json:"mcap"`
		Ev                float64 `json:"ev"`
		Pe                float64 `json:"pe"`
		Pbv               float64 `json:"pbv"`
		Eps               float64 `json:"eps"`
		Bookvalue         float64 `json:"bookvalue"`
		Ebit              float64 `json:"ebit"`
		Ebitda            float64 `json:"ebitda"`
		EvSales           float64 `json:"ev_sales"`
		EvEbitda          float64 `json:"ev_ebitda"`
		Netincomemargin   float64 `json:"netincomemargin"`
		Grossincomemargin float64 `json:"grossincomemargin"`
		Ebitdamargin      float64 `json:"ebitdamargin"`
		Epsdiluted        float64 `json:"epsdiluted"`
		Netsales          float64 `json:"netsales"`
		Netprofit         float64 `json:"netprofit"`
		Cogs              float64 `json:"cogs"`
	} `json:"data"`
	Message string `json:"message"`
}

var db *sql.DB
var api_key = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1bmlxdWVfbmFtZSI6Imluc2JhYXBpcyIsInJvbGUiOiJBZG1pbiIsIm5iZiI6MTczNDk2MDMwMiwiZXhwIjoxNzY2NjY5MTAyLCJpYXQiOjE3MzQ5NjAzMDIsImlzcyI6Imh0dHA6Ly9sb2NhbGhvc3Q6NTAxOTEiLCJhdWQiOiJodHRwOi8vbG9jYWxob3N0OjUwMTkxIn0.UqzyAKBMcDMmPL-kgaZtnusOAWOuB3v1tVIu_PZsJp8"

func main() {
	var err error
	file, err := os.Create("golang_findamental_cmots_quaterly_log.log")
	if err != nil {
		log.Println(err)
	}
	log.SetOutput(file)
	db, err = sql.Open("mysql", "admin:saumitrasuparn@tcp(alerttradedb.czqug0e2in8p.ap-south-1.rds.amazonaws.com:3306)/alert_trade_db")
	if err != nil {
		log.Println(err)
	}

	for {
		get_all_quaterly_ratios()
		time.Sleep(time.Hour * 24)
	}

}

func get_all_quaterly_ratios() {
	var unfolded_quaterly stuct_quaterly_ratios
	rows, err := db.Query("call stp_get_all_cocodes_cmots()")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	list_cocodes := []int{}
	for rows.Next() {
		var cocode int
		rows.Scan(&cocode)
		list_cocodes = append(list_cocodes, cocode)
	}
	for i := 0; i < len(list_cocodes); i++ {
		str_cocode := strconv.Itoa(list_cocodes[i])
		url := fmt.Sprintf("https://insbaapis.cmots.com/api/QuarterlyRatio/%s/S/", str_cocode)
		fmt.Println(url)
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
		json.Unmarshal(res_body, &unfolded_quaterly)
		if len(unfolded_quaterly.Data) == 0 {
			error_message := fmt.Sprintf("TTM data not found for cocode: %s", str_cocode)
			log.Println(error_message)
			continue
		}
		fmt.Println(string(res_body))

		for i := 0; i < len(unfolded_quaterly.Data); i++ {
			data := unfolded_quaterly.Data[i]
			_, err := db.Exec("call stp_upsert_quarterly_financials_cmots_final(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", data.CoCode, data.Qtrend, data.Mcap, data.Ev, data.Pe, data.Pbv, data.Eps, data.Bookvalue, data.Ebit, data.Ebitda, data.EvSales, data.EvEbitda, data.Netincomemargin, data.Grossincomemargin, data.Ebitdamargin, data.Epsdiluted, data.Netsales, data.Netprofit, data.Cogs)
			if err != nil {
				log.Println(err)
			}
		}

	}
	defer db.Close()
}
