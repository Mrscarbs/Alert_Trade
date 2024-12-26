package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type screen_master struct {
	Screen_name string
	Screen      string
	Screen_id   int64
}

var db, _ = sql.Open("mysql", "admin:saumitrasuparn@tcp(alerttradedb.czqug0e2in8p.ap-south-1.rds.amazonaws.com:3306)/alert_trade_db")
var file, _ = os.Create("screen_master.log")

func main() {
	router := gin.Default()
	router.GET("/screens_details", screens_details)
	router.Run("0.0.0.0:8080")
}

func screens_details(c *gin.Context) {
	log.SetOutput(file)
	var screen_obj []screen_master
	rows, err := db.Query("call stp_GetAllscreenMaster()")

	if err != nil {
		log.Println(err)
	}

	for rows.Next() {
		var Screen_name string
		var screen string
		var screen_id int64
		err := rows.Scan(&Screen_name, &screen, &screen_id)
		if err != nil {
			log.Println(err)
		}
		screen_ab := screen_master{Screen_name: Screen_name, Screen: screen, Screen_id: screen_id}
		screen_obj = append(screen_obj, screen_ab)
	}
	defer rows.Close()
	c.IndentedJSON(http.StatusOK, screen_obj)
}
