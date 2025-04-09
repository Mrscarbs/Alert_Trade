package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type CombinedFinancialData struct {
	Company_name                 string
	CO_CODE                      int     `json:"co_code"`
	TTMAson                      int     `json:"ttm_ason"`
	MCAP1                        float64 `json:"mcap1"`
	EV1                          float64 `json:"ev1"`
	PE1                          float64 `json:"pe1"`
	PBV1                         float64 `json:"pbv1"`
	DIVYIELD1                    float64 `json:"div_yield1"`
	EPS1                         float64 `json:"eps1"`
	BookValue1                   float64 `json:"book_value1"`
	ROA_TTM1                     float64 `json:"roa_ttm1"`
	ROE_TTM1                     float64 `json:"roe_ttm1"`
	ROCE_TTM1                    float64 `json:"roce_ttm1"`
	EBIT_TTM1                    float64 `json:"ebit_ttm1"`
	EBITDA_TTM1                  float64 `json:"ebitda_ttm1"`
	EV_Sales_TTM1                float64 `json:"ev_sales_ttm1"`
	EV_EBITDA_TTM1               float64 `json:"ev_ebitda_ttm1"`
	NetIncomeMargin_TTM1         float64 `json:"net_income_margin_ttm1"`
	GrossIncomeMargin_TTM1       float64 `json:"gross_income_margin_ttm1"`
	AssetTurnover_TTM1           float64 `json:"asset_turnover_ttm1"`
	CurrentRatio_TTM1            float64 `json:"current_ratio_ttm1"`
	Debt_Equity_TTM1             float64 `json:"debt_equity_ttm1"`
	Sales_TotalAssets_TTM1       float64 `json:"sales_total_assets_ttm1"`
	NetDebt_EBITDA_TTM1          float64 `json:"net_debt_ebitda_ttm1"`
	EBITDA_Margin_TTM1           float64 `json:"ebitda_margin_ttm1"`
	TotalShareHoldersEquity_TTM1 float64 `json:"total_shareholders_equity_ttm1"`
	ShorttermDebt_TTM1           float64 `json:"short_term_debt_ttm1"`
	LongtermDebt_TTM1            float64 `json:"long_term_debt_ttm1"`
	SharesOutstanding1           float64 `json:"shares_outstanding1"`
	EPSDiluted1                  float64 `json:"eps_diluted1"`
	NetSales1                    float64 `json:"net_sales1"`
	Netprofit1                   float64 `json:"net_profit1"`
	AnnualDividend1              float64 `json:"annual_dividend1"`
	COGS1                        float64 `json:"cogs1"`
	PEGRatio_TTM1                float64 `json:"peg_ratio_ttm1"`
	DividendPayout_TTM1          float64 `json:"dividend_payout_ttm1"`
	MCAP2                        float64 `json:"mcap2"`
	EV2                          float64 `json:"ev2"`
	PE2                          float64 `json:"pe2"`
	PBV2                         float64 `json:"pbv2"`
	EPS2                         float64 `json:"eps2"`
	BookValue2                   float64 `json:"book_value2"`
	EBIT2                        float64 `json:"ebit2"`
	EBITDA2                      float64 `json:"ebitda2"`
	EV_Sales2                    float64 `json:"ev_sales2"`
	EV_EBITDA2                   float64 `json:"ev_ebitda2"`
	NetIncomeMargin2             float64 `json:"net_income_margin2"`
	GrossIncomeMargin2           float64 `json:"gross_income_margin2"`
	EBITDAMargin2                float64 `json:"ebitda_margin2"`
	EPSDiluted2                  float64 `json:"eps_diluted2"`
	NetSales2                    float64 `json:"net_sales2"`
	Netprofit2                   float64 `json:"net_profit2"`
	COGS2                        float64 `json:"cogs2"`
	MCAP3                        float64 `json:"mcap3"`
	EV3                          float64 `json:"ev3"`
	PE3                          float64 `json:"pe3"`
	PBV3                         float64 `json:"pbv3"`
	EPS3                         float64 `json:"eps3"`
	BookValue3                   float64 `json:"book_value3"`
	EBIT3                        float64 `json:"ebit3"`
	EBITDA3                      float64 `json:"ebitda3"`
	EV_Sales3                    float64 `json:"ev_sales3"`
	EV_EBITDA3                   float64 `json:"ev_ebitda3"`
	NetIncomeMargin3             float64 `json:"net_income_margin3"`
	GrossIncomeMargin3           float64 `json:"gross_income_margin3"`
	EBITDAMargin3                float64 `json:"ebitda_margin3"`
	EPSDiluted3                  float64 `json:"eps_diluted3"`
	NetSales3                    float64 `json:"net_sales3"`
	Netprofit3                   float64 `json:"net_profit3"`
	COGS3                        float64 `json:"cogs3"`
	MCAP4                        float64 `json:"mcap4"`
	EV4                          float64 `json:"ev4"`
	PE4                          float64 `json:"pe4"`
	PBV4                         float64 `json:"pbv4"`
	EPS4                         float64 `json:"eps4"`
	BookValue4                   float64 `json:"book_value4"`
	EBIT4                        float64 `json:"ebit4"`
	EBITDA4                      float64 `json:"ebitda4"`
	EV_Sales4                    float64 `json:"ev_sales4"`
	EV_EBITDA4                   float64 `json:"ev_ebitda4"`
	NetIncomeMargin4             float64 `json:"net_income_margin4"`
	GrossIncomeMargin4           float64 `json:"gross_income_margin4"`
	EBITDAMargin4                float64 `json:"ebitda_margin4"`
	EPSDiluted4                  float64 `json:"eps_diluted4"`
	NetSales4                    float64 `json:"net_sales4"`
	Netprofit4                   float64 `json:"net_profit4"`
	COGS4                        float64 `json:"cogs4"`
}
type upstox_response struct {
	Broker_name string
	Status      string `json:"status"`
	Data        []struct {
		Isin                     string  `json:"isin"`
		CncUsedQuantity          int     `json:"cnc_used_quantity"`
		CollateralType           string  `json:"collateral_type"`
		CompanyName              string  `json:"company_name"`
		Haircut                  float64 `json:"haircut"`
		Product                  string  `json:"product"`
		Quantity                 int     `json:"quantity"`
		Tradingsymbol            string  `json:"tradingsymbol"`
		LastPrice                float64 `json:"last_price"`
		ClosePrice               float64 `json:"close_price"`
		Pnl                      float64 `json:"pnl"`
		DayChange                float64 `json:"day_change"`
		DayChangePercentage      float64 `json:"day_change_percentage"`
		InstrumentToken          string  `json:"instrument_token"`
		AveragePrice             float64 `json:"average_price"`
		CollateralQuantity       int     `json:"collateral_quantity"`
		CollateralUpdateQuantity int     `json:"collateral_update_quantity"`
		TradingSymbol            string  `json:"trading_symbol"`
		T1Quantity               int     `json:"t1_quantity"`
		Exchange                 string  `json:"exchange"`
	} `json:"data"`
}

