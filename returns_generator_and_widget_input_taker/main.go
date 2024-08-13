package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-gota/gota/dataframe"
	_ "github.com/go-sql-driver/mysql"
)

var mutex sync.Mutex
var log_file, _ = os.Create("return_user_position_log.log")

type user_position_details struct {
	EntityID   int     `json:"entity_id"`
	Ticker     string  `json:"ticker"`
	Quantity   float64 `json:"quantity"`
	TradePrice float64 `json:"trade_price"`
	TradeType  int     `json:"trade_type"`
	TradeID    int     `json:"trade_id"`
	Timestamp  int     `json:"timestamp"`
}

type trade_object struct {
	Closes            float64
	High              float64
	Low               float64
	Timestamp         int
	Position_price    float64
	Position_quantity int
}

type symbols_entity_details struct {
	List_symbols []string
	Entity_id    int
}
type trade_respnse_final struct {
	Equity        []float64
	Trade_details []trade_object
}

func main() {

	router := gin.Default()
	router.POST("/insert_user_position_detail", insert_user_position_detail)
	router.POST("/get_historical_performance", get_equity_and_historical_performance)
	router.Run("localhost:8080")

}

func insert_user_position_detail(c *gin.Context) {

	log.SetOutput(log_file)
	var user_positions user_position_details

	err_binding := c.BindJSON(&user_positions)

	if err_binding != nil {
		log.Println(err_binding)
	}
	db, err_db_open := sql.Open("mysql", "root:Karma100%@tcp(localhost:3306)/alert_trade_db")

	if err_db_open != nil {
		log.Println(err_db_open)
	}
	_, err_call_stp_insert_user_position_details := db.Exec("call stp_insert_user_position_details(?,?,?,?,?,?)", user_positions.EntityID, user_positions.Ticker, user_positions.Quantity, user_positions.TradePrice, user_positions.TradeType, user_positions.Timestamp)
	if err_call_stp_insert_user_position_details != nil {
		log.Println(err_call_stp_insert_user_position_details)
	}
	db.Close()
}
func InsertionSort(arr []trade_object) {
	for i := 1; i < len(arr); i++ {
		key := arr[i].Timestamp
		key2 := arr[i]
		j := i - 1

		// Move elements of arr[0..i-1], that are greater than key,
		// to one position ahead of their current position
		for j >= 0 && arr[j].Timestamp > key {
			arr[j+1] = arr[j]
			j = j - 1
		}
		arr[j+1] = key2
	}
}
func get_equity_and_historical_performance(c *gin.Context) {
	var wg sync.WaitGroup
	log.SetOutput(log_file)

	var user_stock_list symbols_entity_details
	from, to := getFormattedTimes_and_one_year_ago()
	err := c.BindJSON(&user_stock_list)
	if err != nil {
		log.Println(err)
	}
	var stock_equity_map = make(map[string]trade_respnse_final)
	// var trade_list_map = make(map[string][]trade_object)
	for i := 0; i < len(user_stock_list.List_symbols); i++ {
		wg.Add(1)

		go func(symbol string, wg *sync.WaitGroup, entity_id int) {
			trade_list, timestamps := get_historical_data(symbol, from, to, "5min", 1)

			trade_list_final, final_timestamps := get_position_history_symbol(1, symbol, entity_id, trade_list, timestamps)
			fmt.Println(final_timestamps)
			// sort.Ints(final_timestamps)
			InsertionSort(trade_list_final)
			fmt.Println(len(trade_list_final))
			fmt.Println(trade_list_final)

			fmt.Println("check3")

			for d := 0; d < len(trade_list_final)-1; d++ {
				if trade_list_final[d].Timestamp == trade_list_final[d+1].Timestamp && trade_list_final[d].Closes != 0 {
					trade_list_final = append(trade_list_final[:d], trade_list_final[d+1:]...)
					d--
				}
			}

			fmt.Println("check1")
			fmt.Println(len(trade_list_final))
			sum := 0
			// trade_list_map[symbol+"_prices"] = trade_list_final
			var arr_equity = []float64{}
			for i, val := range trade_list_final {
				price := val.Position_price
				if price == 0 {
					price = val.Closes
				}
				if i == 0 {

					sum = val.Position_quantity
				} else if i > 0 {
					sum = sum + val.Position_quantity
				}
				var current_equity = float64(sum) * price
				arr_equity = append(arr_equity, current_equity)
				mutex.Lock()
				trade_response := trade_respnse_final{Equity: arr_equity, Trade_details: trade_list_final}
				stock_equity_map[symbol] = trade_response
				mutex.Unlock()

			}

			wg.Done()
			fmt.Println("check2")

		}(user_stock_list.List_symbols[i], &wg, user_stock_list.Entity_id)
	}

	for i := 0; i < len(user_stock_list.List_symbols); i++ {

		wg.Add(1)

		go func(symbol string, wg *sync.WaitGroup, entity_id int) {
			trade_list, timestamps := get_historical_data(symbol, from, to, "5min", 1)

			trade_list_final, final_timestamps := get_position_history_symbol(2, symbol, entity_id, trade_list, timestamps)
			fmt.Println(final_timestamps)
			// sort.Ints(final_timestamps)
			InsertionSort(trade_list_final)
			fmt.Println(len(trade_list_final))
			fmt.Println(trade_list_final)

			fmt.Println("check5")

			for d := 0; d < len(trade_list_final)-1; d++ {
				if trade_list_final[d].Timestamp == trade_list_final[d+1].Timestamp && trade_list_final[d].Closes != 0 {
					trade_list_final = append(trade_list_final[:d], trade_list_final[d+1:]...)
					d--
				}
			}

			fmt.Println("check6")
			fmt.Println(len(trade_list_final))

			// trade_list_map[symbol+"_prices"] = trade_list_final
			var arr_equity_2 = []float64{}
			var current_equity_2 float64
			for i, val := range trade_list_final {

				price := val.Position_price

				if price == 0 {
					price = val.Closes
				}
				if i == 0 {
					current_equity_2 = float64(val.Position_quantity) * val.Position_price
					// arr_equity_2 = append(arr_equity_2, current_equity_2)

				} else {
					var prev_price float64
					prev_equity := arr_equity_2[i-1]

					if trade_list_final[i-1].Position_price == 0 {
						prev_price = trade_list_final[i-1].Closes

					} else {
						prev_price = trade_list_final[i-1].Position_price
					}
					quantity := val.Position_quantity
					current_equity_2 = (prev_equity + (((prev_price - price) / prev_price) * prev_equity)) + (float64(quantity) * price)

				}
				// var current_equity = float64(sum) * price
				arr_equity_2 = append(arr_equity_2, current_equity_2)
				mutex.Lock()
				trade_response := trade_respnse_final{Equity: arr_equity_2, Trade_details: trade_list_final}
				stock_equity_map[symbol+"_short"] = trade_response
				mutex.Unlock()

			}

			wg.Done()
			fmt.Println("check4")

		}(user_stock_list.List_symbols[i], &wg, user_stock_list.Entity_id)

	}
	wg.Wait()

	// c.IndentedJSON(http.StatusOK, trade_list_final)
	// c.IndentedJSON(http.StatusOK, final_timestamps)
	// c.IndentedJSON(http.StatusOK, arr_equity)
	// var list_json = []interface{}{}
	// list_json = append(list_json, stock_equity_map, trade_list_map)

	c.IndentedJSON(http.StatusOK, stock_equity_map)

}
func get_historical_data(symbol string, from string, to string, interval string, api_id int) ([]trade_object, []int) {
	var acess_token string
	var expires_in int

	db, err_db_historical_data := sql.Open("mysql", "root:Karma100%@tcp(localhost:3306)/alert_trade_db")
	if err_db_historical_data != nil {
		log.Println(err_db_historical_data)
	}

	db.QueryRow("call stp_get_access_token_api_id(?)", api_id).Scan(&acess_token, &expires_in, &api_id)

	url := fmt.Sprintf("https://history.truedata.in/getbars?symbol=%s&from=%s&to=%s&response=csv&interval=%s", symbol, from, to, interval)
	req, err_request_historical_data := http.NewRequest("GET", url, nil)
	if err_request_historical_data != nil {
		log.Println("Error creating request:", err_request_historical_data)

	}

	req.Header.Add("Authorization", "Bearer "+acess_token)
	req.Header.Add("Accept", "application/json")
	fmt.Println(acess_token)

	res, err_response_historical_data := http.DefaultClient.Do(req)
	if err_response_historical_data != nil {
		log.Println(err_response_historical_data)
	}

	body, err_reading_body_historical_data := io.ReadAll(res.Body)
	if err_reading_body_historical_data != nil {
		log.Println(err_reading_body_historical_data)
	}

	body_string := string(body)

	df := dataframe.ReadCSV(strings.NewReader(body_string))
	fmt.Println(df)
	closes := df.Col("close")
	times := df.Col("timestamp")
	lows := df.Col("low")
	highs := df.Col("high")
	var trade_list []trade_object
	var list_timestamp []int
	for i := 0; i < closes.Len(); i++ {
		var struc trade_object
		close := closes.Elem(i).Float()
		timestr := times.Elem(i).String()
		low := lows.Elem(i).Float()
		high := highs.Elem(i).Float()
		const layout = "2006-01-02T15:04:05"
		t, err := time.Parse(layout, timestr)
		if err != nil {
			log.Fatal(err)
		}

		// Convert the time.Time object to a Unix timestamp
		unixTimestamp := t.Unix()
		list_timestamp = append(list_timestamp, int(unixTimestamp))

		struc.Closes = close
		struc.High = high
		struc.Low = low
		struc.Timestamp = int(unixTimestamp)
		trade_list = append(trade_list, struc)
	}
	return trade_list, list_timestamp

}

