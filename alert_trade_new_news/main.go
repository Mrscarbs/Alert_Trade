package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type news_reponse_form struct {
	Kind string `json:"kind"`
	URL  struct {
		Type     string `json:"type"`
		Template string `json:"template"`
	} `json:"url"`
	Queries struct {
		Request []struct {
			Title          string `json:"title"`
			TotalResults   string `json:"totalResults"`
			SearchTerms    string `json:"searchTerms"`
			Count          int    `json:"count"`
			StartIndex     int    `json:"startIndex"`
			InputEncoding  string `json:"inputEncoding"`
			OutputEncoding string `json:"outputEncoding"`
			Safe           string `json:"safe"`
			Cx             string `json:"cx"`
		} `json:"request"`
		NextPage []struct {
			Title          string `json:"title"`
			TotalResults   string `json:"totalResults"`
			SearchTerms    string `json:"searchTerms"`
			Count          int    `json:"count"`
			StartIndex     int    `json:"startIndex"`
			InputEncoding  string `json:"inputEncoding"`
			OutputEncoding string `json:"outputEncoding"`
			Safe           string `json:"safe"`
			Cx             string `json:"cx"`
		} `json:"nextPage"`
	} `json:"queries"`
	Context struct {
		Title string `json:"title"`
	} `json:"context"`
	SearchInformation struct {
		SearchTime            float64 `json:"searchTime"`
		FormattedSearchTime   string  `json:"formattedSearchTime"`
		TotalResults          string  `json:"totalResults"`
		FormattedTotalResults string  `json:"formattedTotalResults"`
	} `json:"searchInformation"`
	Items []struct {
		Kind             string `json:"kind"`
		Title            string `json:"title"`
		HTMLTitle        string `json:"htmlTitle"`
		Link             string `json:"link"`
		DisplayLink      string `json:"displayLink"`
		Snippet          string `json:"snippet"`
		HTMLSnippet      string `json:"htmlSnippet"`
		FormattedURL     string `json:"formattedUrl"`
		HTMLFormattedURL string `json:"htmlFormattedUrl"`
		Pagemap          struct {
			CseThumbnail []struct {
				Src    string `json:"src"`
				Width  string `json:"width"`
				Height string `json:"height"`
			} `json:"cse_thumbnail"`
			Metatags []struct {
				OgImage                       string `json:"og:image"`
				OgType                        string `json:"og:type"`
				TwitterCard                   string `json:"twitter:card"`
				TwitterTitle                  string `json:"twitter:title"`
				TwitterURL                    string `json:"twitter:url"`
				Author                        string `json:"author"`
				OgTitle                       string `json:"og:title"`
				DailymotionDomainVerification string `json:"dailymotion-domain-verification"`
				OgDescription                 string `json:"og:description"`
				TwitterCreator                string `json:"twitter:creator"`
				TwitterImageSrc               string `json:"twitter:image:src"`
				TwitterSite                   string `json:"twitter:site"`
				Viewport                      string `json:"viewport"`
				TwitterDescription            string `json:"twitter:description"`
				OgURL                         string `json:"og:url"`
			} `json:"metatags"`
			CseImage []struct {
				Src string `json:"src"`
			} `json:"cse_image"`
		} `json:"pagemap"`
	} `json:"items"`
}

var logfile, _ = os.Create("news_logs.log")

func main() {
	router := gin.Default()
	router.GET("/get_stockwise_news", get_stockwise_news)
	router.Run("0.0.0.0:8080")
}

func get_stockwise_news(c *gin.Context) {
	stock, _ := c.GetQuery("stock")
	api, search_id := get_search_id_and_api_key(logfile)
	news_res := get_stock_news(stock+" latest india buisness news", search_id, api)
	c.IndentedJSON(http.StatusOK, news_res)

}
func get_search_id_and_api_key(log_file *os.File) (string, string) {
	var api_id int
	var provider string
	var api_key string
	var search_endine_id string
	var start_time time.Time
	var last_updated_time time.Time

	log.SetOutput(log_file)
	db, err := sql.Open("mysql", "admin:saumitrasuparn@tcp(alerttradedb.czqug0e2in8p.ap-south-1.rds.amazonaws.com:3306)/alert_trade_db")
	if err != nil {
		log.Println(err)
	}
	db.QueryRow("call alert_trade_db.stp_get_api_config(?)", 3).Scan(&api_id, &provider, &api_key, &search_endine_id, &start_time, &last_updated_time)
	defer db.Close()
	fmt.Println(api_key)
	fmt.Println("-------------------------------")

	fmt.Println(search_endine_id)
	fmt.Println("----------------------------------")
	return api_key, search_endine_id
}

func get_stock_news(query string, search_engine_id string, api_key string) news_reponse_form {
	var news_response news_reponse_form
	encodedQuery := url.QueryEscape(query)
	var url string = fmt.Sprintf("https://www.googleapis.com/customsearch/v1?q=%s&cx=%s&key=%s", encodedQuery, search_engine_id, api_key)

	res, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(string(body))
	err = json.Unmarshal(body, &news_response)

	if err != nil {
		log.Println(err)
	}

	return news_response
}
