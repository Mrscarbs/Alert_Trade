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

var log_file, _ = os.Create("additional_application_apis_log.log")

type res_company_name_cocode struct {
	Symbol  string
	Co_code int64
}

type BulkDeal struct {
	Datetime   string  `json:"datetime"` // You can use time.Time instead of string if you need datetime parsing
	COCode     int     `json:"co_code"`
	ScripCode  string  `json:"scripcode"`
	Serial     int     `json:"serial"`
	ScripName  string  `json:"scripname"`
	ClientName string  `json:"clientname"`
	BuySell    string  `json:"buysell"`
	QtyShares  int     `json:"qtyshares"`
	AvgPrice   float64 `json:"avg_price"`
}

type BlockDeal struct {
	Datetime   string  `json:"datetime"`
	COCode     int     `json:"co_code"`
	ScripCode  string  `json:"scripcode"`
	Serial     int     `json:"serial"`
	ScripName  string  `json:"scripname"`
	ClientName string  `json:"clientname"`
	BuySell    string  `json:"buysell"`
	QtyShares  int     `json:"qtyshares"`
	AvgPrice   float64 `json:"avg_price"`
}

type Shareholding struct {
	ID        int    `json:"id"`
	Candidate string `json:"candidate"`
	ShareName string `json:"share_name"`
	Amount    int    `json:"amount"`
}

func main() {

	router := gin.Default()
	router.GET("delete_position", delete_position)
	router.GET("get_cocde_symbol", get_cocde_symbol)
	router.GET("get_bulk_deals", get_bulk_deals)
	router.GET("get_block_deals", get_block_deals)
	router.GET("get_bulk_deals_cocode", get_bulk_deals_cocode)
	router.GET("get_block_deals_cocode", get_block_deals_cocode)
	router.GET("get_shareholding_mp_mla", get_shareholding_mp_mla)
	router.Run("0.0.0.0:8091")

}

func delete_position(c *gin.Context) {
	log.SetOutput(log_file)
	trade_id, _ := c.GetQuery("trade_id")
	trade_id_int, err := strconv.Atoi(trade_id)
	if err != nil {
		log.Println(err)
	}
	db, err := sql.Open("mysql", "admin:saumitrasuparn@tcp(alerttradedb.czqug0e2in8p.ap-south-1.rds.amazonaws.com:3306)/alert_trade_db")

	if err != nil {
		log.Println(err)
	}

	db.Exec("call stp_delete_user_position_by_trade_id(?)", trade_id_int)
}

func get_cocde_symbol(c *gin.Context) {
	var company_detailt res_company_name_cocode
	log.SetOutput(log_file)
	comp, _ := c.GetQuery("company_name")

	db, err := sql.Open("mysql", "admin:saumitrasuparn@tcp(alerttradedb.czqug0e2in8p.ap-south-1.rds.amazonaws.com:3306)/alert_trade_db")

	if err != nil {
		log.Println(err)
	}

	db.QueryRow("call stp_get_symbol_and_code_by_companyname(?)", comp).Scan(&company_detailt.Symbol, &company_detailt.Co_code)

	c.IndentedJSON(http.StatusOK, company_detailt)
}

