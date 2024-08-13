package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type FinancialData struct {
	Symbol                        string  `json:"symbol"`
	Date                          string  `json:"date"`
	CalendarYear                  string  `json:"calendarYear"`
	Period                        string  `json:"period"`
	CurrentRatio                  float64 `json:"currentRatio"`
	QuickRatio                    float64 `json:"quickRatio"`
	CashRatio                     float64 `json:"cashRatio"`
	DaysOfSalesOutstanding        float64 `json:"daysOfSalesOutstanding"`
	DaysOfInventoryOutstanding    float64 `json:"daysOfInventoryOutstanding"`
	OperatingCycle                float64 `json:"operatingCycle"`
	DaysOfPayablesOutstanding     float64 `json:"daysOfPayablesOutstanding"`
	CashConversionCycle           float64 `json:"cashConversionCycle"`
	GrossProfitMargin             float64 `json:"grossProfitMargin"`
	OperatingProfitMargin         float64 `json:"operatingProfitMargin"`
	PretaxProfitMargin            float64 `json:"pretaxProfitMargin"`
	NetProfitMargin               float64 `json:"netProfitMargin"`
	EffectiveTaxRate              float64 `json:"effectiveTaxRate"`
	ReturnOnAssets                float64 `json:"returnOnAssets"`
	ReturnOnEquity                float64 `json:"returnOnEquity"`
	ReturnOnCapitalEmployed       float64 `json:"returnOnCapitalEmployed"`
	NetIncomePerEBT               float64 `json:"netIncomePerEBT"`
	EbtPerEbit                    float64 `json:"ebtPerEbit"`
	EbitPerRevenue                float64 `json:"ebitPerRevenue"`
	DebtRatio                     float64 `json:"debtRatio"`
	DebtEquityRatio               float64 `json:"debtEquityRatio"`
	LongTermDebtToCapitalization  float64 `json:"longTermDebtToCapitalization"`
	TotalDebtToCapitalization     float64 `json:"totalDebtToCapitalization"`
	InterestCoverage              float64 `json:"interestCoverage"`
	CashFlowToDebtRatio           float64 `json:"cashFlowToDebtRatio"`
	CompanyEquityMultiplier       float64 `json:"companyEquityMultiplier"`
	ReceivablesTurnover           float64 `json:"receivablesTurnover"`
	PayablesTurnover              float64 `json:"payablesTurnover"`
	InventoryTurnover             float64 `json:"inventoryTurnover"`
	FixedAssetTurnover            float64 `json:"fixedAssetTurnover"`
	AssetTurnover                 float64 `json:"assetTurnover"`
	OperatingCashFlowPerShare     float64 `json:"operatingCashFlowPerShare"`
	FreeCashFlowPerShare          float64 `json:"freeCashFlowPerShare"`
	CashPerShare                  float64 `json:"cashPerShare"`
	PayoutRatio                   float64 `json:"payoutRatio"`
	OperatingCashFlowSalesRatio   float64 `json:"operatingCashFlowSalesRatio"`
	FreeCashFlowOperatingCashFlow float64 `json:"freeCashFlowOperatingCashFlowRatio"`
	CashFlowCoverageRatios        float64 `json:"cashFlowCoverageRatios"`
	ShortTermCoverageRatios       float64 `json:"shortTermCoverageRatios"`
	CapitalExpenditureCoverage    float64 `json:"capitalExpenditureCoverageRatio"`
	DividendPaidAndCapexCoverage  float64 `json:"dividendPaidAndCapexCoverageRatio"`
	DividendPayoutRatio           float64 `json:"dividendPayoutRatio"`
	PriceBookValueRatio           float64 `json:"priceBookValueRatio"`
	PriceToBookRatio              float64 `json:"priceToBookRatio"`
	PriceToSalesRatio             float64 `json:"priceToSalesRatio"`
	PriceEarningsRatio            float64 `json:"priceEarningsRatio"`
	PriceToFreeCashFlowsRatio     float64 `json:"priceToFreeCashFlowsRatio"`
	PriceToOperatingCashFlows     float64 `json:"priceToOperatingCashFlowsRatio"`
	PriceCashFlowRatio            float64 `json:"priceCashFlowRatio"`
	PriceEarningsToGrowthRatio    float64 `json:"priceEarningsToGrowthRatio"`
	PriceSalesRatio               float64 `json:"priceSalesRatio"`
	DividendYield                 float64 `json:"dividendYield"`
	EnterpriseValueMultiple       float64 `json:"enterpriseValueMultiple"`
	PriceFairValue                float64 `json:"priceFairValue"`
}
type GrowthData struct {
	Date                                   string  `json:"date"`
	Symbol                                 string  `json:"symbol"`
	CalendarYear                           string  `json:"calendarYear"`
	Period                                 string  `json:"period"`
	GrowthRevenue                          float64 `json:"growthRevenue"`
	GrowthCostOfRevenue                    float64 `json:"growthCostOfRevenue"`
	GrowthGrossProfit                      float64 `json:"growthGrossProfit"`
	GrowthGrossProfitRatio                 float64 `json:"growthGrossProfitRatio"`
	GrowthResearchAndDevelopmentExpenses   float64 `json:"growthResearchAndDevelopmentExpenses"`
	GrowthGeneralAndAdministrativeExpenses float64 `json:"growthGeneralAndAdministrativeExpenses"`
	GrowthSellingAndMarketingExpenses      float64 `json:"growthSellingAndMarketingExpenses"`
	GrowthOtherExpenses                    float64 `json:"growthOtherExpenses"`
	GrowthOperatingExpenses                float64 `json:"growthOperatingExpenses"`
	GrowthCostAndExpenses                  float64 `json:"growthCostAndExpenses"`
	GrowthInterestExpense                  float64 `json:"growthInterestExpense"`
	GrowthDepreciationAndAmortization      float64 `json:"growthDepreciationAndAmortization"`
	GrowthEBITDA                           float64 `json:"growthEBITDA"`
	GrowthEBITDARatio                      float64 `json:"growthEBITDARatio"`
	GrowthOperatingIncome                  float64 `json:"growthOperatingIncome"`
	GrowthOperatingIncomeRatio             float64 `json:"growthOperatingIncomeRatio"`
	GrowthTotalOtherIncomeExpensesNet      float64 `json:"growthTotalOtherIncomeExpensesNet"`
	GrowthIncomeBeforeTax                  float64 `json:"growthIncomeBeforeTax"`
	GrowthIncomeBeforeTaxRatio             float64 `json:"growthIncomeBeforeTaxRatio"`
	GrowthIncomeTaxExpense                 float64 `json:"growthIncomeTaxExpense"`
	GrowthNetIncome                        float64 `json:"growthNetIncome"`
	GrowthNetIncomeRatio                   float64 `json:"growthNetIncomeRatio"`
	GrowthEPS                              float64 `json:"growthEPS"`
	GrowthEPSDiluted                       float64 `json:"growthEPSDiluted"`
	GrowthWeightedAverageShsOut            float64 `json:"growthWeightedAverageShsOut"`
	GrowthWeightedAverageShsOutDil         float64 `json:"growthWeightedAverageShsOutDil"`
}
type dcf_advanced struct {
	Year                      string  `json:"year"`
	Symbol                    string  `json:"symbol"`
	Revenue                   float64 `json:"revenue"`
	RevenuePercentage         float64 `json:"revenuePercentage"`
	Ebitda                    float64 `json:"ebitda"`
	EbitdaPercentage          float64 `json:"ebitdaPercentage"`
	Ebit                      float64 `json:"ebit"`
	EbitPercentage            float64 `json:"ebitPercentage"`
	Depreciation              float64 `json:"depreciation"`
	DepreciationPercentage    float64 `json:"depreciationPercentage"`
	TotalCash                 float64 `json:"totalCash"`
	TotalCashPercentage       float64 `json:"totalCashPercentage"`
	Receivables               float64 `json:"receivables"`
	ReceivablesPercentage     float64 `json:"receivablesPercentage"`
	Inventories               float64 `json:"inventories"`
	InventoriesPercentage     float64 `json:"inventoriesPercentage"`
	Payable                   float64 `json:"payable"`
	PayablePercentage         float64 `json:"payablePercentage"`
	CapitalExpenditure        float64 `json:"capitalExpenditure"`
	CapitalExpenditurePercent float64 `json:"capitalExpenditurePercentage"`
	Price                     float64 `json:"price"`
	Beta                      float64 `json:"beta"`
	DilutedSharesOutstanding  float64 `json:"dilutedSharesOutstanding"`
	CostOfDebt                float64 `json:"costofDebt"`
	TaxRate                   float64 `json:"taxRate"`
	AfterTaxCostOfDebt        float64 `json:"afterTaxCostOfDebt"`
	RiskFreeRate              float64 `json:"riskFreeRate"`
	MarketRiskPremium         float64 `json:"marketRiskPremium"`
	CostOfEquity              float64 `json:"costOfEquity"`
	TotalDebt                 float64 `json:"totalDebt"`
	TotalEquity               float64 `json:"totalEquity"`
	TotalCapital              float64 `json:"totalCapital"`
	DebtWeighting             float64 `json:"debtWeighting"`
	EquityWeighting           float64 `json:"equityWeighting"`
	Wacc                      float64 `json:"wacc"`
	TaxRateCash               float64 `json:"taxRateCash"`
	Ebiat                     float64 `json:"ebiat"`
	Ufcf                      float64 `json:"ufcf"`
	SumPvUfcf                 float64 `json:"sumPvUfcf"`
	LongTermGrowthRate        float64 `json:"longTermGrowthRate"`
	TerminalValue             float64 `json:"terminalValue"`
	PresentTerminalValue      float64 `json:"presentTerminalValue"`
	EnterpriseValue           float64 `json:"enterpriseValue"`
	NetDebt                   float64 `json:"netDebt"`
	EquityValue               float64 `json:"equityValue"`
	EquityValuePerShare       float64 `json:"equityValuePerShare"`
	FreeCashFlowT1            float64 `json:"freeCashFlowT1"`
}
type LeveredDCF struct {
	Year                         string  `json:"year"`
	Symbol                       string  `json:"symbol"`
	Revenue                      float64 `json:"revenue"`
	RevenuePercentage            float64 `json:"revenuePercentage"`
	CapitalExpenditure           float64 `json:"capitalExpenditure"`
	CapitalExpenditurePercentage float64 `json:"capitalExpenditurePercentage"`
	Price                        float64 `json:"price"`
	Beta                         float64 `json:"beta"`
	DilutedSharesOutstanding     float64 `json:"dilutedSharesOutstanding"`
	CostOfDebt                   float64 `json:"costofDebt"`
	TaxRate                      float64 `json:"taxRate"`
	AfterTaxCostOfDebt           float64 `json:"afterTaxCostOfDebt"`
	RiskFreeRate                 float64 `json:"riskFreeRate"`
	MarketRiskPremium            float64 `json:"marketRiskPremium"`
	CostOfEquity                 float64 `json:"costOfEquity"`
	TotalDebt                    float64 `json:"totalDebt"`
	TotalEquity                  float64 `json:"totalEquity"`
	TotalCapital                 float64 `json:"totalCapital"`
	DebtWeighting                float64 `json:"debtWeighting"`
	EquityWeighting              float64 `json:"equityWeighting"`
	Wacc                         float64 `json:"wacc"`
	OperatingCashFlow            float64 `json:"operatingCashFlow"`
	PvLfcf                       float64 `json:"pvLfcf"`
	SumPvLfcf                    float64 `json:"sumPvLfcf"`
	LongTermGrowthRate           float64 `json:"longTermGrowthRate"`
	FreeCashFlow                 float64 `json:"freeCashFlow"`
	TerminalValue                float64 `json:"terminalValue"`
	PresentTerminalValue         float64 `json:"presentTerminalValue"`
	EnterpriseValue              float64 `json:"enterpriseValue"`
	NetDebt                      float64 `json:"netDebt"`
	EquityValue                  float64 `json:"equityValue"`
	EquityValuePerShare          float64 `json:"equityValuePerShare"`
	FreeCashFlowT1               float64 `json:"freeCashFlowT1"`
	OperatingCashFlowPercentage  float64 `json:"operatingCashFlowPercentage"`
}
type wrong_json struct {
	message string
}

