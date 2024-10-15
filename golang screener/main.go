package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var log_file, _ = os.Create("screener_log.log")

// Combined struct for `tbl_financial_data_cmots` and `tbl_cash_flow_data_cmots` with duplicate `co_code` handled
// Combined struct for `tbl_financial_data_cmots` and `tbl_cash_flow_data_cmots` without pointer types
type FinancialDataCombined struct {
	// Columns from `tbl_financial_data_cmots`
	CoCode                     int     `json:"co_code"`
	TTMAson                    int     `json:"ttm_ason"`
	Mcap                       float64 `json:"mcap"`
	Ev                         float64 `json:"ev"`
	Pe                         float64 `json:"pe"`
	Pbv                        float64 `json:"pbv"`
	DivYield                   float64 `json:"div_yield"`
	Eps                        float64 `json:"eps"`
	BookValue                  float64 `json:"book_value"`
	RoaTTM                     float64 `json:"roa_ttm"`
	RoeTTM                     float64 `json:"roe_ttm"`
	RoceTTM                    float64 `json:"roce_ttm"`
	EbitTTM                    float64 `json:"ebit_ttm"`
	EbitdaTTM                  float64 `json:"ebitda_ttm"`
	EVSalesTTM                 float64 `json:"ev_sales_ttm"`
	EVEBITDATTM                float64 `json:"ev_ebitda_ttm"`
	NetIncomeMarginTTM         float64 `json:"net_income_margin_ttm"`
	GrossIncomeMarginTTM       float64 `json:"gross_income_margin_ttm"`
	AssetTurnoverTTM           float64 `json:"asset_turnover_ttm"`
	CurrentRatioTTM            float64 `json:"current_ratio_ttm"`
	DebtEquityTTM              float64 `json:"debt_equity_ttm"`
	SalesTotalAssetsTTM        float64 `json:"sales_total_assets_ttm"`
	NetDebtEBITDATTM           float64 `json:"net_debt_ebitda_ttm"`
	EBITDA_MarginTTM           float64 `json:"ebitda_margin_ttm"`
	TotalShareholdersEquityTTM float64 `json:"total_shareholders_equity_ttm"`
	ShorttermDebtTTM           float64 `json:"short_term_debt_ttm"`
	LongtermDebtTTM            float64 `json:"long_term_debt_ttm"`
	SharesOutstanding          float64 `json:"shares_outstanding"`
	EPSDiluted                 float64 `json:"eps_diluted"`
	NetSales                   float64 `json:"net_sales"`
	NetProfit                  float64 `json:"net_profit"`
	AnnualDividend             float64 `json:"annual_dividend"`
	Cogs                       float64 `json:"cogs"`
	PegRatioTTM                float64 `json:"peg_ratio_ttm"`
	DividendPayoutTTM          float64 `json:"dividend_payout_ttm"`

	// Columns from `tbl_cash_flow_data_cmots`
	CoCode2              int     `json:"co_code2"` // Renamed co_code from tbl_cash_flow_data_cmots
	Yrc                  int     `json:"yrc"`
	CashFlowPerShare     float64 `json:"cash_flow_per_share"`
	PriceToCashFlowRatio float64 `json:"price_to_cash_flow_ratio"`
	FreeCashFlowPerShare float64 `json:"free_cash_flow_per_share"`
	PriceToFreeCashFlow  float64 `json:"price_to_free_cash_flow"`
	FreeCashFlowYield    float64 `json:"free_cash_flow_yield"`
	SalesToCashFlowRatio float64 `json:"sales_to_cash_flow_ratio"`
}

func main() {

	router := gin.Default()
	router.GET("/screener", screener)
	router.Run("localhost:8088")
}

func screener(c *gin.Context) {
	query, _ := c.GetQuery("qu")
	fmt.Println(query)
	db, err_db_open := sql.Open("mysql", "root:Karma100%@tcp(alerttrade.cbgqgqswkxrn.eu-north-1.rds.amazonaws.com:3306)/alert_trade_db")
	log.SetOutput(log_file)
	if err_db_open != nil {
		log.Println(err_db_open)
	}
	defer db.Close()
	rows, err := db.Query("call stp_GetFinancialDataByCondition_cmots_2(?)", query)
	if err != nil {
		log.Println(err)
	}
	list_struct_screen := []FinancialDataCombined{}
	for rows.Next() {
		var financial_datacom FinancialDataCombined
		err := rows.Scan(
			&financial_datacom.CoCode,
			&financial_datacom.TTMAson,
			&financial_datacom.Mcap,
			&financial_datacom.Ev,
			&financial_datacom.Pe,
			&financial_datacom.Pbv,
			&financial_datacom.DivYield,
			&financial_datacom.Eps,
			&financial_datacom.BookValue,
			&financial_datacom.RoaTTM,
			&financial_datacom.RoeTTM,
			&financial_datacom.RoceTTM,
			&financial_datacom.EbitTTM,
			&financial_datacom.EbitdaTTM,
			&financial_datacom.EVSalesTTM,
			&financial_datacom.EVEBITDATTM,
			&financial_datacom.NetIncomeMarginTTM,
			&financial_datacom.GrossIncomeMarginTTM,
			&financial_datacom.AssetTurnoverTTM,
			&financial_datacom.CurrentRatioTTM,
			&financial_datacom.DebtEquityTTM,
			&financial_datacom.SalesTotalAssetsTTM,
			&financial_datacom.NetDebtEBITDATTM,
			&financial_datacom.EBITDA_MarginTTM,
			&financial_datacom.TotalShareholdersEquityTTM,
			&financial_datacom.ShorttermDebtTTM,
			&financial_datacom.LongtermDebtTTM,
			&financial_datacom.SharesOutstanding,
			&financial_datacom.EPSDiluted,
			&financial_datacom.NetSales,
			&financial_datacom.NetProfit,
			&financial_datacom.AnnualDividend,
			&financial_datacom.Cogs,
			&financial_datacom.PegRatioTTM,
			&financial_datacom.DividendPayoutTTM,
			&financial_datacom.CoCode2,
			&financial_datacom.Yrc,
			&financial_datacom.CashFlowPerShare,
			&financial_datacom.PriceToCashFlowRatio,
			&financial_datacom.FreeCashFlowPerShare,
			&financial_datacom.PriceToFreeCashFlow,
			&financial_datacom.FreeCashFlowYield,
			&financial_datacom.SalesToCashFlowRatio,
		)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Error scanning row"})
			return
		}
		list_struct_screen = append(list_struct_screen, financial_datacom)
	}
	c.IndentedJSON(http.StatusOK, list_struct_screen)

}
