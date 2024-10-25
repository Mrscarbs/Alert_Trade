package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var log_file, _ = os.Create("mf_inc_dec_log.log")

type change_val struct {
	Change float64
}

func main() {
	router := gin.Default()
	router.GET("/get_company_wise_mf_increment_decrement", get_company_wise_mf_increment_decrement)
	router.Run("0.0.0.0:8090")

}

func get_company_wise_mf_increment_decrement(c *gin.Context) {
	var change1 float64
	code, _ := c.GetQuery("co_code")
	num, _ := strconv.Atoi(code)
	db, err_db_open := sql.Open("mysql", "root:Karma100%@tcp(alerttrade.cbgqgqswkxrn.eu-north-1.rds.amazonaws.com:3306)/alert_trade_db")
	log.SetOutput(log_file)
	if err_db_open != nil {
		log.Println(err_db_open)
	}
	err := db.QueryRow("call stp_subtract_shares_mf(?)", num).Scan(&change1)
	if err != nil {
		log.Println(err)
	}
	change_obj := change_val{Change: change1}
	c.IndentedJSON(http.StatusOK, change_obj)
}
