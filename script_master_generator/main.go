package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-gota/gota/dataframe"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	file, err := os.Create("script_master.log")
	if err != nil {
		fmt.Println(err)
	}
	i := 2
	for i > 1 {
		username, password := get_userid_pass(file, 1)
		fmt.Println(username)
		fmt.Println(password)
		url := fmt.Sprintf("https://api.truedata.in/getAllSymbols?segment=eq&user=%s&password=%s&token=true&csv=true&csvHeader=true&ticksize=true&companyname=true&isin=true&limit=20000&circuit=true", username, password)
		fmt.Println(url)
		get_script_master_true_data(url, file)
		fmt.Println("going to sleep")
		time.Sleep(time.Hour * 6)
	}

}

func get_userid_pass(log_file *os.File, api_id int) (string, string) {
	log.SetOutput(log_file)
	var provider string
	var username string
	var password string
	var start_time string
	var last_update_time string
	db, err_db_open := sql.Open("mysql", "admin:saumitrasuparn@tcp(alerttradedb.czqug0e2in8p.ap-south-1.rds.amazonaws.com:3306)/alert_trade_db")

	if err_db_open != nil {
		log.Println(err_db_open)
	}

	err_db_stp := db.QueryRow("call alert_trade_db.stp_get_api_config(?)", api_id).Scan(&api_id, &provider, &username, &password, &start_time, &last_update_time)
	if err_db_stp != nil {
		log.Println(err_db_stp)
	}
	defer db.Close()
	return username, password
}

func get_script_master_true_data(url string, log_file *os.File) {
	log.SetOutput(log_file)
	res, err := http.Get(url)

	if err != nil {
		log.Println(err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}
	// log.Println(string(body))

	df := dataframe.ReadCSV(strings.NewReader(string(body)))

	fmt.Println(df)
	symbolid := df.Col("symbolid")
	symbol := df.Col("symbol")
	series := df.Col("series")
	isin := df.Col("isin")
	exchange := df.Col("exchange")
	// lotsize := df.Col("lotsize")
	// expiry := df.Col("expiry")
	// strike := df.Col("strike")
	metastocksymbol := df.Col("metastocksymbol")
	symbolalias := df.Col("symbolalias")
	token := df.Col("token")
	company := df.Col("company")
	ticksize := df.Col("ticksize")
	circuit := df.Col("circuit")
	db, err_db_open := sql.Open("mysql", "admin:saumitrasuparn@tcp(alerttradedb.czqug0e2in8p.ap-south-1.rds.amazonaws.com:3306)/alert_trade_db")

	if err_db_open != nil {
		log.Println(err_db_open)
	}
	for i := 0; i < symbolid.Len(); i++ {
		_, err_exec := db.Exec("call stp_insert_into_symbol_master(?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
			symbolid.Elem(i).String(),
			symbol.Elem(i).Float(),
			series.Elem(i).String(),
			isin.Elem(i).String(),
			exchange.Elem(i).String(),
			0,
			"",
			"",
			metastocksymbol.Elem(i).String(),
			symbolalias.Elem(i).String(),
			token.Elem(i).String(),
			company.Elem(i).String(),
			ticksize.Elem(i).String(),
			circuit.Elem(i).String(),
		)
		if err_exec != nil {
			log.Println("Error executing stored procedure:", err_exec)
		}
	}

}
