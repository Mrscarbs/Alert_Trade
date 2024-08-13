package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-gota/gota/dataframe"
	_ "github.com/go-sql-driver/mysql"
)

var log_file, _ = os.Create("backtester_log.log")

type Condition struct {
	Indicator            string
	Parameters           []int  `json:"parameters"`
	Logic                string `json:"logic"`
	Compare              []int  `json:"compare"` // Use interface{} to handle both strings and string arrays
	Indicator_calculated [][]float64
}

// Define the main structure for buy and sell positions
type TradingStrategy struct {
	BuyPosition  []Condition `json:"buy_position"`
	SellPosition []Condition `json:"sell_position"`
}

// type composite_condition struct {
// 	indicator       string
// 	indicator_value [][]float64
// 	compare         [][]float64
// }

type length_master_possition struct {
	arr []int
	len int
}

func main() {
	router := gin.Default()
	router.POST("/strategy_creator", strategy_creator)
	router.Run("localhost:8080")
}
func strategy_creator(c *gin.Context) {
	var strategy_json TradingStrategy
	c.BindJSON(&strategy_json)
	dataframe := get_price_data("63MOONS", "60min", 1)
	fmt.Println(dataframe)
	fmt.Println(strategy_json)
	master_indicator_map := map_creator()
	prices := dataframe.Col("close").Float()

	for i := 0; i < len(strategy_json.BuyPosition); i++ {
		fmt.Println(master_indicator_map[strategy_json.BuyPosition[i].Indicator])
		indicator_func := master_indicator_map[strategy_json.BuyPosition[i].Indicator].function
		indicator := indicator_func(prices, strategy_json.BuyPosition[i].Parameters)
		fmt.Println(indicator)
		indicator = append(indicator, prices)
		strategy_json.BuyPosition[i].Indicator_calculated = indicator
	}
	for i := 0; i < len(strategy_json.SellPosition); i++ {
		fmt.Println(master_indicator_map[strategy_json.SellPosition[i].Indicator])
		indicator_func := master_indicator_map[strategy_json.SellPosition[i].Indicator].function
		indicator := indicator_func(prices, strategy_json.SellPosition[i].Parameters)
		fmt.Println(indicator)
		indicator = append(indicator, prices)
		strategy_json.SellPosition[i].Indicator_calculated = indicator
	}

	fmt.Println(strategy_json)
	fmt.Println("=====================================================================================================================")
	map_positions := make(map[string][]int)
	var total_position_map [][]int
	for i := 0; i < len(strategy_json.BuyPosition); i++ {
		if strategy_json.BuyPosition[i].Logic == "crossover" {
			t := strconv.Itoa(i)
			posistion_indexes := Crossover(strategy_json.BuyPosition[i].Indicator_calculated, strategy_json.BuyPosition[i].Compare)
			map_positions["position"+t] = posistion_indexes
			total_position_map = append(total_position_map, posistion_indexes)
		}
		if strategy_json.BuyPosition[i].Logic == ">" {
			t := strconv.Itoa(i)
			posistion_indexes := Greater_than(strategy_json.BuyPosition[i].Indicator_calculated, strategy_json.BuyPosition[i].Compare)
			map_positions["position_"+t] = posistion_indexes
			total_position_map = append(total_position_map, posistion_indexes)

		}
		if strategy_json.BuyPosition[i].Logic == "<" {
			t := strconv.Itoa(i)
			posistion_indexes := Less_than(strategy_json.BuyPosition[i].Indicator_calculated, strategy_json.BuyPosition[i].Compare)
			map_positions["position_"+t] = posistion_indexes
			total_position_map = append(total_position_map, posistion_indexes)

		}
	}

	fmt.Println("-------------------------------------------------------------------------------------------")
	fmt.Println(map_positions)
	fmt.Println(total_position_map)
	length_indicatoras := []length_master_possition{}
	for i := 0; i < len(total_position_map); i++ {

		obj := length_master_possition{arr: total_position_map[i], len: len(total_position_map[i])}
		length_indicatoras = append(length_indicatoras, obj)
	}
	quick_sort(0, len(length_indicatoras)-1, length_indicatoras)
	fmt.Println(length_indicatoras)
	main_check_list := length_indicatoras[0].arr
	for i := 0; i < len(main_check_list); i++ {
		check_target := main_check_list[i]
		for z := 1; z < len(length_indicatoras); z++ {
			ret := binary_search(length_indicatoras[z].arr, check_target)
			if ret == -1 {
				main_check_list[i] = 0
			}
		}

	}
	fmt.Println(main_check_list)
}

func binary_search(arr []int, target int) int {
	low := 0
	high := len(arr) - 1
	for low <= high {
		middle := low + (high-low)/2
		middle_val := arr[middle]

		if target > middle_val {
			low = middle + 1
		} else if target < middle_val {
			high = middle - 1
		} else {
			return middle
		}

	}
	return -1
}

func partition_quick_sort(start int, end int, arr []length_master_possition) int {
	pivot := arr[end]
	i := start - 1
	for j := start; j < end; j++ {
		if arr[j].len < pivot.len {
			i++
			temp := arr[i]
			arr[i] = arr[j]
			arr[j] = temp
		}
	}
	i++
	temp := arr[i]
	arr[i] = arr[end]
	arr[end] = temp
	return i
}

func quick_sort(start int, end int, arr []length_master_possition) {
	if end <= start {
		return
	}

	pivot := partition_quick_sort(start, end, arr)
	quick_sort(start, pivot-1, arr)
	quick_sort(pivot+1, end, arr)

}

func get_price_data(symbol string, interval string, api_id int) dataframe.DataFrame {
	var acess_token string
	var expires_in int
	log.SetOutput(log_file)
	db, err_db_historical_data := sql.Open("mysql", "root:Karma100%@tcp(localhost:3306)/alert_trade_db")
	if err_db_historical_data != nil {
		log.Println(err_db_historical_data)
	}
	from, to := getFormattedTimes_and_one_year_ago()
	db.QueryRow("call stp_get_access_token_api_id(?)", 1).Scan(&acess_token, &expires_in, &api_id)
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
	// fmt.Println(string(body))
	fmt.Println(string(body))
	df := dataframe.ReadCSV(strings.NewReader(string(body)))
	// fmt.Println(df)
	return df

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
