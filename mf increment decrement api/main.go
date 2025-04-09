package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var log_file, _ = os.Create("mf_inc_dec_log.log")

type change_val struct {
	Change float64
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
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow all origins (change this to restrict access)
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, // Preflight request cache duration
	}))
	router.GET("/get_company_wise_mf_increment_decrement", get_company_wise_mf_increment_decrement)
	router.Run("0.0.0.0:8090")

}

func get_company_wise_mf_increment_decrement(c *gin.Context) {
	var change1 float64
	code, _ := c.GetQuery("co_code")
	num, _ := strconv.Atoi(code)
	log.SetOutput(log_file)

	err := db.QueryRow("call stp_subtract_shares_mf(?)", num).Scan(&change1)
	if err != nil {
		log.Println(err)
	}
	change_obj := change_val{Change: change1}
	c.IndentedJSON(http.StatusOK, change_obj)
}
