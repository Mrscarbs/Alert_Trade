package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type CombinedFinancialData struct {
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

var log_file, _ = os.Create("golang_key_metrics_log.log")

func main() {
	router := gin.Default()
	router.GET("/get_key_metrics", get_key_metrics)
	router.Run("0.0.0.0:8080")

}

func get_key_metrics(c *gin.Context) {
	var data CombinedFinancialData
	co_code, _ := c.GetQuery("co_code")
	log.SetOutput(log_file)
	db, err := sql.Open("mysql", "root:Karma100%@tcp(alerttrade.cbgqgqswkxrn.eu-north-1.rds.amazonaws.com:3306)/alert_trade_db")
	if err != nil {
		log.Println(err)
	}
	err_db := db.QueryRow("call stp_get_combined_financial_data_by_cocode(?)", co_code).Scan(
		&data.CO_CODE, &data.TTMAson,
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

	if err_db != nil {
		log.Println(err)
	}

	c.IndentedJSON(http.StatusOK, data)

}