var file, _ = os.Create("screener_data_uploader.log")

func main() {
	fmt.Println("screener_data_generator")
	router := gin.Default()
	router.GET("/get_ratios", get_data_financials)
	router.GET("/get_income_growth", get_income_growth)
	router.GET("/advanced_discounted_cash_flow", advanced_discounted_cash_flow)
	router.GET("/advanced_levered_discounted_cash_flow", advanced_levered_discounted_cash_flow)
	router.Run("localhost:8080")

}
func advanced_levered_discounted_cash_flow(c *gin.Context) {
	api_key := get_api_key_fundamentals(4)

	symbol, _ := c.GetQuery("symbol")

	client_id, _ := c.GetQuery("client_id")

	if client_id == "Karma100" {

		var json_data []LeveredDCF
		log.SetOutput(file)

		url := fmt.Sprintf("https://financialmodelingprep.com/api/v4/advanced_levered_discounted_cash_flow?symbol=%s&apikey=%s", symbol, api_key)

		res, err := http.Get(url)
		if err != nil {
			log.Println(err)
		}
		data, err := io.ReadAll(res.Body)
		if err != nil {
			log.Println(err)
		}

		err = json.Unmarshal(data, &json_data)
		if err != nil {
			log.Println(err)
		}

		c.IndentedJSON(http.StatusOK, json_data)
	} else {
		message := wrong_json{message: "wrong client id"}
		c.IndentedJSON(http.StatusBadRequest, message)
	}

}
func advanced_discounted_cash_flow(c *gin.Context) {
	api_key := get_api_key_fundamentals(4)

	symbol, _ := c.GetQuery("symbol")

	client_id, _ := c.GetQuery("client_id")

	if client_id == "Karma100" {

		var json_data []dcf_advanced
		log.SetOutput(file)

		url := fmt.Sprintf("https://financialmodelingprep.com/api/v4/advanced_discounted_cash_flow?symbol=%s&apikey=%s", symbol, api_key)

		res, err := http.Get(url)
		if err != nil {
			log.Println(err)
		}
		data, err := io.ReadAll(res.Body)
		if err != nil {
			log.Println(err)
		}

		err = json.Unmarshal(data, &json_data)
		if err != nil {
			log.Println(err)
		}

		c.IndentedJSON(http.StatusOK, json_data)
	} else {
		message := wrong_json{message: "wrong client id"}
		c.IndentedJSON(http.StatusBadRequest, message)
	}

}
func get_income_growth(c *gin.Context) {
	api_key := get_api_key_fundamentals(4)

	symbol, _ := c.GetQuery("symbol")

	period, _ := c.GetQuery("period")

	client_id, _ := c.GetQuery("client_id")

	if client_id == "Karma100" {

		var json_data []GrowthData
		log.SetOutput(file)
		url := fmt.Sprintf("https://financialmodelingprep.com/api/v3/income-statement-growth/%s?period=%s&apikey=%s&limit=%s", symbol, period, api_key, "10")

		res, err := http.Get(url)
		if err != nil {
			log.Println(err)
		}
		data, err := io.ReadAll(res.Body)
		if err != nil {
			log.Println(err)
		}

		err = json.Unmarshal(data, &json_data)
		if err != nil {
			log.Println(err)
		}

		c.IndentedJSON(http.StatusOK, json_data)
	} else {
		message := wrong_json{message: "wrong client id"}
		c.IndentedJSON(http.StatusBadRequest, message)
	}

}
func get_data_financials(c *gin.Context) {
	api_key := get_api_key_fundamentals(4)

	symbol, _ := c.GetQuery("symbol")

	period, _ := c.GetQuery("period")

	client_id, _ := c.GetQuery("client_id")

	if client_id == "Karma100" {

		var json_data []FinancialData
		log.SetOutput(file)
		url := fmt.Sprintf("https://financialmodelingprep.com/api/v3/ratios/%s?period=%s&apikey=%s&limit=%s", symbol, period, api_key, "10")

		res, err := http.Get(url)
		if err != nil {
			log.Println(err)
		}
		data, err := io.ReadAll(res.Body)
		if err != nil {
			log.Println(err)
		}

		err = json.Unmarshal(data, &json_data)
		if err != nil {
			log.Println(err)
		}

		c.IndentedJSON(http.StatusOK, json_data)
	} else {
		message := wrong_json{message: "wrong client id"}
		c.IndentedJSON(http.StatusBadRequest, message)
	}

}
func get_api_key_fundamentals(api_id int) string {
	db, err := sql.Open("mysql", "root:Karma100%@/alert_trade_db")
	log.SetOutput(file)
	if err != nil {
		log.Println(err)
	}

	var provider string
	var api_key string
	var secret_key string
	var n_start_time int64
	var n_last_update_time int64

	db.QueryRow("call stp_get_api_config(?)", api_id).Scan(&api_id, &provider, &api_key, &secret_key, &n_start_time, &n_last_update_time)
	return api_key
}