type user_form struct {
	sub_user_id        int64
	acess_token        string
	broker_name_string string
	api_key            string
}
type angle_one_resp struct {
	Broker_name string
	Status      bool   `json:"status"`
	Message     string `json:"message"`
	Errorcode   string `json:"errorcode"`
	Data        struct {
		Holdings []struct {
			Tradingsymbol      string      `json:"tradingsymbol"`
			Exchange           string      `json:"exchange"`
			Isin               string      `json:"isin"`
			T1Quantity         int         `json:"t1quantity"`
			Realisedquantity   int         `json:"realisedquantity"`
			Quantity           int         `json:"quantity"`
			Authorisedquantity int         `json:"authorisedquantity"`
			Product            string      `json:"product"`
			Collateralquantity interface{} `json:"collateralquantity"`
			Collateraltype     interface{} `json:"collateraltype"`
			Haircut            float64     `json:"haircut"`
			Averageprice       float64     `json:"averageprice"`
			Ltp                float64     `json:"ltp"`
			Symboltoken        string      `json:"symboltoken"`
			Close              float64     `json:"close"`
			Profitandloss      float64     `json:"profitandloss"`
			Pnlpercentage      float64     `json:"pnlpercentage"`
		} `json:"holdings"`
		Totalholding struct {
			Totalholdingvalue  int     `json:"totalholdingvalue"`
			Totalinvvalue      int     `json:"totalinvvalue"`
			Totalprofitandloss float64 `json:"totalprofitandloss"`
			Totalpnlpercentage float64 `json:"totalpnlpercentage"`
		} `json:"totalholding"`
	} `json:"data"`
}
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RequestBody struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type ResponseBody struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}
type llm_response struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role        string        `json:"role"`
			Content     string        `json:"content"`
			Refusal     interface{}   `json:"refusal"`
			Annotations []interface{} `json:"annotations"`
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
	ServiceTier string `json:"service_tier"`
}

