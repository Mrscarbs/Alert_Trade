package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Acess_token_response_frame struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	UserName    string `json:"userName"`
	Expires     string `json:".expires"`
	Issued      string `json:".issued"`
}

func main() {

	log_file, _ := os.Create("acess_token_logs.log")
	log.SetOutput(log_file)
	log.Println("generating token")
	username, password := Get_id_pass(1, log_file)

	i := 1
	for i < 2 {
		token_frame := generate_access_token(username, password, "password", log_file)

		expiry := update_insert_acess_token_db(token_frame.AccessToken, 1, token_frame.ExpiresIn, log_file)

		time.Sleep(time.Second * time.Duration(expiry))
		log.Println("new_token_generated")

	}

}

func Get_id_pass(api_id int, log_file *os.File) (string, string) {

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

func generate_access_token(username string, password string, grantType string, log_file *os.File) Acess_token_response_frame {
	var token_frame Acess_token_response_frame
	log.SetOutput(log_file)
	serviceURL := "https://auth.truedata.in/token"

	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)
	data.Set("grant_type", grantType)

	req, err_req_generator := http.NewRequest("POST", serviceURL, bytes.NewBufferString(data.Encode()))
	if err_req_generator != nil {
		log.Println(err_req_generator)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err_access_token_resp := http.DefaultClient.Do(req)
	if err_access_token_resp != nil {
		log.Println(err_access_token_resp)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	json.Unmarshal(body, &token_frame)

	return token_frame

}

func update_insert_acess_token_db(acess_token string, api_id int, Exipers_in int, log_file *os.File) int {

	log.SetOutput(log_file)
	current_time := time.Now()
	current_timestamp := current_time.Unix()

	expire_timestamp := current_timestamp + int64(Exipers_in)

	db, err_db_open := sql.Open("mysql", "admin:saumitrasuparn@tcp(alerttradedb.czqug0e2in8p.ap-south-1.rds.amazonaws.com:3306)/alert_trade_db")

	if err_db_open != nil {
		log.Println(err_db_open)
	}

	_, err_db_stp := db.Exec("call alert_trade_db.stp_update_access_token(?,?,?)", api_id, acess_token, expire_timestamp)
	if err_db_stp != nil {
		log.Println(err_db_stp)
	}
	db.Close()
	return Exipers_in
}
