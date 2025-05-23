package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

type token_details struct {
	sub_id       int64
	access_token string
	status       int
}
type AutoGenerated struct {
	Status string `json:"status"`
	Data   []struct {
		Isin                     string  `json:"isin"`
		CncUsedQuantity          int     `json:"cnc_used_quantity"`
		CollateralType           string  `json:"collateral_type"`
		CompanyName              string  `json:"company_name"`
		Haircut                  float64 `json:"haircut"`
		Product                  string  `json:"product"`
		Quantity                 int     `json:"quantity"`
		Tradingsymbol            string  `json:"tradingsymbol"`
		LastPrice                float64 `json:"last_price"`
		ClosePrice               float64 `json:"close_price"`
		Pnl                      float64 `json:"pnl"`
		DayChange                float64 `json:"day_change"`
		DayChangePercentage      float64 `json:"day_change_percentage"`
		InstrumentToken          string  `json:"instrument_token"`
		AveragePrice             float64 `json:"average_price"`
		CollateralQuantity       int     `json:"collateral_quantity"`
		CollateralUpdateQuantity int     `json:"collateral_update_quantity"`
		TradingSymbol            string  `json:"trading_symbol"`
		T1Quantity               int     `json:"t1_quantity"`
		Exchange                 string  `json:"exchange"`
	} `json:"data"`
}

type detailed_portfolio struct {
	stocks []string
	cocode []int64
}

var wg sync.WaitGroup

func main() {
	var subuser_map = make(map[int64]detailed_portfolio)
	log_file, err := os.Create("mail_system.log")
	if err != nil {
		fmt.Println(err)
	}
	db, err_db_open := sql.Open("mysql", "admin:saumitrasuparn@tcp(alerttradedb.czqug0e2in8p.ap-south-1.rds.amazonaws.com:3306)/alert_trade_db")
	log.SetOutput(log_file)
	if err_db_open != nil {
		log.Println(err_db_open)
	}
	rows, err := db.Query("call tbl_GetAccessTokenStatus()")
	if err != nil {
		log.Println(err)
	}

	var list_details []token_details
	for rows.Next() {
		var token_details_obj token_details

		rows.Scan(&token_details_obj.sub_id, &token_details_obj.access_token, &token_details_obj.status)
		list_details = append(list_details, token_details_obj)
	}
	fmt.Println(list_details)

	for i := 0; i < len(list_details); i++ {
		fmt.Println(list_details[i])
		obj := list_details[i]
		wg.Add(1)
		go get_portfolio_holdings(&wg, obj, log_file, &subuser_map)
		if obj.status == 1 {
			fmt.Println("sucess")
		} else {
			fmt.Println("failed")

		}
	}
	wg.Wait()
	fmt.Println(subuser_map)

}

func get_portfolio_holdings(wg *sync.WaitGroup, details token_details, log_file *os.File, map_details *map[int64]detailed_portfolio) {
	var portfolio_obj AutoGenerated
	accessToken := details.access_token
	log.SetOutput(log_file)
	url := "https://api.upstox.com/v2/portfolio/long-term-holdings"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	db, err_db_open := sql.Open("mysql", "admin:saumitrasuparn@tcp(alerttradedb.czqug0e2in8p.ap-south-1.rds.amazonaws.com:3306)/alert_trade_db")
	log.SetOutput(log_file)
	if err_db_open != nil {
		log.Println(err_db_open)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}

	json.Unmarshal(body, &portfolio_obj)

	var code int64
	var port detailed_portfolio
	for i := 0; i < len(portfolio_obj.Data); i++ {

		db.QueryRow("call stp_Getcocodebycompany(?)", portfolio_obj.Data[i].TradingSymbol).Scan(&code)
		port.stocks = append(port.stocks, portfolio_obj.Data[i].TradingSymbol)
		port.cocode = append(port.cocode, code)
		fmt.Println(code)
	}
	(*map_details)[details.sub_id] = port

	fmt.Println(portfolio_obj)
	fmt.Println(details.sub_id)

	wg.Done()
}
