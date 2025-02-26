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
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type SubUser struct {
	NSubUserID   int            `json:"n_sub_user_id"`
	NUserID      int            `json:"n_user_id"`
	BrokerName   string         `json:"broker_name_string"`
	NStateID     int            `json:"n_state_id"`
	BrokerUserID sql.NullString `json:"broker_user_id"` // Use sql.NullString
	AccessToken  sql.NullString `json:"access_token"`
	RefreshToken sql.NullString `json:"refresh_token"`
	NExpiresAt   int64          `json:"n_expires_at"`
	NLastLogin   int64          `json:"n_last_login"`
	NStatus      int            `json:"n_status"`
	LoginLink    sql.NullString `json:"login_link_string"`
	APIKey       string         `json:"api_key_string"`
	SecretKey    string         `json:"secret_key_string"`
	NCreatedAt   int64          `json:"n_created_at"`
	NUpdatedAt   int64          `json:"n_updated_at"`
}
type upstox_response struct {
	Status string `json:"status"`
	Data   []struct {
		Isin                     string  `json:"isin"`
		CncUsedQuantity          int     `json:"cnc_used_quantity"`
		CollateralType           string  `json:"collateral_type"`
		CompanyName              string  `json:"company_name"`
		Haircut                  float64 `json:"haircut"`
		Product                  string  `json:"product"`
		Quantity                 int     `json:"quantity"`
		TradingSymbol            string  `json:"trading_symbol"`
		Tradingsymbol            string  `json:"tradingsymbol"`
		LastPrice                float64 `json:"last_price"`
		ClosePrice               float64 `json:"close_price"`
		Pnl                      float64 `json:"pnl"`
		DayChange                int     `json:"day_change"`
		DayChangePercentage      int     `json:"day_change_percentage"`
		InstrumentToken          string  `json:"instrument_token"`
		AveragePrice             float64 `json:"average_price"`
		CollateralQuantity       int     `json:"collateral_quantity"`
		CollateralUpdateQuantity int     `json:"collateral_update_quantity"`
		T1Quantity               int     `json:"t1_quantity"`
		Exchange                 string  `json:"exchange"`
	} `json:"data"`
	Price_action []candles_upstox
	Key_metric   []FinancialData
}

type candles_upstox struct {
	Status string `json:"status"`
	Data   struct {
		Candles [][]interface{} `json:"candles"`
	} `json:"data"`
}

type FinancialData struct {
	CO_CODE                     int     `json:"co_code"`
	TTMAson                     int     `json:"ttm_ason"`
	MCAP                        float64 `json:"mcap"`
	EV                          float64 `json:"ev"`
	PE                          float64 `json:"pe"`
	PBV                         float64 `json:"pbv"`
	DIVYIELD                    float64 `json:"div_yield"`
	EPS                         float64 `json:"eps"`
	BookValue                   float64 `json:"book_value"`
	ROA_TTM                     float64 `json:"roa_ttm"`
	ROE_TTM                     float64 `json:"roe_ttm"`
	ROCE_TTM                    float64 `json:"roce_ttm"`
	EBIT_TTM                    float64 `json:"ebit_ttm"`
	EBITDA_TTM                  float64 `json:"ebitda_ttm"`
	EV_Sales_TTM                float64 `json:"ev_sales_ttm"`
	EV_EBITDA_TTM               float64 `json:"ev_ebitda_ttm"`
	NetIncomeMargin_TTM         float64 `json:"net_income_margin_ttm"`
	GrossIncomeMargin_TTM       float64 `json:"gross_income_margin_ttm"`
	AssetTurnover_TTM           float64 `json:"asset_turnover_ttm"`
	CurrentRatio_TTM            float64 `json:"current_ratio_ttm"`
	Debt_Equity_TTM             float64 `json:"debt_equity_ttm"`
	Sales_TotalAssets_TTM       float64 `json:"sales_total_assets_ttm"`
	NetDebt_EBITDA_TTM          float64 `json:"net_debt_ebitda_ttm"`
	EBITDA_Margin_TTM           float64 `json:"ebitda_margin_ttm"`
	TotalShareHoldersEquity_TTM float64 `json:"total_shareholders_equity_ttm"`
	ShorttermDebt_TTM           float64 `json:"short_term_debt_ttm"`
	LongtermDebt_TTM            float64 `json:"long_term_debt_ttm"`
	SharesOutstanding           float64 `json:"shares_outstanding"`
	EPSDiluted                  float64 `json:"eps_diluted"`
	NetSales                    float64 `json:"net_sales"`
	Netprofit                   float64 `json:"net_profit"`
	AnnualDividend              float64 `json:"annual_dividend"`
	COGS                        float64 `json:"cogs"`
	PEGRatio_TTM                float64 `json:"peg_ratio_ttm"`
	DividendPayout_TTM          float64 `json:"dividend_payout_ttm"`
}

