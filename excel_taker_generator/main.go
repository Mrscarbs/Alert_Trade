package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/xuri/excelize/v2"
)

func main() {

	file, err := os.Create("excel_taker.log")
	log.SetOutput(file)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("exceltaker initializing")
	companies := get_excel_column_extractor("Company")
	date := get_excel_column_extractor("Date")
	price := get_excel_column_extractor("Price")
	Quantity := get_excel_column_extractor("Quantity")
	time := get_excel_column_extractor("Trade Time")
	Side := get_excel_column_extractor("Side")
	fmt.Println(companies)
	fmt.Println(date)
	fmt.Println(price)
	fmt.Println(Quantity)
	fmt.Println(time)
	fmt.Println(Side)
	unix_times := []int64{}
	for i := 0; i < len(date); i++ {
		unix_time := combineDateAndTime(date[i], time[i])
		unix_times = append(unix_times, unix_time)
	}
	fmt.Println(unix_times)
	var side_num = 0
	for i := 0; i < len(companies); i++ {
		db, err := sql.Open("mysql", "root:Karma100%@tcp(localhost:3306)/alert_trade_db")
		if err != nil {
			log.Println(err)
		}

		if Side[i] == "Buy" {
			side_num = 1
		} else {
			side_num = 2
		}
		cleanedValue := strings.Replace(price[i], "₹", "", -1)
		cleanedValue = strings.Replace(cleanedValue, ",", "", -1)

		// Convert the cleaned string to a float64
		floatValue_price, err := strconv.ParseFloat(cleanedValue, 64)
		if err != nil {
			log.Println(err)
		}
		_, err_db := db.Exec("call stp_insert_user_position_details(?,?,?,?,?,?)", 2, companies[i], Quantity[i], floatValue_price, side_num, unix_times[i])
		if err != nil {
			log.Println(err_db)
		}
	}

}
func get_excel_column_extractor(name string) []string {

	f, err := excelize.OpenFile("trade_2425_6XA39S (1).xlsx")
	if err != nil {
		log.Println(err)
	}

	rows, err := f.GetRows("TRADE")

	if err != nil {
		log.Println(err)
	}
	var companies = []string{}
	var star_index int
	var company_index int
	for index, row := range rows {
		fmt.Println(row)
		for i := 0; i < len(row); i++ {
			if row[i] == name {
				company_index = i
				star_index = index
				break
			}

		}

	}

	for index, row := range rows {
		if index > star_index && len(row) == 0 {
			break
		}
		if len(row) > company_index && index > star_index {

			companies = append(companies, row[company_index])
		}

	}
	fmt.Println("hiiiiiii:  ", companies)
	return companies
}
func combineDateAndTime(dateStr, timeStr string) int64 {
	// Layout for parsing the date and time
	dateLayout := "02-01-2006"
	timeLayout := "15:04:05"

	// Parse the date
	date, err := time.Parse(dateLayout, dateStr)
	if err != nil {
		log.Println("Error parsing date:", err)
		return 0
	}

	// Parse the time
	parsedTime, err := time.Parse(timeLayout, timeStr)
	if err != nil {
		log.Println("Error parsing time:", err)
		return 0
	}

	// Combine the date and time
	combinedDateTime := time.Date(date.Year(), date.Month(), date.Day(), parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(), 0, time.UTC)

	// Return the Unix timestamp
	return combinedDateTime.Unix()
}