func main() {
	var err error
	fmt.Println("llm_api_2")
	file, err_file := os.Create("alert_trade_llm_api2_log.log")
	if err_file != nil {
		log.Println(err_file)
	}
	log.SetOutput(file)
	db, err = sql.Open("mysql", "admin:saumitrasuparn@tcp(alerttradedb.czqug0e2in8p.ap-south-1.rds.amazonaws.com:3306)/alert_trade_db")
	if err != nil {
		log.Println(err)
	}
	router := gin.Default()
	router.GET("/chat_with_portfolio", llm_api_final)
	router.Run("0.0.0.0:8080")

}

func llm_api_final(c *gin.Context) {
	query, _ := c.GetQuery("query")
	eid, _ := c.GetQuery("eid")
	eid_int, err := strconv.Atoi(eid)
	if err != nil {
		log.Println(err)
	}
	list_user_all_details := []any{}
	user_list := get_portfolio_details(eid_int)
	_, user_holdings, financial_para := consolidated_acoount_generator(user_list)
	final_json, err := json.Marshal(user_holdings)
	if err != nil {
		log.Println(err)
	}
	final_json_string := string(final_json)
	list_user_all_details = append(list_user_all_details, financial_para, final_json_string)
	llm_resp := invoke_llm(query, list_user_all_details)
	c.IndentedJSON(http.StatusOK, llm_resp)

}

func get_portfolio_details(user_id int) []user_form {
	var user_details_llist []user_form
	rows, err := db.Query("call stp_GetSubUserandToken(?)", user_id)
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		var sub_user user_form
		rows.Scan(&sub_user.sub_user_id, &sub_user.acess_token, &sub_user.broker_name_string, &sub_user.api_key)
		user_details_llist = append(user_details_llist, sub_user)
	}
	fmt.Println(user_details_llist)
	return user_details_llist
}

func get_upstox_holdings(acess_token string) upstox_response {
	url := "https://api.upstox.com/v2/portfolio/long-term-holdings"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)

	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+acess_token)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)

	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)

	}
	fmt.Println(body)
	var res_bod upstox_response
	json.Unmarshal(body, &res_bod)
	fmt.Println(res_bod)
	res_bod.Broker_name = "upstox"
	return res_bod
}

func get_angle_one_holdings(acess_token string, api_key string) angle_one_resp {
	url := "https://apiconnect.angelone.in/rest/secure/angelbroking/portfolio/v1/getAllHolding"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)

	}
	mac_address := getMACAddress()
	fmt.Println(mac_address)
	req.Header.Add("Authorization", "Bearer "+acess_token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-UserType", "USER")
	req.Header.Add("X-SourceID", "WEB")
	req.Header.Add("X-MACAddress", mac_address)
	req.Header.Add("X-PrivateKey", api_key)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)

	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)

	}
	fmt.Println(string(body))
	var resp angle_one_resp
	json.Unmarshal(body, &resp)
	resp.Broker_name = "angleone"
	return resp
}

