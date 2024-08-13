package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-gota/gota/dataframe"
	_ "github.com/go-sql-driver/mysql"
)

type actions struct {
	Symbol      string
	Company     string
	Series      string
	Purpose     string
	Face_value  string
	Exdate      string
	Record_date string
	Bc_start    string
	Bc_end      string
	Ratio       float64
}

var main_log, _ = os.Create("main_corp_log.log")

func main() {
	router := gin.Default()
	router.GET("/get_corp_action", get_corp_action)
	router.Run("localhost:8081")

}

func get_acess_token(log_file *os.File) (string, int64) {
	log.SetOutput(log_file)
	db, err := sql.Open("mysql", "root:Karma100%@/alert_trade_db")
	if err != nil {
		log.Println(err)
	}
	var acess_token string
	var expiry int64
	var api_ID int
	defer db.Close()
	db.QueryRow("call stp_get_access_token_api_id(?)", 1).Scan(&acess_token, &expiry, &api_ID)

	return acess_token, expiry
}

func get_corp_action(c *gin.Context) {
	ticker, _ := c.GetQuery("symbol")
	log.SetOutput(main_log)
	acess_token, _ := get_acess_token(main_log)
	url := fmt.Sprintf("https://history.truedata.in/getcorpaction?symbol=%s&response=csv", ticker)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("authorization", "Bearer "+acess_token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}
	var action_corp []actions
	log.Println(string(data))
	df := dataframe.ReadCSV(strings.NewReader(string(data)))
	fmt.Println(df)
	symbol := df.Col("symbol")
	company := df.Col("company")
	series := df.Col("series")
	purpose := df.Col("purpose")
	face_value := df.Col("face_value")
	ex_date := df.Col("ex_date")
	record_date := df.Col("record_date")
	bc_start := df.Col("bc_start")
	bc_end := df.Col("bc_end")
	ratio := df.Col("ratio")
	for i := 0; i < company.Len(); i++ {
		elem_corp_actions := actions{
			Symbol:      symbol.Elem(i).String(),
			Company:     company.Elem(i).String(),
			Series:      series.Elem(i).String(),
			Purpose:     purpose.Elem(i).String(),
			Face_value:  face_value.Elem(i).String(),
			Exdate:      ex_date.Elem(i).String(),
			Record_date: record_date.Elem(i).String(),
			Bc_start:    bc_start.Elem(i).String(),
			Bc_end:      bc_end.Elem(i).String(),
			Ratio:       ratio.Elem(i).Float(),
		}
		action_corp = append(action_corp, elem_corp_actions)
	}
	fmt.Println(action_corp)
	c.IndentedJSON(http.StatusOK, action_corp)
}