func get_bulk_deals(c *gin.Context) {

	log.SetOutput(log_file)

	db, err := sql.Open("mysql", "admin:saumitrasuparn@tcp(alerttradedb.czqug0e2in8p.ap-south-1.rds.amazonaws.com:3306)/alert_trade_db")

	if err != nil {
		log.Println(err)
	}

	bulk_list := []BulkDeal{}

	rows, err := db.Query("call stp_get_all_bulk_deals()")
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		var bulk_object BulkDeal
		err := rows.Scan(&bulk_object.Datetime, &bulk_object.COCode, &bulk_object.ScripCode, &bulk_object.Serial, &bulk_object.ScripName, &bulk_object.ClientName, &bulk_object.BuySell, &bulk_object.QtyShares, &bulk_object.AvgPrice)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Row scan error"})
			return
		}
		bulk_list = append(bulk_list, bulk_object)
	}
	c.IndentedJSON(http.StatusOK, bulk_list)

}
func get_block_deals(c *gin.Context) {
	log.SetOutput(log_file)

	db, err := sql.Open("mysql", "admin:saumitrasuparn@tcp(alerttradedb.czqug0e2in8p.ap-south-1.rds.amazonaws.com:3306)/alert_trade_db")
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}
	defer db.Close()

	block_list := []BlockDeal{}

	rows, err := db.Query("call stp_get_all_block_deals()")
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Query execution error"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var block_object BlockDeal
		err := rows.Scan(&block_object.Datetime, &block_object.COCode, &block_object.ScripCode, &block_object.Serial, &block_object.ScripName, &block_object.ClientName, &block_object.BuySell, &block_object.QtyShares, &block_object.AvgPrice)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Row scan error"})
			return
		}
		block_list = append(block_list, block_object)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Rows iteration error"})
		return
	}

	c.IndentedJSON(http.StatusOK, block_list)
}
func get_bulk_deals_cocode(c *gin.Context) {

	log.SetOutput(log_file)
	cocode, _ := c.GetQuery("co_code")
	cocode_int, err := strconv.Atoi(cocode)
	if err != nil {
		log.Println(err)
	}
	db, err := sql.Open("mysql", "admin:saumitrasuparn@tcp(alerttradedb.czqug0e2in8p.ap-south-1.rds.amazonaws.com:3306)/alert_trade_db")

	if err != nil {
		log.Println(err)
	}

	bulk_list := []BulkDeal{}

	rows, err := db.Query("call stp_get_bulk_deals_by_cocode(?)", cocode_int)
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		var bulk_object BulkDeal
		err := rows.Scan(&bulk_object.Datetime, &bulk_object.COCode, &bulk_object.ScripCode, &bulk_object.Serial, &bulk_object.ScripName, &bulk_object.ClientName, &bulk_object.BuySell, &bulk_object.QtyShares, &bulk_object.AvgPrice)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Row scan error"})
			return
		}
		bulk_list = append(bulk_list, bulk_object)
	}
	c.IndentedJSON(http.StatusOK, bulk_list)

}
func get_block_deals_cocode(c *gin.Context) {

	log.SetOutput(log_file)
	cocode, _ := c.GetQuery("co_code")
	cocode_int, err := strconv.Atoi(cocode)
	if err != nil {
		log.Println(err)
	}
	db, err := sql.Open("mysql", "admin:saumitrasuparn@tcp(alerttradedb.czqug0e2in8p.ap-south-1.rds.amazonaws.com:3306)/alert_trade_db")

	if err != nil {
		log.Println(err)
	}

	bulk_list := []BlockDeal{}

	rows, err := db.Query("call stp_get_block_deals_by_cocode(?)", cocode_int)
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		var bulk_object BlockDeal
		err := rows.Scan(&bulk_object.Datetime, &bulk_object.COCode, &bulk_object.ScripCode, &bulk_object.Serial, &bulk_object.ScripName, &bulk_object.ClientName, &bulk_object.BuySell, &bulk_object.QtyShares, &bulk_object.AvgPrice)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Row scan error"})
			return
		}
		bulk_list = append(bulk_list, bulk_object)
	}
	c.IndentedJSON(http.StatusOK, bulk_list)

}

func get_shareholding_mp_mla(c *gin.Context) {

	comp_name, _ := c.GetQuery("comp_name")

	db, err := sql.Open("mysql", "admin:saumitrasuparn@tcp(alerttradedb.czqug0e2in8p.ap-south-1.rds.amazonaws.com:3306)/alert_trade_db")

	if err != nil {
		log.Println(err)
	}
	rows, err := db.Query("call stp_get_shareholdings_mp_mla_by_share_name(?)", comp_name)
	if err != nil {
		log.Println(err)
	}

	list_shareholding := []Shareholding{}
	for rows.Next() {
		var share_holding_mp Shareholding
		rows.Scan(&share_holding_mp.ID, &share_holding_mp.Candidate, &share_holding_mp.ShareName, &share_holding_mp.Amount)
		list_shareholding = append(list_shareholding, share_holding_mp)
	}
	c.IndentedJSON(http.StatusOK, list_shareholding)
}