func getMACAddress() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "00:00:00:00:00:00"
	}
	for _, i := range interfaces {
		if len(i.HardwareAddr) > 0 {
			return i.HardwareAddr.String()
		}
	}
	return "00:00:00:00:00:00"
}
func consolidated_acoount_generator(user_list []user_form) ([]CombinedFinancialData, []any, string) {
	list_isins := []string{}
	var wg sync.WaitGroup

	var list_user_holding_details []any
	for i := 0; i < len(user_list); i++ {

		sub_user := user_list[i]
		wg.Add(1)
		go func(sub_user_in user_form, wg *sync.WaitGroup) {
			defer wg.Done()
			if sub_user_in.broker_name_string == "upstox" {
				fmt.Println("upstox")
				res := get_upstox_holdings(sub_user_in.acess_token)
				list_user_holding_details = append(list_user_holding_details, res)
				for i := 0; i < len(res.Data); i++ {
					instrument_token := res.Data[i].InstrumentToken
					instrument_token_formated := strings.TrimPrefix(instrument_token, "NSE_EQ|")
					list_isins = append(list_isins, instrument_token_formated)
				}
			}
			if sub_user_in.broker_name_string == "angelone" {
				fmt.Println("angelone")
				angleone_portfolio := get_angle_one_holdings(sub_user_in.acess_token, sub_user.api_key)
				list_user_holding_details = append(list_user_holding_details, angleone_portfolio)
				for i := 0; i < len(angleone_portfolio.Data.Holdings); i++ {

					instrument_token := angleone_portfolio.Data.Holdings[i].Isin
					list_isins = append(list_isins, instrument_token)
				}
			}

		}(sub_user, &wg)

	}

	wg.Wait()
	fmt.Println(list_isins)
	list_cocode := []int{}
	list_combined_financial_data := []CombinedFinancialData{}
	for i := 0; i < len(list_isins); i++ {
		var data CombinedFinancialData
		var cocode int
		err := db.QueryRow("call stp_GetCoCodeByISIN_cmots(?)", list_isins[i]).Scan(&cocode)
		if err != nil {
			log.Println(err)
		}
		var company_name string
		db.QueryRow("call stp_GetCompanyNameByCoCode(?)", cocode).Scan(&company_name)
		data.Company_name = company_name
		list_cocode = append(list_cocode, cocode)
		data.CO_CODE = cocode
		err_db := db.QueryRow("call stp_get_combined_financial_data_by_cocode_3(?)", cocode).Scan(
			&data.TTMAson,
			&data.MCAP1, &data.EV1, &data.PE1, &data.PBV1, &data.DIVYIELD1, &data.EPS1,
			&data.BookValue1, &data.ROA_TTM1, &data.ROE_TTM1, &data.ROCE_TTM1, &data.EBIT_TTM1,
			&data.EBITDA_TTM1, &data.EV_Sales_TTM1, &data.EV_EBITDA_TTM1, &data.NetIncomeMargin_TTM1,
			&data.GrossIncomeMargin_TTM1, &data.AssetTurnover_TTM1, &data.CurrentRatio_TTM1,
			&data.Debt_Equity_TTM1, &data.Sales_TotalAssets_TTM1, &data.NetDebt_EBITDA_TTM1,
			&data.EBITDA_Margin_TTM1, &data.TotalShareHoldersEquity_TTM1, &data.ShorttermDebt_TTM1,
			&data.LongtermDebt_TTM1, &data.SharesOutstanding1, &data.EPSDiluted1, &data.NetSales1,
			&data.Netprofit1, &data.AnnualDividend1, &data.COGS1, &data.PEGRatio_TTM1,
			&data.DividendPayout_TTM1, &data.MCAP2, &data.EV2, &data.PE2, &data.PBV2,
			&data.EPS2, &data.BookValue2, &data.EBIT2, &data.EBITDA2, &data.EV_Sales2,
			&data.EV_EBITDA2, &data.NetIncomeMargin2, &data.GrossIncomeMargin2, &data.EBITDAMargin2,
			&data.EPSDiluted2, &data.NetSales2, &data.Netprofit2, &data.COGS2,
			&data.MCAP3, &data.EV3, &data.PE3, &data.PBV3, &data.EPS3, &data.BookValue3,
			&data.EBIT3, &data.EBITDA3, &data.EV_Sales3, &data.EV_EBITDA3, &data.NetIncomeMargin3,
			&data.GrossIncomeMargin3, &data.EBITDAMargin3, &data.EPSDiluted3, &data.NetSales3,
			&data.Netprofit3, &data.COGS3, &data.MCAP4, &data.EV4, &data.PE4,
			&data.PBV4, &data.EPS4, &data.BookValue4, &data.EBIT4, &data.EBITDA4,
			&data.EV_Sales4, &data.EV_EBITDA4, &data.NetIncomeMargin4, &data.GrossIncomeMargin4,
			&data.EBITDAMargin4, &data.EPSDiluted4, &data.NetSales4, &data.Netprofit4,
			&data.COGS4,
		)
		list_combined_financial_data = append(list_combined_financial_data, data)
		if err_db != nil {
			log.Println(err_db)
		}
	}
	fmt.Println(list_cocode)
	fmt.Println(list_combined_financial_data)
	fmt.Println(list_user_holding_details)
	var financials_para string
	financials_para = "these are all the stocks and their fundamental details from the users portfolio "

	for i := 0; i < len(list_combined_financial_data); i++ {
		data_string_stock := fmt.Sprintf(
			"the cocode for %s is %v,the ttm_ason for company %s is %v,the mcap1 for company %s is %v,"+
				"the ev1 for company %s is %v,the pe1 for company %s is %v,the pbv1 for company %s is %v,"+
				"the div_yield1 for company %s is %v,the eps1 for company %s is %v,the book_value1 for company %s is %v,"+
				"the roa_ttm1 for company %s is %v,the roe_ttm1 for company %s is %v,the roce_ttm1 for company %s is %v,"+
				"the ebit_ttm1 for company %s is %v,the ebitda_ttm1 for company %s is %v,the ev_sales_ttm1 for company %s is %v,"+
				"the ev_ebitda_ttm1 for company %s is %v,the net_income_margin_ttm1 for company %s is %v,"+
				"the gross_income_margin_ttm1 for company %s is %v,the asset_turnover_ttm1 for company %s is %v,"+
				"the current_ratio_ttm1 for company %s is %v,the debt_equity_ttm1 for company %s is %v,"+
				"the sales_total_assets_ttm1 for company %s is %v,the net_debt_ebitda_ttm1 for company %s is %v,"+
				"the ebitda_margin_ttm1 for company %s is %v,the total_shareholders_equity_ttm1 for company %s is %v,"+
				"the short_term_debt_ttm1 for company %s is %v,the long_term_debt_ttm1 for company %s is %v,"+
				"the shares_outstanding1 for company %s is %v,the eps_diluted1 for company %s is %v,"+
				"the net_sales1 for company %s is %v,the net_profit1 for company %s is %v,"+
				"the annual_dividend1 for company %s is %v,the cogs1 for company %s is %v,"+
				"the peg_ratio_ttm1 for company %s is %v,the dividend_payout_ttm1 for company %s is %v,",
			list_combined_financial_data[i].Company_name, list_combined_financial_data[i].CO_CODE,
			list_combined_financial_data[i].Company_name, list_combined_financial_data[i].TTMAson,
			list_combined_financial_data[i].Company_name, list_combined_financial_data[i].MCAP1,
			list_combined_financial_data[i].Company_name, list_combined_financial_data[i].EV1,
			list_combined_financial_data[i].Company_name, list_combined_financial_data[i].PE1,
			list_combined_financial_data[i].Company_name, list_combined_financial_data[i].PBV1,
			list_combined_financial_data[i].Company_name, list_combined_financial_data[i].DIVYIELD1,
			list_combined_financial_data[i].Company_name, list_combined_financial_data[i].EPS1,
			list_combined_financial_data[i].Company_name, list_combined_financial_data[i].BookValue1,
			list_combined_financial_data[i].Company_name, list_combined_financial_data[i].ROA_TTM1,
			list_combined_financial_data[i].Company_name, list_combined_financial_data[i].ROE_TTM1,
			list_combined_financial_data[i].Company_name, list_combined_financial_data[i].ROCE_TTM1,
			list_combined_financial_data[i].Company_name, list_combined_financial_data[i].EBIT_TTM1,
			list_combined_financial_data[i].Company_name, list_combined_financial_data[i].EBITDA_TTM1,
			list_combined_financial_data[i].Company_name, list_combined_financial_data[i].EV_Sales_TTM1,
			list_combined_financial_data[i].Company_name, list_combined_financial_data[i].EV_EBITDA_TTM1,
			list_combined_financial_data[i].Company_name, list_combined_financial_data[i].NetIncomeMargin_TTM1,
			list_combined_financial_data[i].Company_name, list_combined_financial_data[i].GrossIncomeMargin_TTM1,
			list_combined_financial_data[i].Company_name, list_combined_financial_data[i].AssetTurnover_TTM1,
			list_combined_financial_data[i].Company_name, list_combined_financial_data[i].CurrentRatio_TTM1,
			list_combined_financial_data[i].Company_name, list_combined_financial_data[i].Debt_Equity_TTM1,
			list_combined_financial_data[i].Company_name, list_combined_financial_data[i].Sales_TotalAssets_TTM1,
			list_combined_financial_data[i].Company_name, list_combined_financial_data[i].NetDebt_EBITDA_TTM1,
			list_combined_financial_data[i].Company_name, list_combined_financial_data[i].EBITDA_Margin_TTM1,
			list_combined_financial_data[i].Company_name, list_combined_financial_data[i].TotalShareHoldersEquity_TTM1,
			list_combined_financial_data[i].Company_name, list_combined_financial_data[i].ShorttermDebt_TTM1,
			list_combined_financial_data[i].Company_name, list_combined_financial_data[i].LongtermDebt_TTM1,
			list_combined_financial_data[i].Company_name, list_combined_financial_data[i].SharesOutstanding1,
			list_combined_financial_data[i].Company_name, list_combined_financial_data[i].EPSDiluted1,
			list_combined_financial_data[i].Company_name, list_combined_financial_data[i].NetSales1,
			list_combined_financial_data[i].Company_name, list_combined_financial_data[i].Netprofit1,
			list_combined_financial_data[i].Company_name, list_combined_financial_data[i].AnnualDividend1,
			list_combined_financial_data[i].Company_name, list_combined_financial_data[i].COGS1,
			list_combined_financial_data[i].Company_name, list_combined_financial_data[i].PEGRatio_TTM1,
			list_combined_financial_data[i].Company_name, list_combined_financial_data[i].DividendPayout_TTM1,
		)
		financials_para = financials_para + data_string_stock
	}
	fmt.Println(financials_para)
	return list_combined_financial_data, list_user_holding_details, financials_para

}

func invoke_llm(query string, list_user_all_details []any) llm_response {
	url := "https://api.openai.com/v1/chat/completions"
	system_prompt := fmt.Sprintf("You are a helpful trading and investment assistant.This is all the details about the userts stock portfolio and the fundamental data for all of his stocks pls answer questions based on these details only. if the user asks for return for a perticular stock the formula will be (ltp(last traded price)-avg price(average price))*quantity,if the user asks for return percentage for a perticular stock the formula will be ((ltp(last traded price)-avg price(average price))/avg price(average price))*100, also show calculations where nessecary: %v", list_user_all_details)
	fmt.Println(system_prompt)
	requestBody := RequestBody{
		Model: "gpt-4o",
		Messages: []Message{
			{Role: "system", Content: system_prompt},
			{Role: "user", Content: query},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Println(err)
	}
	api_key := os.Getenv("apikey_openai")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+api_key)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(string(body))
	var llm_response llm_response
	err_llm := json.Unmarshal(body, &llm_response)
	if err != nil {
		log.Println(err_llm)
	}
	return llm_response
}