func getFormattedTimes_and_one_year_ago() (string, string) {

	timeFormat := "060102T15:04:05"

	currentTime := time.Now()

	oneYearAgo := currentTime.AddDate(-1, 0, 0)

	currentTimeStr := currentTime.Format(timeFormat)
	oneYearAgoStr := oneYearAgo.Format(timeFormat)
	fmt.Println(currentTimeStr)
	fmt.Println(oneYearAgoStr)
	return oneYearAgoStr, currentTimeStr
}

func get_position_history_symbol(position_type int, symbol string, entity_id int, trade_list []trade_object, list_timestamps []int) ([]trade_object, []int) {
	db, err_get_userpositions := sql.Open("mysql", "root:Karma100%@tcp(localhost:3306)/alert_trade_db")
	if err_get_userpositions != nil {
		log.Println(err_get_userpositions)
	}
	rows, err_db_get_position_detail := db.Query("call stp_get_user_position_details(?,?,?)", entity_id, position_type, symbol)
	if err_db_get_position_detail != nil {
		fmt.Println(err_db_get_position_detail)
	}
	var (
		EntityID   int
		Ticker     string
		Quantity   int
		TradePrice float64
		TradeType  int
		TradeID    int
		Timestamp  int
	)

	for rows.Next() {
		var struc trade_object
		rows.Scan(&EntityID, &Ticker, &Quantity, &TradePrice, &TradeType, &TradeID, &Timestamp)

		struc.Timestamp = Timestamp
		list_timestamps = append(list_timestamps, Timestamp)
		struc.Position_quantity = int(Quantity)

		struc.Position_price = TradePrice
		trade_list = append(trade_list, struc)

	}
	return trade_list, list_timestamps
}