type prompt struct {
	Question            string
	Broker_details      []any
	Structural_overview string
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}
type gpt_answer struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string      `json:"role"`
			Content string      `json:"content"`
			Refusal interface{} `json:"refusal"`
		} `json:"message"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens        int `json:"prompt_tokens"`
		CompletionTokens    int `json:"completion_tokens"`
		TotalTokens         int `json:"total_tokens"`
		PromptTokensDetails struct {
			CachedTokens int `json:"cached_tokens"`
			AudioTokens  int `json:"audio_tokens"`
		} `json:"prompt_tokens_details"`
		CompletionTokensDetails struct {
			ReasoningTokens          int `json:"reasoning_tokens"`
			AudioTokens              int `json:"audio_tokens"`
			AcceptedPredictionTokens int `json:"accepted_prediction_tokens"`
			RejectedPredictionTokens int `json:"rejected_prediction_tokens"`
		} `json:"completion_tokens_details"`
	} `json:"usage"`
	ServiceTier       string `json:"service_tier"`
	SystemFingerprint string `json:"system_fingerprint"`
}

var db *sql.DB
var log_file, _ = os.Create("llm_Api.log")

func main() {
	var err error
	db, err = sql.Open("mysql", "admin:saumitrasuparn@tcp(alerttradedb.czqug0e2in8p.ap-south-1.rds.amazonaws.com:3306)/alert_trade_db")
	if err != nil {
		log.Println(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	router := gin.Default()
	router.GET("llm_invoke", llm_invoke)
	router.Run("localhost:8080")
}

func llm_invoke(c *gin.Context) {
	user_id, _ := c.GetQuery("user_id")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	log.SetOutput(log_file)

	var subuser_list []SubUser
	rows, err := db.Query("call stp_GetSubUsersByUserId(?)", user_id)
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		var subUser SubUser
		err := rows.Scan(
			&subUser.NSubUserID, &subUser.NUserID, &subUser.BrokerName,
			&subUser.NStateID, &subUser.BrokerUserID, &subUser.AccessToken,
			&subUser.RefreshToken, &subUser.NExpiresAt, &subUser.NLastLogin,
			&subUser.NStatus, &subUser.LoginLink, &subUser.APIKey,
			&subUser.SecretKey, &subUser.NCreatedAt, &subUser.NUpdatedAt,
		)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to scan row", "details": err.Error()})
			return
		}
		subuser_list = append(subuser_list, subUser)
	}
	var mut sync.Mutex
	var wg sync.WaitGroup
	list_broker_response := []any{}
	for _, val := range subuser_list {
		if val.BrokerName == "upstox" {
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				acess_token := val.AccessToken.String
				upstox_resp := upstox_details(acess_token)
				mut.Lock()
				list_broker_response = append(list_broker_response, upstox_resp)
				mut.Unlock()
				wg.Done()
			}(&wg)

		}
	}

	wg.Wait()
	// c.IndentedJSON(http.StatusOK, list_broker_response)
	structural_explainatio := "The `broker_details` struct represents the complete response from Upstox's API, containing financial and market data. The primary components include `data`, `price_action`, and `key_metric`. The `data` field is a slice containing details about various stocks, such as their trading symbols, quantities held, prices, profit and loss (PnL), and collateral details. Each stock entry includes identifiers like `isin` and `instrument_token`, along with financial information like `average_price`, `last_price`, and `day_change_percentage`. The `price_action` field contains a slice of `candles_upstox`, which holds OHLCV (Open, High, Low, Close, Volume) data in a nested array format, where each entry represents a candlestick for a specific time interval. Each candlestick entry follows the structure `[timestamp, open_price, high_price, low_price, close_price, volume, 0]`, ensuring historical price tracking for individual stocks in a chronological sequence. The `key_metric` field contains a slice of `FinancialData`, providing in-depth financial metrics such as market capitalization (`MCAP`), earnings per share (`EPS`), book value, return ratios (`ROA_TTM`, `ROE_TTM`), debt-equity ratio, and valuation multiples (`PE`, `PBV`, `EV/EBITDA`). These metrics give a comprehensive view of a stock’s fundamental financial health, helping traders and analysts assess its investment potential. The struct is designed to capture stock-specific details, historical price movements, and key financial indicators in a structured manner.pls use broker details for reference regarding price action and key metrics or other things"

	for {
		message_type, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		detail_question := prompt{Question: string(message), Broker_details: list_broker_response, Structural_overview: structural_explainatio}
		json_question, err := json.Marshal(detail_question)
		fmt.Println(string(json_question))
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		answer := call_gpt_api(string(json_question))

		err = conn.WriteMessage(message_type, []byte(answer))

		if err != nil {
			log.Println(err)
			return
		}

	}

}

func upstox_details(access_token string) upstox_response {
	url := "https://api.upstox.com/v2/portfolio/long-term-holdings"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
	}
	req.Header.Add("Authorization", "Bearer "+access_token)
	req.Header.Add("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}
	var response_struct upstox_response

	fmt.Println(string(body))

	json.Unmarshal(body, &response_struct)
	fmt.Println(response_struct.Price_action)

	for _, val := range response_struct.Data {
		re := upstox_intraday_candel(val.InstrumentToken)
		isin := removePrefix(val.InstrumentToken)
		re_key_metrics := get_key_metrics_isin(isin)

		response_struct.Price_action = append(response_struct.Price_action, re)
		response_struct.Key_metric = append(response_struct.Key_metric, re_key_metrics)

	}
	return response_struct

}

func upstox_intraday_candel(instrument_key string) candles_upstox {
	toDate := time.Now().Format("2006-01-02")

	fromDate := time.Now().AddDate(0, 0, -10).Format("2006-01-02")
	url := fmt.Sprintf("https://api.upstox.com/v2/historical-candle/%s/30minute/%s/%s", instrument_key, toDate, fromDate)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
	}
	req.Header.Add("Accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	var candles candles_upstox
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}
	json.Unmarshal(body, &candles)

	return candles
}

func get_key_metrics_isin(isin string) FinancialData {

	var key_metrics FinancialData
	var code int
	db.QueryRow("call stp_GetCoCodeByISIN_cmots(?)", isin).Scan(&code)
	err := db.QueryRow("call stp_get_financial_data_by_cocode(?)", code).Scan(
		&key_metrics.CO_CODE, &key_metrics.TTMAson, &key_metrics.MCAP, &key_metrics.EV, &key_metrics.PE, &key_metrics.PBV,
		&key_metrics.DIVYIELD, &key_metrics.EPS, &key_metrics.BookValue, &key_metrics.ROA_TTM, &key_metrics.ROE_TTM,
		&key_metrics.ROCE_TTM, &key_metrics.EBIT_TTM, &key_metrics.EBITDA_TTM, &key_metrics.EV_Sales_TTM,
		&key_metrics.EV_EBITDA_TTM, &key_metrics.NetIncomeMargin_TTM, &key_metrics.GrossIncomeMargin_TTM,
		&key_metrics.AssetTurnover_TTM, &key_metrics.CurrentRatio_TTM, &key_metrics.Debt_Equity_TTM,
		&key_metrics.Sales_TotalAssets_TTM, &key_metrics.NetDebt_EBITDA_TTM, &key_metrics.EBITDA_Margin_TTM,
		&key_metrics.TotalShareHoldersEquity_TTM, &key_metrics.ShorttermDebt_TTM, &key_metrics.LongtermDebt_TTM,
		&key_metrics.SharesOutstanding, &key_metrics.EPSDiluted, &key_metrics.NetSales, &key_metrics.Netprofit,
		&key_metrics.AnnualDividend, &key_metrics.COGS, &key_metrics.PEGRatio_TTM, &key_metrics.DividendPayout_TTM)

	if err != nil {
		log.Println("Error fetching financial data:", err)

	}

	return key_metrics

}

func removePrefix(input string) string {
	return strings.TrimPrefix(input, "NSE_EQ|")
}

func call_gpt_api(prompt string) string {
	// Get OpenAI API Key from Environment Variables
	apiKey := os.Getenv("OPENAI_API_KEY")

	// OpenAI API URL
	url := "https://api.openai.com/v1/chat/completions"

	// Create the Request Payload
	requestBody := OpenAIRequest{
		Model: "gpt-4o",
		Messages: []Message{
			{Role: "developer", Content: "You are a helpful trading assistant.pls dont give example calculaitions do the calculations your self based on the broker details from the user"},
			{Role: "user", Content: prompt},
		},
	}

	// Convert Request to JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)

	}

	// Create HTTP Request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)

	}

	// Set Headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// ✅ Use `http.DefaultClient.Do()` to Send the Request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)

	}
	defer resp.Body.Close() // Close response body

	// Read Response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)

	}
	var answer gpt_answer
	// Print Raw GPT-4o Response
	// fmt.Println("GPT-4o Response:", string(body))
	err = json.Unmarshal(body, &answer)
	if err != nil {
		log.Println(err)
	}
	return answer.Choices[0].Message.Content
}
