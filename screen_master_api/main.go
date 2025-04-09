package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type screen_master struct {
	Screen_name string
	Screen      string
	Screen_id   int64
}

var file, _ = os.Create("screen_master.log")
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
