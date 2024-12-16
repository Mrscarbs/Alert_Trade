package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var log_file, _ = os.Create("buysell_log.log")

type order_response struct {
	Status string `json:"status"`
	Data   struct {
		OrderID string `json:"order_id"`
	} `json:"data"`
}

type OrderRequest struct {
	Quantity          int     `json:"quantity"`
	Product           string  `json:"product"`
	Validity          string  `json:"validity"`
	Price             float64 `json:"price"`
	Tag               string  `json:"tag"`
	InstrumentToken   string  `json:"instrument_token"`
	OrderType         string  `json:"order_type"`
	TransactionType   string  `json:"transaction_type"`
	DisclosedQuantity int     `json:"disclosed_quantity"`
	TriggerPrice      float64 `json:"trigger_price"`
	IsAmo             bool    `json:"is_amo"`
}

func main() {
	router := gin.Default()
	router.GET("/place_order", place_order)
	router.GET("/cancel_order", cancel_order)
	router.Run("0.0.0.0:8080")
}

func place_order(c *gin.Context) {
	subID, _ := c.GetQuery("sub_id")
	userID, _ := c.GetQuery("id")
	log.SetOutput(log_file)

	// Connect to the database
	dsn := "admin:saumitrasuparn@tcp(alerttradedb.czqug0e2in8p.ap-south-1.rds.amazonaws.com:3306)/alert_trade_db"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Println("Error connecting to database:", err)
		c.JSON(500, gin.H{"error": "Database connection failed"})
		return
	}
	defer db.Close()

	// Fetch access token
	var accessToken string
	err = db.QueryRow("call stp_getAccessToken_broker(?, ?)", subID, userID).Scan(&accessToken)
	if err != nil {
		log.Println("Error fetching access token:", err)
		c.JSON(500, gin.H{"error": "Failed to retrieve access token"})
		return
	}
	fmt.Println(accessToken)

	quantity, _ := c.GetQuery("quantity")
	product, _ := c.GetQuery("Product")
	validity, _ := c.GetQuery("Validity")
	price, _ := c.GetQuery("Price")
	tag, _ := c.GetQuery("Tag")
	instrumentToken, _ := c.GetQuery("InstrumentToken")
	orderType, _ := c.GetQuery("OrderType")
	transactionType, _ := c.GetQuery("TransactionType")
	disclosedQuantity, _ := c.GetQuery("DisclosedQuantity")
	triggerPrice, _ := c.GetQuery("TriggerPrice")
	isAmo, _ := c.GetQuery("IsAmo")

	// Convert query values to appropriate types
	quantityInt, _ := strconv.Atoi(quantity)
	priceFloat, _ := strconv.ParseFloat(price, 64)
	disclosedQuantityInt, _ := strconv.Atoi(disclosedQuantity)
	triggerPriceFloat, _ := strconv.ParseFloat(triggerPrice, 64)
	isAmoBool, _ := strconv.ParseBool(isAmo)

	// Prepare the order request body
	order := OrderRequest{
		Quantity:          quantityInt,
		Product:           product,
		Validity:          validity,
		Price:             priceFloat,
		Tag:               tag,
		InstrumentToken:   instrumentToken,
		OrderType:         orderType,
		TransactionType:   transactionType,
		DisclosedQuantity: disclosedQuantityInt,
		TriggerPrice:      triggerPriceFloat,
		IsAmo:             isAmoBool,
	}

	requestBody, err := json.Marshal(order)
	fmt.Println(string(requestBody))
	if err != nil {
		log.Println("Error marshalling request body:", err)
		c.JSON(500, gin.H{"error": "Failed to create order request"})
		return
	}

	// Make the API request
	url := "https://api-hft.upstox.com/v2/order/place"
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println("Error creating HTTP request:", err)
		c.JSON(500, gin.H{"error": "Failed to create HTTP request"})
		return
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Println("Error making HTTP request:", err)
		c.JSON(500, gin.H{"error": "Failed to execute API call"})
		return
	}
	defer res.Body.Close()

	// Read the response
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Error reading API response:", err)
		c.JSON(500, gin.H{"error": "Failed to read API response"})
		return
	}
	response_obj := order_response{}
	json.Unmarshal(responseBody, &response_obj)
	c.IndentedJSON(http.StatusOK, response_obj)

	fmt.Println(string(responseBody))
}

func cancel_order(c *gin.Context) {
	order_id, _ := c.GetQuery("order_id")

	// orderid_int, _ := strconv.Atoi(order_id)
	subID, _ := c.GetQuery("sub_id")
	userID, _ := c.GetQuery("id")
	log.SetOutput(log_file)

	// Connect to the database
	dsn := "admin:saumitrasuparn@tcp(alerttradedb.czqug0e2in8p.ap-south-1.rds.amazonaws.com:3306)/alert_trade_db"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Println("Error connecting to database:", err)
		c.JSON(500, gin.H{"error": "Database connection failed"})
		return
	}
	defer db.Close()

	// Fetch access token
	var accessToken string
	err = db.QueryRow("call stp_getAccessToken_broker(?, ?)", subID, userID).Scan(&accessToken)
	if err != nil {
		log.Println("Error fetching access token:", err)
		c.JSON(500, gin.H{"error": "Failed to retrieve access token"})
		return
	}
	fmt.Println(accessToken)
	url := fmt.Sprintf("https://api-hft.upstox.com/v2/order/cancel?order_id=%s", order_id)

	req, err := http.NewRequest("DELETE", url, nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		log.Println(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(string(body))
	var reponse order_response

	json.Unmarshal(body, &reponse)

	c.IndentedJSON(http.StatusOK, reponse)

}
