package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
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
	Growth        growth_2
	Stability     stability_2
	Liquidity     liquidity_2
}
type growth struct {
	Success bool `json:"success"`
	Data    []struct {
		CoCode         float64 `json:"co_code"`
		Yrc            float64 `json:"yrc"`
		Netsalesgrowth float64 `json:"netsalesgrowth"`
		Ebitdagrowth   float64 `json:"ebitdagrowth"`
		Ebitgrowth     float64 `json:"ebitgrowth"`
		Patgrowth      float64 `json:"patgrowth"`
		Eps            float64 `json:"eps"`
	} `json:"data"`
	Message string `json:"message"`
}
type growth_2 struct {
	Success bool `json:"success"`
	Data    []struct {
		CoCode         Metric `json:"co_code"`
		Yrc            Metric `json:"yrc"`
		Netsalesgrowth Metric `json:"netsalesgrowth"`
		Ebitdagrowth   Metric `json:"ebitdagrowth"`
		Ebitgrowth     Metric `json:"ebitgrowth"`
		Patgrowth      Metric `json:"patgrowth"`
		Eps            Metric `json:"eps"`
	} `json:"data"`
	Message string `json:"message"`
	Total   int
}
type stability struct {
	Success bool `json:"success"`
	Data    []struct {
		CoCode          float64 `json:"co_code"`
		YRC             float64 `json:"YRC"`
		TotalDebtEquity float64 `json:"TotalDebt_Equity"`
		CurrentRatio    float64 `json:"CurrentRatio"`
		QuickRatio      float64 `json:"QuickRatio"`
		InterestCover   float64 `json:"InterestCover"`
		TotalDebtMCap   float64 `json:"TotalDebt_MCap"`
	} `json:"data"`
	Message string `json:"message"`
}
type stability_2 struct {
	Success bool `json:"success"`
	Data    []struct {
		CoCode          Metric `json:"co_code"`
		YRC             Metric `json:"YRC"`
		TotalDebtEquity Metric `json:"TotalDebt_Equity"`
		CurrentRatio    Metric `json:"CurrentRatio"`
		QuickRatio      Metric `json:"QuickRatio"`
		InterestCover   Metric `json:"InterestCover"`
		TotalDebtMCap   Metric `json:"TotalDebt_MCap"`
	} `json:"data"`
	Message string `json:"message"`
	Total   int
}
type liquidity struct {
	Success bool `json:"success"`
	Data    []struct {
		CoCode                           float64 `json:"co_code"`
		Yrc                              float64 `json:"yrc"`
		LoansToDeposits                  float64 `json:"loans_to_deposits"`
		CashToDeposits                   float64 `json:"cash_to_deposits"`
		InvestmentTodeposits             float64 `json:"investment_todeposits"`
		IncloanToDeposit                 float64 `json:"incloan_to_deposit"`
		CreditToDeposits                 float64 `json:"credit_to_deposits"`
		InterestexpendedToInterestearned float64 `json:"interestexpended_to_interestearned"`
		InterestincomeToTotalfunds       float64 `json:"interestincome_to_totalfunds"`
		InterestexpendedToTotalfunds     float64 `json:"interestexpended_to_totalfunds"`
		Casa                             float64 `json:"casa"`
	} `json:"data"`
	Message string `json:"message"`
}
type liquidity_2 struct {
	Success bool `json:"success"`
	Data    []struct {
		CoCode                           Metric `json:"co_code"`
		Yrc                              Metric `json:"yrc"`
		LoansToDeposits                  Metric `json:"loans_to_deposits"`
		CashToDeposits                   Metric `json:"cash_to_deposits"`
		InvestmentTodeposits             Metric `json:"investment_todeposits"`
		IncloanToDeposit                 Metric `json:"incloan_to_deposit"`
		CreditToDeposits                 Metric `json:"credit_to_deposits"`
		InterestexpendedToInterestearned Metric `json:"interestexpended_to_interestearned"`
		InterestincomeToTotalfunds       Metric `json:"interestincome_to_totalfunds"`
		InterestexpendedToTotalfunds     Metric `json:"interestexpended_to_totalfunds"`
		Casa                             Metric `json:"casa"`
	} `json:"data"`
	Message string `json:"message"`
	Total   int
}

