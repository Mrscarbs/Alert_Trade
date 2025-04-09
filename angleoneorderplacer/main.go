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
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type OrderStatusRequest struct {
	Variety string `json:"variety"`
	OrderID string `json:"orderid"`
}

type PlaceOrderReqBody struct {
	Variety         string `json:"variety"`
	Tradingsymbol   string `json:"tradingsymbol"`
	Symboltoken     string `json:"symboltoken"`
	Transactiontype string `json:"transactiontype"`
	Exchange        string `json:"exchange"`
	Ordertype       string `json:"ordertype"`
	Producttype     string `json:"producttype"`
	Duration        string `json:"duration"`
	Price           string `json:"price"`
	Squareoff       string `json:"squareoff"`
	Stoploss        string `json:"stoploss"`
	Quantity        string `json:"quantity"`
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
	router.POST("/place_order_angleone", placeOrderAngleOne)
	router.POST("/cancelOrderAngleOne", cancelOrderAngleOne)

	// Start the server
	router.Run("0.0.0.0:8080") // Runs on http://localhost:8080
}

func placeOrderAngleOne(c *gin.Context) {
	subID, _ := c.GetQuery("sub_id")
	userID, _ := c.GetQuery("id")

	// Fetch access token
	var accessToken string
	var err error
	err = db.QueryRow("call stp_getAccessToken_broker(?, ?)", subID, userID).Scan(&accessToken)
	if err != nil {
		log.Println("Error fetching access token:", err)
		c.JSON(500, gin.H{"error": "Failed to retrieve access token"})
		return
	}
	fmt.Println(accessToken)
	var reqBody PlaceOrderReqBody

	// Parse JSON request body
	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Convert struct to JSON
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.Println("JSON Encoding Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode JSON"})
		return
	}

	// Angel One API URL
	url := "https://apiconnect.angelone.in/rest/secure/angelbroking/order/v1/placeOrder"

	// Create HTTP request using http.NewRequest
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		log.Println("Request Creation Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Add required headers

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken) // Replace with actual token
	req.Header.Set("X-UserType", "USER")
	req.Header.Set("X-SourceID", "WEB")
	req.Header.Set("X-ClientLocalIP", "192.168.1.1")     // Replace with actual IP
	req.Header.Set("X-ClientPublicIP", "YOUR_PUBLIC_IP") // Replace with actual IP
	req.Header.Set("X-MACAddress", "YOUR_MAC_ADDRESS")   // Replace with actual MAC Address
	req.Header.Set("X-PrivateKey", "YOUR_API_KEY")       // Replace with Angel One API key

	// Execute HTTP request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Request Execution Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Request failed"})
		return
	}
	defer res.Body.Close()

	// Read response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Response Read Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	// Print response for debugging
	fmt.Println("Response from Angel One:", string(body))

	// Send the response from Angel One back to the user
	c.Data(res.StatusCode, "application/json", body)
}

func cancelOrderAngleOne(c *gin.Context) {

	subID, _ := c.GetQuery("sub_id")
	userID, _ := c.GetQuery("id")

	// Fetch access token
	var accessToken string
	var err error
	err = db.QueryRow("call stp_getAccessToken_broker(?, ?)", subID, userID).Scan(&accessToken)
	if err != nil {
		log.Println("Error fetching access token:", err)
		c.JSON(500, gin.H{"error": "Failed to retrieve access token"})
		return
	}
	fmt.Println(accessToken)
	var reqBody OrderStatusRequest

	// Parse JSON request body
	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Convert struct to JSON
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.Println("JSON Encoding Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode JSON"})
		return
	}

	// Angel One API URL
	url := "https://apiconnect.angelone.in/rest/secure/angelbroking/order/v1/cancelOrder"

	// Create HTTP request using http.NewRequest
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		log.Println("Request Creation Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Add required headers

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken) // Replace with actual token
	req.Header.Set("X-UserType", "USER")
	req.Header.Set("X-SourceID", "WEB")
	req.Header.Set("X-ClientLocalIP", "192.168.1.1")     // Replace with actual IP
	req.Header.Set("X-ClientPublicIP", "YOUR_PUBLIC_IP") // Replace with actual IP
	req.Header.Set("X-MACAddress", "YOUR_MAC_ADDRESS")   // Replace with actual MAC Address
	req.Header.Set("X-PrivateKey", "YOUR_API_KEY")       // Replace with Angel One API key

	// Execute HTTP request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Request Execution Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Request failed"})
		return
	}
	defer res.Body.Close()

	// Read response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Response Read Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	// Print response for debugging
	fmt.Println("Response from Angel One:", string(body))

	// Send the response from Angel One back to the user
	c.Data(res.StatusCode, "application/json", body)
}
