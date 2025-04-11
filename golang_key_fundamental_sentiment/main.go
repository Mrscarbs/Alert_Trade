package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Metric struct {
	Value    float64 `json:"value"`
	Polarity string  `json:"polarity"`
}

type valuation_ratio struct {
	Pe                Metric `json:"pe"`
	Pb                Metric `json:"pb"`
	Price_to_cashflow Metric `json:"price_to_cashflow"`
	Peg_ratio         Metric `json:"peg_ratio"`
	Total_score       int
}

type company_profitablity struct {
	Roe                 Metric `json:"roe"`
	Roa                 Metric `json:"roa"`
	Roce                Metric `json:"roce"`
	Net_incom_margin    Metric `json:"net_income_margin"`
	Gross_income_margin Metric `json:"gross_income_margin"`
	Total_score         int
}

type company_debt struct {
	Debt_to_equity Metric `json:"debt_to_equity"`
	Short_termdebt Metric `json:"short_term_debt"`
	Long_termdebt  Metric `json:"long_term_debt"`
	Total_score    int
}

type FinancialAndCashFlowData struct {
	// From tbl_financial_data_cmots
	CO_CODE                     int
	TTMAson                     int
	MCAP                        float64
	EV                          float64
	PE                          float64
	PBV                         float64
	DIVYIELD                    float64
	EPS                         float64
	BookValue                   float64
	ROA_TTM                     float64
	ROE_TTM                     float64
	ROCE_TTM                    float64
	EBIT_TTM                    float64
	EBITDA_TTM                  float64
	EV_Sales_TTM                float64
	EV_EBITDA_TTM               float64
	NetIncomeMargin_TTM         float64
	GrossIncomeMargin_TTM       float64
	AssetTurnover_TTM           float64
	CurrentRatio_TTM            float64
	Debt_Equity_TTM             float64
	Sales_TotalAssets_TTM       float64
	NetDebt_EBITDA_TTM          float64
	EBITDA_Margin_TTM           float64
	TotalShareHoldersEquity_TTM float64
	ShorttermDebt_TTM           float64
	LongtermDebt_TTM            float64
	SharesOutstanding           float64
	EPSDiluted                  float64
	NetSales                    float64
	Netprofit                   float64
	AnnualDividend              float64
	COGS                        float64
	PEGRatio_TTM                float64
	DividendPayout_TTM          float64

	// From tbl_cash_flow_data_cmots
	CO_CODE_CashFlow     int
	YRC                  int
	CashFlowPerShare     float64
	PricetoCashFlowRatio float64
	FreeCashFlowperShare float64
	PricetoFreeCashFlow  float64
	FreeCashFlowYield    float64
	Salestocashflowratio float64
}

type combined_sentiment_struct struct {
	Valuation     valuation_ratio
	Profitability company_profitablity
	Debt          company_debt
}

var db *sql.DB
var dsn string

