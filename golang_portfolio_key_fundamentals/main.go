package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
)

type portfolio_json struct {
	stocks_cocode []float64
	stock_value   []float64
}

func main() {
	file, err := os.Create("golang_portfolio_key_fundamentals_log.log")
	if err != nil {
		fmt.Println(err)
	}
	log.SetOutput(file)
	router := gin.Default()
	router.POST("/portfolio_fundamental_sentiment", portfolio_fundamental_sentiment)
	router.Run("localhost:8080")
}

func portfolio_fundamental_sentiment(c *gin.Context) {
	var stockdetails portfolio_json

	err := c.BindJSON(&stockdetails)
	if err != nil {
		log.Println(err)
	}
	total_porfolti_val := 0.0
	var list_stock_alocations []float64
	for i := 0; i < len(stockdetails.stocks_cocode); i++ {
		total_porfolti_val = total_porfolti_val + stockdetails.stock_value[i]

	}
	for i := 0; i < len(stockdetails.stocks_cocode); i++ {
		alocation := stockdetails.stock_value[i] / total_porfolti_val
		list_stock_alocations = append(list_stock_alocations, alocation)
	}
	var wg sync.WaitGroup
	for i := 0; i < len(stockdetails.stocks_cocode); i++ {
		wg.Add(1)
		go func(i int, wg1 *sync.WaitGroup) {
			url := fmt.Sprintf("https://uatapi-keyindctrs.traderspilot.com/get_fundamental_sentiment?cocode=%v", stockdetails.stocks_cocode[i])
			res, err := http.Get(url)
			if err != nil {
				log.Println(err)
			}
			body, err := io.ReadAll(res.Body)
			fmt.Println(string(body))
			wg1.Done()
		}(i, &wg)

	}
	wg.Wait()

}