var db *sql.DB
var dsn string
var api_key = os.Getenv("api_key_cmots")

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

	growth := growth_ratios(cocode)
	stability := stability_ratios(cocode)
	liquidity := liquidity_ratios(cocode)

	// Debug/log output (optional)
	fmt.Println("Valuation:", valuation)
	fmt.Println("Profitability:", profitabilty)
	fmt.Println("Debt:", debt)

	var combied_sentiment combined_sentiment_struct

	combied_sentiment.Valuation = valuation
	combied_sentiment.Debt = debt
	combied_sentiment.Profitability = profitabilty
	combied_sentiment.Growth = growth
	combied_sentiment.Stability = stability
	combied_sentiment.Liquidity = liquidity
	c.IndentedJSON(http.StatusOK, combied_sentiment)

}

func growth_ratios(cocode string) growth_2 {

	url := fmt.Sprintf("https://insbaapis.cmots.com/api/GrowthRatio/%s/S", cocode)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("Authorization", "Bearer "+api_key)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}
	var growth_struct growth
	fmt.Println(string(body))
	fmt.Println("done")
	err = json.Unmarshal(body, &growth_struct)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(growth_struct)

	var growth_struct_2 growth_2
	var growth_data struct {
		CoCode         Metric `json:"co_code"`
		Yrc            Metric `json:"yrc"`
		Netsalesgrowth Metric `json:"netsalesgrowth"`
		Ebitdagrowth   Metric `json:"ebitdagrowth"`
		Ebitgrowth     Metric `json:"ebitgrowth"`
		Patgrowth      Metric `json:"patgrowth"`
		Eps            Metric `json:"eps"`
	}
	growth_struct_2.Data = append(growth_struct_2.Data, growth_data)
	growth_struct_2.Data[0].CoCode.Value = growth_struct.Data[0].CoCode
	growth_struct_2.Data[0].Yrc.Value = growth_struct.Data[0].Yrc
	growth_struct_2.Data[0].Netsalesgrowth.Value = growth_struct.Data[0].Netsalesgrowth
	growth_struct_2.Data[0].Ebitdagrowth.Value = growth_struct.Data[0].Ebitdagrowth
	growth_struct_2.Data[0].Ebitgrowth.Value = growth_struct.Data[0].Ebitgrowth
	growth_struct_2.Data[0].Patgrowth.Value = growth_struct.Data[0].Patgrowth
	growth_struct_2.Data[0].Eps.Value = growth_struct.Data[0].Eps
	growth_score := 0

	// Evaluating Net Sales Growth
	// Typically market considers >15% excellent, 5-15% good, 0-5% moderate, <0% concerning
	if growth_struct_2.Data[0].Netsalesgrowth.Value < 0 {
		growth_struct_2.Data[0].Netsalesgrowth.Polarity = "bad"
		growth_score = growth_score - 1
	} else if growth_struct_2.Data[0].Netsalesgrowth.Value < 0.05 {
		growth_struct_2.Data[0].Netsalesgrowth.Polarity = "neutral"
		growth_score = growth_score + 0
	} else if growth_struct_2.Data[0].Netsalesgrowth.Value < 0.15 {
		growth_struct_2.Data[0].Netsalesgrowth.Polarity = "good"
		growth_score = growth_score + 1
	} else {
		growth_struct_2.Data[0].Netsalesgrowth.Polarity = "excellent"
		growth_score = growth_score + 2
	}

	// Evaluating EBITDA Growth
	// EBITDA growth expectations are similar to revenue but slightly higher
	if growth_struct_2.Data[0].Ebitdagrowth.Value < 0 {
		growth_struct_2.Data[0].Ebitdagrowth.Polarity = "bad"
		growth_score = growth_score - 1
	} else if growth_struct_2.Data[0].Ebitdagrowth.Value < 0.08 {
		growth_struct_2.Data[0].Ebitdagrowth.Polarity = "neutral"
		growth_score = growth_score + 0
	} else if growth_struct_2.Data[0].Ebitdagrowth.Value < 0.18 {
		growth_struct_2.Data[0].Ebitdagrowth.Polarity = "good"
		growth_score = growth_score + 1
	} else {
		growth_struct_2.Data[0].Ebitdagrowth.Polarity = "excellent"
		growth_score = growth_score + 2
	}

	// Evaluating EBIT Growth
	// Similar thresholds to EBITDA
	if growth_struct_2.Data[0].Ebitgrowth.Value < 0 {
		growth_struct_2.Data[0].Ebitgrowth.Polarity = "bad"
		growth_score = growth_score - 1
	} else if growth_struct_2.Data[0].Ebitgrowth.Value < 0.08 {
		growth_struct_2.Data[0].Ebitgrowth.Polarity = "neutral"
		growth_score = growth_score + 0
	} else if growth_struct_2.Data[0].Ebitgrowth.Value < 0.18 {
		growth_struct_2.Data[0].Ebitgrowth.Polarity = "good"
		growth_score = growth_score + 1
	} else {
		growth_struct_2.Data[0].Ebitgrowth.Polarity = "excellent"
		growth_score = growth_score + 2
	}

	// Evaluating PAT (Profit After Tax) Growth
	// Net income growth can be more volatile, so thresholds are wider
	if growth_struct_2.Data[0].Patgrowth.Value < 0 {
		growth_struct_2.Data[0].Patgrowth.Polarity = "bad"
		growth_score = growth_score - 1
	} else if growth_struct_2.Data[0].Patgrowth.Value < 0.10 {
		growth_struct_2.Data[0].Patgrowth.Polarity = "neutral"
		growth_score = growth_score + 0
	} else if growth_struct_2.Data[0].Patgrowth.Value < 0.20 {
		growth_struct_2.Data[0].Patgrowth.Polarity = "good"
		growth_score = growth_score + 1
	} else {
		growth_struct_2.Data[0].Patgrowth.Polarity = "excellent"
		growth_score = growth_score + 2
	}

	// Evaluating EPS Growth
	// EPS growth is closely watched and thresholds are similar to PAT
	if growth_struct_2.Data[0].Eps.Value < 0 {
		growth_struct_2.Data[0].Eps.Polarity = "bad"
		growth_score = growth_score - 1
	} else if growth_struct_2.Data[0].Eps.Value < 0.10 {
		growth_struct_2.Data[0].Eps.Polarity = "neutral"
		growth_score = growth_score + 0
	} else if growth_struct_2.Data[0].Eps.Value < 0.20 {
		growth_struct_2.Data[0].Eps.Polarity = "good"
		growth_score = growth_score + 1
	} else {
		growth_struct_2.Data[0].Eps.Polarity = "excellent"
		growth_score = growth_score + 2
	}
	growth_struct_2.Total = growth_score
	return growth_struct_2
}
func stability_ratios(cocode string) stability_2 {

	url := fmt.Sprintf("https://insbaapis.cmots.com/api/FinancialStabilityRatios/%s/S", cocode)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("Authorization", "Bearer "+api_key)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}
	var stability_struct stability
	fmt.Println(string(body))
	fmt.Println("done")
	err = json.Unmarshal(body, &stability_struct)
	if err != nil {
		log.Println(err)
	}

	var stability_struct_2 stability_2
	var stability_data struct {
		CoCode          Metric `json:"co_code"`
		YRC             Metric `json:"YRC"`
		TotalDebtEquity Metric `json:"TotalDebt_Equity"`
		CurrentRatio    Metric `json:"CurrentRatio"`
		QuickRatio      Metric `json:"QuickRatio"`
		InterestCover   Metric `json:"InterestCover"`
		TotalDebtMCap   Metric `json:"TotalDebt_MCap"`
	}
	stability_struct_2.Data = append(stability_struct_2.Data, stability_data)
	stability_struct_2.Data[0].CoCode.Value = stability_struct.Data[0].CoCode
	stability_struct_2.Data[0].YRC.Value = stability_struct.Data[0].YRC
	stability_struct_2.Data[0].TotalDebtEquity.Value = stability_struct.Data[0].TotalDebtEquity
	stability_struct_2.Data[0].CurrentRatio.Value = stability_struct.Data[0].CurrentRatio
	stability_struct_2.Data[0].QuickRatio.Value = stability_struct.Data[0].QuickRatio
	stability_struct_2.Data[0].InterestCover.Value = stability_struct.Data[0].InterestCover
	stability_struct_2.Data[0].TotalDebtMCap.Value = stability_struct.Data[0].TotalDebtMCap
	stability_score := 0

	// Evaluating Total Debt to Equity
	// Lower is better, industry standard thresholds
	if stability_struct_2.Data[0].TotalDebtEquity.Value > 2.0 {
		stability_struct_2.Data[0].TotalDebtEquity.Polarity = "bad"
		stability_score = stability_score - 1
	} else if stability_struct_2.Data[0].TotalDebtEquity.Value > 1.0 {
		stability_struct_2.Data[0].TotalDebtEquity.Polarity = "neutral"
		stability_score = stability_score + 0
	} else if stability_struct_2.Data[0].TotalDebtEquity.Value > 0.5 {
		stability_struct_2.Data[0].TotalDebtEquity.Polarity = "good"
		stability_score = stability_score + 1
	} else {
		stability_struct_2.Data[0].TotalDebtEquity.Polarity = "excellent"
		stability_score = stability_score + 2
	}

	// Evaluating Current Ratio
	// Typically, 1.5-3.0 is considered healthy
	if stability_struct_2.Data[0].CurrentRatio.Value < 1.0 {
		stability_struct_2.Data[0].CurrentRatio.Polarity = "bad"
		stability_score = stability_score - 1
	} else if stability_struct_2.Data[0].CurrentRatio.Value < 1.5 {
		stability_struct_2.Data[0].CurrentRatio.Polarity = "neutral"
		stability_score = stability_score + 0
	} else if stability_struct_2.Data[0].CurrentRatio.Value < 3.0 {
		stability_struct_2.Data[0].CurrentRatio.Polarity = "good"
		stability_score = stability_score + 1
	} else {
		// Too high might indicate inefficient use of assets
		stability_struct_2.Data[0].CurrentRatio.Polarity = "neutral"
		stability_score = stability_score + 0
	}

	// Evaluating Quick Ratio
	// Generally, a quick ratio of 1.0 or higher is good
	if stability_struct_2.Data[0].QuickRatio.Value < 0.7 {
		stability_struct_2.Data[0].QuickRatio.Polarity = "bad"
		stability_score = stability_score - 1
	} else if stability_struct_2.Data[0].QuickRatio.Value < 1.0 {
		stability_struct_2.Data[0].QuickRatio.Polarity = "neutral"
		stability_score = stability_score + 0
	} else if stability_struct_2.Data[0].QuickRatio.Value < 2.0 {
		stability_struct_2.Data[0].QuickRatio.Polarity = "good"
		stability_score = stability_score + 1
	} else {
		stability_struct_2.Data[0].QuickRatio.Polarity = "excellent"
		stability_score = stability_score + 2
	}

	// Evaluating Interest Coverage Ratio
	// Higher is better; below 1.5 is concerning
	if stability_struct_2.Data[0].InterestCover.Value < 1.5 {
		stability_struct_2.Data[0].InterestCover.Polarity = "bad"
		stability_score = stability_score - 1
	} else if stability_struct_2.Data[0].InterestCover.Value < 3.0 {
		stability_struct_2.Data[0].InterestCover.Polarity = "neutral"
		stability_score = stability_score + 0
	} else if stability_struct_2.Data[0].InterestCover.Value < 5.0 {
		stability_struct_2.Data[0].InterestCover.Polarity = "good"
		stability_score = stability_score + 1
	} else {
		stability_struct_2.Data[0].InterestCover.Polarity = "excellent"
		stability_score = stability_score + 2
	}

	// Evaluating Total Debt to Market Cap
	// Lower is generally better
	if stability_struct_2.Data[0].TotalDebtMCap.Value > 0.6 {
		stability_struct_2.Data[0].TotalDebtMCap.Polarity = "bad"
		stability_score = stability_score - 1
	} else if stability_struct_2.Data[0].TotalDebtMCap.Value > 0.3 {
		stability_struct_2.Data[0].TotalDebtMCap.Polarity = "neutral"
		stability_score = stability_score + 0
	} else if stability_struct_2.Data[0].TotalDebtMCap.Value > 0.1 {
		stability_struct_2.Data[0].TotalDebtMCap.Polarity = "good"
		stability_score = stability_score + 1
	} else {
		stability_struct_2.Data[0].TotalDebtMCap.Polarity = "excellent"
		stability_score = stability_score + 2
	}
	stability_struct_2.Total = stability_score
	return stability_struct_2
}