func main() {

	var err error
	file, errfile := os.Create("golang_key_fundamental_sentiment_logs.log")
	if errfile != nil {
		fmt.Println(errfile)
	}
	dsn = os.Getenv("dsn")
	log.SetOutput(file)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Println(err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(time.Minute * 5)

	router := gin.Default()
	router.GET("/get_fundamental_sentiment", get_stockwise_details)
	router.Run("0.0.0.0:8080")
}

func get_stockwise_details(c *gin.Context) {
	cocode, _ := c.GetQuery("cocode")
	var financisdata FinancialAndCashFlowData
	err := db.QueryRow("CALL stp_get_full_data_by_cocode(?)", cocode).Scan(
		// ... all your existing Scan fields
		&financisdata.CO_CODE,
		&financisdata.TTMAson,
		&financisdata.MCAP,
		&financisdata.EV,
		&financisdata.PE,
		&financisdata.PBV,
		&financisdata.DIVYIELD,
		&financisdata.EPS,
		&financisdata.BookValue,
		&financisdata.ROA_TTM,
		&financisdata.ROE_TTM,
		&financisdata.ROCE_TTM,
		&financisdata.EBIT_TTM,
		&financisdata.EBITDA_TTM,
		&financisdata.EV_Sales_TTM,
		&financisdata.EV_EBITDA_TTM,
		&financisdata.NetIncomeMargin_TTM,
		&financisdata.GrossIncomeMargin_TTM,
		&financisdata.AssetTurnover_TTM,
		&financisdata.CurrentRatio_TTM,
		&financisdata.Debt_Equity_TTM,
		&financisdata.Sales_TotalAssets_TTM,
		&financisdata.NetDebt_EBITDA_TTM,
		&financisdata.EBITDA_Margin_TTM,
		&financisdata.TotalShareHoldersEquity_TTM,
		&financisdata.ShorttermDebt_TTM,
		&financisdata.LongtermDebt_TTM,
		&financisdata.SharesOutstanding,
		&financisdata.EPSDiluted,
		&financisdata.NetSales,
		&financisdata.Netprofit,
		&financisdata.AnnualDividend,
		&financisdata.COGS,
		&financisdata.PEGRatio_TTM,
		&financisdata.DividendPayout_TTM,
		&financisdata.CO_CODE_CashFlow,
		&financisdata.YRC,
		&financisdata.CashFlowPerShare,
		&financisdata.PricetoCashFlowRatio,
		&financisdata.FreeCashFlowperShare,
		&financisdata.PricetoFreeCashFlow,
		&financisdata.FreeCashFlowYield,
		&financisdata.Salestocashflowratio,
	)
	if err != nil {
		log.Println("Scan error:", err)
		return
	}

	var valuation valuation_ratio
	var profitabilty company_profitablity
	var debt company_debt

	// Valuation
	valuation_score := 0
	valuation.Pb.Value = financisdata.PBV
	if valuation.Pb.Value > 3 {
		valuation.Pb.Polarity = "bad"
		valuation_score = valuation_score - 1
	} else if valuation.Pb.Value > 1 {
		valuation.Pb.Polarity = "neutral"
		valuation_score = valuation_score + 0
	} else {
		valuation.Pb.Polarity = "good"
		valuation_score = valuation_score + 1
	}

	valuation.Pe.Value = financisdata.PE
	if valuation.Pe.Value > 25 {
		valuation.Pe.Polarity = "bad"
		valuation_score = valuation_score - 1
	} else if valuation.Pe.Value > 15 {
		valuation.Pe.Polarity = "neutral"
		valuation_score = valuation_score + 0
	} else {
		valuation.Pe.Polarity = "good"
		valuation_score = valuation_score + 1
	}

	valuation.Peg_ratio.Value = financisdata.PEGRatio_TTM
	if valuation.Peg_ratio.Value > 2 {
		valuation.Peg_ratio.Polarity = "bad"
		valuation_score = valuation_score - 1
	} else if valuation.Peg_ratio.Value > 1 {
		valuation.Peg_ratio.Polarity = "neutral"
		valuation_score = valuation_score + 0
	} else {
		valuation.Peg_ratio.Polarity = "good"
		valuation_score = valuation_score + 1
	}

	valuation.Price_to_cashflow.Value = financisdata.PricetoCashFlowRatio
	if valuation.Price_to_cashflow.Value > 20 {
		valuation.Price_to_cashflow.Polarity = "bad"
		valuation_score = valuation_score - 1
	} else if valuation.Price_to_cashflow.Value > 10 {
		valuation.Price_to_cashflow.Polarity = "neutral"
		valuation_score = valuation_score + 0
	} else {
		valuation.Price_to_cashflow.Polarity = "good"
		valuation_score = valuation_score + 1
	}
	valuation.Total_score = valuation_score

	// Profitability
	profitability_score := 0

	profitabilty.Roe.Value = financisdata.ROE_TTM
	if profitabilty.Roe.Value > 20 {
		profitabilty.Roe.Polarity = "good"
		profitability_score = profitability_score + 1
	} else if profitabilty.Roe.Value > 10 {
		profitabilty.Roe.Polarity = "neutral"
		profitability_score = profitability_score + 0
	} else {
		profitabilty.Roe.Polarity = "bad"
		profitability_score = profitability_score - 1
	}

	profitabilty.Roa.Value = financisdata.ROA_TTM
	if profitabilty.Roa.Value > 7 {
		profitabilty.Roa.Polarity = "good"
		profitability_score = profitability_score + 1
	} else if profitabilty.Roa.Value > 3 {
		profitabilty.Roa.Polarity = "neutral"
		profitability_score = profitability_score + 0
	} else {
		profitabilty.Roa.Polarity = "bad"
		profitability_score = profitability_score - 1
	}

	profitabilty.Roce.Value = financisdata.ROCE_TTM
	if profitabilty.Roce.Value > 20 {
		profitabilty.Roce.Polarity = "good"
		profitability_score = profitability_score + 1
	} else if profitabilty.Roce.Value > 10 {
		profitabilty.Roce.Polarity = "neutral"
		profitability_score = profitability_score + 0
	} else {
		profitabilty.Roce.Polarity = "bad"
		profitability_score = profitability_score - 1
	}

	profitabilty.Net_incom_margin.Value = financisdata.NetIncomeMargin_TTM
	if profitabilty.Net_incom_margin.Value > 20 {
		profitabilty.Net_incom_margin.Polarity = "good"
		profitability_score = profitability_score + 1
	} else if profitabilty.Net_incom_margin.Value > 10 {
		profitabilty.Net_incom_margin.Polarity = "neutral"
		profitability_score = profitability_score + 0
	} else {
		profitabilty.Net_incom_margin.Polarity = "bad"
		profitability_score = profitability_score - 1
	}

	profitabilty.Gross_income_margin.Value = financisdata.GrossIncomeMargin_TTM
	if profitabilty.Gross_income_margin.Value > 50 {
		profitabilty.Gross_income_margin.Polarity = "good"
		profitability_score = profitability_score + 1
	} else if profitabilty.Gross_income_margin.Value > 30 {
		profitabilty.Gross_income_margin.Polarity = "neutral"
		profitability_score = profitability_score + 0
	} else {
		profitabilty.Gross_income_margin.Polarity = "bad"
		profitability_score = profitability_score - 1
	}

	profitabilty.Total_score = profitability_score

	// Debt
	debt_score := 0

	debt.Debt_to_equity.Value = financisdata.Debt_Equity_TTM
	if debt.Debt_to_equity.Value > 2 {
		debt.Debt_to_equity.Polarity = "bad"
		debt_score = debt_score - 1
	} else if debt.Debt_to_equity.Value > 1 {
		debt.Debt_to_equity.Polarity = "neutral"
		debt_score = debt_score + 0
	} else {
		debt.Debt_to_equity.Polarity = "good"
		debt_score = debt_score + 1
	}

	debt.Short_termdebt.Value = financisdata.ShorttermDebt_TTM
	if debt.Short_termdebt.Value > 0 {
		debt.Short_termdebt.Polarity = "bad"
		debt_score = debt_score - 1
	} else {
		debt.Short_termdebt.Polarity = "good"
		debt_score = debt_score + 1
	}

	debt.Long_termdebt.Value = financisdata.LongtermDebt_TTM
	if debt.Long_termdebt.Value > 0 {
		debt.Long_termdebt.Polarity = "bad"
		debt_score = debt_score - 1
	} else {
		debt.Long_termdebt.Polarity = "good"
		debt_score = debt_score + 1
	}

	debt.Total_score = debt_score

	// Debug/log output (optional)
	fmt.Println("Valuation:", valuation)
	fmt.Println("Profitability:", profitabilty)
	fmt.Println("Debt:", debt)

	var combied_sentiment combined_sentiment_struct

	combied_sentiment.Valuation = valuation
	combied_sentiment.Debt = debt
	combied_sentiment.Profitability = profitabilty
	c.IndentedJSON(http.StatusOK, combied_sentiment)

}