func liquidity_ratios(cocode string) liquidity_2 {

	url := fmt.Sprintf("https://insbaapis.cmots.com/api/LiquidityRatios/%s/S", cocode)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("Authorization", "Bearer "+api_key)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}
	var liquidity_struct liquidity
	fmt.Println(string(body))
	fmt.Println("done")
	err = json.Unmarshal(body, &liquidity_struct)
	if err != nil {
		log.Println(err)
	}

	var liquidity_struct_2 liquidity_2
	var liquidity_data struct {
		CoCode                           Metric `json:"co_code"`
		Yrc                              Metric `json:"yrc"`
		LoansToDeposits                  Metric `json:"loans_to_deposits"`
		CashToDeposits                   Metric `json:"cash_to_deposits"`
		InvestmentTodeposits             Metric `json:"investment_todeposits"`
		IncloanToDeposit                 Metric `json:"incloan_to_deposit"`
		CreditToDeposits                 Metric `json:"credit_to_deposits"`
		InterestexpendedToInterestearned Metric `json:"interestexpended_to_interestearned"`
		InterestincomeToTotalfunds       Metric `json:"interestincome_to_totalfunds"`
		InterestexpendedToTotalfunds     Metric `json:"interestexpended_to_totalfunds"`
		Casa                             Metric `json:"casa"`
	}
	liquidity_struct_2.Data = append(liquidity_struct_2.Data, liquidity_data)
	liquidity_struct_2.Data[0].CoCode.Value = liquidity_struct.Data[0].CoCode
	liquidity_struct_2.Data[0].Yrc.Value = liquidity_struct.Data[0].Yrc
	liquidity_struct_2.Data[0].LoansToDeposits.Value = liquidity_struct.Data[0].LoansToDeposits
	liquidity_struct_2.Data[0].CashToDeposits.Value = liquidity_struct.Data[0].CashToDeposits
	liquidity_struct_2.Data[0].InvestmentTodeposits.Value = liquidity_struct.Data[0].InvestmentTodeposits
	liquidity_struct_2.Data[0].IncloanToDeposit.Value = liquidity_struct.Data[0].IncloanToDeposit
	liquidity_struct_2.Data[0].CreditToDeposits.Value = liquidity_struct.Data[0].CreditToDeposits
	liquidity_struct_2.Data[0].InterestexpendedToInterestearned.Value = liquidity_struct.Data[0].InterestexpendedToInterestearned
	liquidity_struct_2.Data[0].InterestincomeToTotalfunds.Value = liquidity_struct.Data[0].InterestincomeToTotalfunds
	liquidity_struct_2.Data[0].InterestexpendedToTotalfunds.Value = liquidity_struct.Data[0].InterestexpendedToTotalfunds
	liquidity_struct_2.Data[0].Casa.Value = liquidity_struct.Data[0].Casa
	liquidity_score := 0

	// Evaluating Loans to Deposits (LTD)
	// Banking standard: 80-90% is considered optimal
	if liquidity_struct_2.Data[0].LoansToDeposits.Value > 1.0 {
		liquidity_struct_2.Data[0].LoansToDeposits.Polarity = "bad"
		liquidity_score = liquidity_score - 1
	} else if liquidity_struct_2.Data[0].LoansToDeposits.Value > 0.9 {
		liquidity_struct_2.Data[0].LoansToDeposits.Polarity = "neutral"
		liquidity_score = liquidity_score + 0
	} else if liquidity_struct_2.Data[0].LoansToDeposits.Value > 0.7 {
		liquidity_struct_2.Data[0].LoansToDeposits.Polarity = "good"
		liquidity_score = liquidity_score + 1
	} else {
		// Too low might indicate underutilization of deposits
		liquidity_struct_2.Data[0].LoansToDeposits.Polarity = "neutral"
		liquidity_score = liquidity_score + 0
	}

	// Evaluating Cash to Deposits
	// Higher is generally better for liquidity, but too high can mean inefficient use of funds
	if liquidity_struct_2.Data[0].CashToDeposits.Value < 0.05 {
		liquidity_struct_2.Data[0].CashToDeposits.Polarity = "bad"
		liquidity_score = liquidity_score - 1
	} else if liquidity_struct_2.Data[0].CashToDeposits.Value < 0.1 {
		liquidity_struct_2.Data[0].CashToDeposits.Polarity = "neutral"
		liquidity_score = liquidity_score + 0
	} else if liquidity_struct_2.Data[0].CashToDeposits.Value < 0.2 {
		liquidity_struct_2.Data[0].CashToDeposits.Polarity = "good"
		liquidity_score = liquidity_score + 1
	} else {
		// Very high cash reserves might be inefficient
		liquidity_struct_2.Data[0].CashToDeposits.Polarity = "neutral"
		liquidity_score = liquidity_score + 0
	}

	// Evaluating Investment to Deposits
	// Should be balanced - too low might mean missed opportunities, too high might indicate risk
	if liquidity_struct_2.Data[0].InvestmentTodeposits.Value < 0.1 {
		liquidity_struct_2.Data[0].InvestmentTodeposits.Polarity = "neutral"
		liquidity_score = liquidity_score + 0
	} else if liquidity_struct_2.Data[0].InvestmentTodeposits.Value < 0.25 {
		liquidity_struct_2.Data[0].InvestmentTodeposits.Polarity = "good"
		liquidity_score = liquidity_score + 1
	} else if liquidity_struct_2.Data[0].InvestmentTodeposits.Value < 0.4 {
		liquidity_struct_2.Data[0].InvestmentTodeposits.Polarity = "excellent"
		liquidity_score = liquidity_score + 2
	} else {
		liquidity_struct_2.Data[0].InvestmentTodeposits.Polarity = "neutral"
		liquidity_score = liquidity_score + 0
	}

	// Evaluating Interbank Loan to Deposit
	// Lower is generally better, indicates less dependence on interbank funding
	if liquidity_struct_2.Data[0].IncloanToDeposit.Value > 0.15 {
		liquidity_struct_2.Data[0].IncloanToDeposit.Polarity = "bad"
		liquidity_score = liquidity_score - 1
	} else if liquidity_struct_2.Data[0].IncloanToDeposit.Value > 0.1 {
		liquidity_struct_2.Data[0].IncloanToDeposit.Polarity = "neutral"
		liquidity_score = liquidity_score + 0
	} else {
		liquidity_struct_2.Data[0].IncloanToDeposit.Polarity = "good"
		liquidity_score = liquidity_score + 1
	}

	// Evaluating Credit to Deposits
	// Similar to LTD, but may include other credit instruments
	if liquidity_struct_2.Data[0].CreditToDeposits.Value > 0.9 {
		liquidity_struct_2.Data[0].CreditToDeposits.Polarity = "bad"
		liquidity_score = liquidity_score - 1
	} else if liquidity_struct_2.Data[0].CreditToDeposits.Value > 0.8 {
		liquidity_struct_2.Data[0].CreditToDeposits.Polarity = "neutral"
		liquidity_score = liquidity_score + 0
	} else {
		liquidity_struct_2.Data[0].CreditToDeposits.Polarity = "good"
		liquidity_score = liquidity_score + 1
	}

	// Evaluating Interest Expended to Interest Earned
	// Lower is better - efficiency metric
	if liquidity_struct_2.Data[0].InterestexpendedToInterestearned.Value > 0.7 {
		liquidity_struct_2.Data[0].InterestexpendedToInterestearned.Polarity = "bad"
		liquidity_score = liquidity_score - 1
	} else if liquidity_struct_2.Data[0].InterestexpendedToInterestearned.Value > 0.6 {
		liquidity_struct_2.Data[0].InterestexpendedToInterestearned.Polarity = "neutral"
		liquidity_score = liquidity_score + 0
	} else if liquidity_struct_2.Data[0].InterestexpendedToInterestearned.Value > 0.5 {
		liquidity_struct_2.Data[0].InterestexpendedToInterestearned.Polarity = "good"
		liquidity_score = liquidity_score + 1
	} else {
		liquidity_struct_2.Data[0].InterestexpendedToInterestearned.Polarity = "excellent"
		liquidity_score = liquidity_score + 2
	}

	// Evaluating Interest Income to Total Funds
	// Higher is generally better - shows earning capability
	if liquidity_struct_2.Data[0].InterestincomeToTotalfunds.Value < 0.06 {
		liquidity_struct_2.Data[0].InterestincomeToTotalfunds.Polarity = "bad"
		liquidity_score = liquidity_score - 1
	} else if liquidity_struct_2.Data[0].InterestincomeToTotalfunds.Value < 0.08 {
		liquidity_struct_2.Data[0].InterestincomeToTotalfunds.Polarity = "neutral"
		liquidity_score = liquidity_score + 0
	} else if liquidity_struct_2.Data[0].InterestincomeToTotalfunds.Value < 0.1 {
		liquidity_struct_2.Data[0].InterestincomeToTotalfunds.Polarity = "good"
		liquidity_score = liquidity_score + 1
	} else {
		liquidity_struct_2.Data[0].InterestincomeToTotalfunds.Polarity = "excellent"
		liquidity_score = liquidity_score + 2
	}

	// Evaluating Interest Expended to Total Funds
	// Lower is better - cost of funds
	if liquidity_struct_2.Data[0].InterestexpendedToTotalfunds.Value > 0.05 {
		liquidity_struct_2.Data[0].InterestexpendedToTotalfunds.Polarity = "bad"
		liquidity_score = liquidity_score - 1
	} else if liquidity_struct_2.Data[0].InterestexpendedToTotalfunds.Value > 0.04 {
		liquidity_struct_2.Data[0].InterestexpendedToTotalfunds.Polarity = "neutral"
		liquidity_score = liquidity_score + 0
	} else if liquidity_struct_2.Data[0].InterestexpendedToTotalfunds.Value > 0.03 {
		liquidity_struct_2.Data[0].InterestexpendedToTotalfunds.Polarity = "good"
		liquidity_score = liquidity_score + 1
	} else {
		liquidity_struct_2.Data[0].InterestexpendedToTotalfunds.Polarity = "excellent"
		liquidity_score = liquidity_score + 2
	}

	// Evaluating CASA (Current Account Savings Account) Ratio
	// Higher is better - indicates lower cost deposits
	if liquidity_struct_2.Data[0].Casa.Value < 0.3 {
		liquidity_struct_2.Data[0].Casa.Polarity = "bad"
		liquidity_score = liquidity_score - 1
	} else if liquidity_struct_2.Data[0].Casa.Value < 0.4 {
		liquidity_struct_2.Data[0].Casa.Polarity = "neutral"
		liquidity_score = liquidity_score + 0
	} else if liquidity_struct_2.Data[0].Casa.Value < 0.5 {
		liquidity_struct_2.Data[0].Casa.Polarity = "good"
		liquidity_score = liquidity_score + 1
	} else {
		liquidity_struct_2.Data[0].Casa.Polarity = "excellent"
		liquidity_score = liquidity_score + 2
	}
	liquidity_struct_2.Total = liquidity_score
	return liquidity_struct_2
}
