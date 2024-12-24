package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const api_key = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1bmlxdWVfbmFtZSI6Imluc2JhYXBpcyIsInJvbGUiOiJBZG1pbiIsIm5iZiI6MTczNDk2MDMwMiwiZXhwIjoxNzY2NjY5MTAyLCJpYXQiOjE3MzQ5NjAzMDIsImlzcyI6Imh0dHA6Ly9sb2NhbGhvc3Q6NTAxOTEiLCJhdWQiOiJodHRwOi8vbG9jYWxob3N0OjUwMTkxIn0.UqzyAKBMcDMmPL-kgaZtnusOAWOuB3v1tVIu_PZsJp8"

type struct_company_master struct {
	Success bool `json:"success"`
	Data    []struct {
		CoCode           float64 `json:"co_code"`
		Bsecode          string  `json:"bsecode"`
		Nsesymbol        string  `json:"nsesymbol"`
		Companyname      string  `json:"companyname"`
		Companyshortname string  `json:"companyshortname"`
		Categoryname     string  `json:"categoryname"`
		Isin             string  `json:"isin"`
		Bsegroup         string  `json:"bsegroup"`
		Mcaptype         string  `json:"mcaptype"`
		Sectorcode       string  `json:"sectorcode"`
		Sectorname       string  `json:"sectorname"`
		Industrycode     string  `json:"industrycode"`
		Industryname     string  `json:"industryname"`
		Bselistedflag    string  `json:"bselistedflag"`
		Nselistedflag    string  `json:"nselistedflag"`
		Displaytype      string  `json:"displaytype"`
		BSEStatus        string  `json:"BSEStatus"`
		NSEStatus        string  `json:"NSEStatus"`
	} `json:"data"`
	Message string `json:"message"`
}
type struct_ttm struct {
	Success bool `json:"success"`
	Data    []struct {
		CoCode                     int     `json:"co_code"`
		TTMAson                    int     `json:"TTMAson"`
		MCAP                       float64 `json:"MCAP"`
		EV                         float64 `json:"EV"`
		PE                         float64 `json:"PE"`
		PBV                        float64 `json:"PBV"`
		DIVYIELD                   float64 `json:"DIVYIELD"`
		EPS                        float64 `json:"EPS"`
		BookValue                  float64 `json:"BookValue"`
		ROATTM                     float64 `json:"ROA_TTM"`
		ROETTM                     float64 `json:"ROE_TTM"`
		ROCETTM                    float64 `json:"ROCE_TTM"`
		EBITTTM                    float64 `json:"EBIT_TTM"`
		EBITDATTM                  float64 `json:"EBITDA_TTM"`
		EVSalesTTM                 float64 `json:"EV_Sales_TTM"`
		EVEBITDATTM                float64 `json:"EV_EBITDA_TTM"`
		NetIncomeMarginTTM         float64 `json:"NetIncomeMargin_TTM"`
		GrossIncomeMarginTTM       float64 `json:"GrossIncomeMargin_TTM"`
		AssetTurnoverTTM           float64 `json:"AssetTurnover_TTM"`
		CurrentRatioTTM            float64 `json:"CurrentRatio_TTM"`
		DebtEquityTTM              float64 `json:"Debt_Equity_TTM"`
		SalesTotalAssetsTTM        float64 `json:"Sales_TotalAssets_TTM"`
		NetDebtEBITDATTM           float64 `json:"NetDebt_EBITDA_TTM"`
		EBITDAMarginTTM            float64 `json:"EBITDA_Margin_TTM"`
		TotalShareHoldersEquityTTM float64 `json:"TotalShareHoldersEquity_TTM"`
		ShorttermDebtTTM           float64 `json:"ShorttermDebt_TTM"`
		LongtermDebtTTM            float64 `json:"LongtermDebt_TTM"`
		SharesOutstanding          int     `json:"SharesOutstanding"`
		EPSDiluted                 float64 `json:"EPSDiluted"`
		NetSales                   float64 `json:"NetSales"`
		Netprofit                  float64 `json:"Netprofit"`
		AnnualDividend             float64 `json:"AnnualDividend"`
		COGS                       float64 `json:"COGS"`
		PEGRatioTTM                float64 `json:"PEGRatio_TTM"`
		DividendPayoutTTM          float64 `json:"DividendPayout_TTM"`
	} `json:"data"`
	Message string `json:"message"`
}
type stuct_quaterly_ratios struct {
	Success bool `json:"success"`
	Data    []struct {
		CoCode            int     `json:"co_code"`
		Qtrend            int     `json:"qtrend"`
		Mcap              float64 `json:"mcap"`
		Ev                float64 `json:"ev"`
		Pe                float64 `json:"pe"`
		Pbv               float64 `json:"pbv"`
		Eps               float64 `json:"eps"`
		Bookvalue         float64 `json:"bookvalue"`
		Ebit              float64 `json:"ebit"`
		Ebitda            float64 `json:"ebitda"`
		EvSales           float64 `json:"ev_sales"`
		EvEbitda          float64 `json:"ev_ebitda"`
		Netincomemargin   float64 `json:"netincomemargin"`
		Grossincomemargin float64 `json:"grossincomemargin"`
		Ebitdamargin      float64 `json:"ebitdamargin"`
		Epsdiluted        float64 `json:"epsdiluted"`
		Netsales          float64 `json:"netsales"`
		Netprofit         float64 `json:"netprofit"`
		Cogs              float64 `json:"cogs"`
	} `json:"data"`
	Message string `json:"message"`
}
type cashflow_ratio struct {
	Success bool `json:"success"`
	Data    []struct {
		CoCode               float64 `json:"co_code"`
		Yrc                  float64 `json:"yrc"`
		Cashflowpershare     float64 `json:"cashflowpershare"`
		Pricetocashflowratio float64 `json:"pricetocashflowratio"`
		Freecashflowpershare float64 `json:"freecashflowpershare"`
		Pricetofreecashflow  float64 `json:"pricetofreecashflow"`
		Freecashflowyield    float64 `json:"freecashflowyield"`
		Salestocashflowratio float64 `json:"salestocashflowratio"`
	} `json:"data"`
	Message string `json:"message"`
}

func main() {
	log_file, err := os.Create("company_master.log")
	if err != nil {
		fmt.Println(err)
	}
	var i = 0
	for i == 0 {
		company_master(log_file)
		get_ttm_ratios(log_file)
		get_quaterly_ratios(log_file)
		get_quaterly_ratios2(log_file)
		get_quaterly_ratios3(log_file)
		get_cash_fow_ratio(log_file)
		time.Sleep(time.Hour * 24)

	}

}

func company_master(log_file *os.File) {
	var unfolded_json struct_company_master
	log.SetOutput(log_file)
	url := "https://insbaapis.cmots.com/api/CompanyMaster"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("Authorization", "Bearer "+api_key)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Println(err)
	}

	res_body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}

	// fmt.Println(string(res_body))

	json.Unmarshal(res_body, &unfolded_json)
	// fmt.Println(unfolded_json)
	db, err_db_open := sql.Open("mysql", "admin:saumitrasuparn@tcp(alerttradedb.czqug0e2in8p.ap-south-1.rds.amazonaws.com:3306)/alert_trade_db")

	if err_db_open != nil {
		log.Println(err_db_open)
	}
	for i := 0; i < len(unfolded_json.Data); i++ {
		data_struct := unfolded_json.Data[i]
		_, err := db.Exec(
			"CALL stp_insert_or_update_company_master_cmots(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
			data_struct.CoCode, data_struct.Bsecode, data_struct.Nsesymbol,
			data_struct.Companyname, data_struct.Companyshortname, data_struct.Categoryname,
			data_struct.Isin, data_struct.Bsegroup, data_struct.Mcaptype,
			data_struct.Sectorcode, data_struct.Sectorname, data_struct.Industrycode,
			data_struct.Industryname, data_struct.Bselistedflag, data_struct.Nselistedflag,
			data_struct.Displaytype, data_struct.BSEStatus, data_struct.NSEStatus,
		)
		if err != nil {
			log.Println("Error executing stored procedure:", err)
		}
	}
	db.Close()
}

func get_ttm_ratios(log_file *os.File) {
	var unfolded_ttm struct_ttm
	db, err_db_open := sql.Open("mysql", "admin:saumitrasuparn@tcp(alerttradedb.czqug0e2in8p.ap-south-1.rds.amazonaws.com:3306)/alert_trade_db")
	log.SetOutput(log_file)
	if err_db_open != nil {
		log.Println(err_db_open)
	}
	rows, err := db.Query("call stp_get_all_cocodes_cmots()")
	if err != nil {
		log.Println(err)
	}
	list_cocodes := []int{}
	for rows.Next() {
		var cocode int
		rows.Scan(&cocode)
		list_cocodes = append(list_cocodes, cocode)
	}

	// fmt.Println(list_cocodes)
	for i := 0; i < len(list_cocodes); i++ {
		str_cocode := strconv.Itoa(list_cocodes[i])
		url := fmt.Sprintf("https://insbaapis.cmots.com/api/DailyRatios/%s/S/", str_cocode)
		fmt.Println(url)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Println(err)
		}
		req.Header.Set("Authorization", "Bearer "+api_key)

		res, err := http.DefaultClient.Do(req)

		if err != nil {
			log.Println(err)
		}

		res_body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Println(err)
		}

		json.Unmarshal(res_body, &unfolded_ttm)
		data_struct := unfolded_ttm.Data[0]
		_, err_exec := db.Exec(
			"CALL stp_insert_or_update_financial_data_cmots(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
			data_struct.CoCode, data_struct.TTMAson, data_struct.MCAP, data_struct.EV, data_struct.PE,
			data_struct.PBV, data_struct.DIVYIELD, data_struct.EPS, data_struct.BookValue, data_struct.ROATTM,
			data_struct.ROETTM, data_struct.ROCETTM, data_struct.EBITTTM, data_struct.EBITDATTM, data_struct.EVSalesTTM,
			data_struct.EVEBITDATTM, data_struct.NetIncomeMarginTTM, data_struct.GrossIncomeMarginTTM,
			data_struct.AssetTurnoverTTM, data_struct.CurrentRatioTTM, data_struct.DebtEquityTTM, data_struct.SalesTotalAssetsTTM,
			data_struct.NetDebtEBITDATTM, data_struct.EBITDAMarginTTM, data_struct.TotalShareHoldersEquityTTM,
			data_struct.ShorttermDebtTTM, data_struct.LongtermDebtTTM, data_struct.SharesOutstanding,
			data_struct.EPSDiluted, data_struct.NetSales, data_struct.Netprofit, data_struct.AnnualDividend,
			data_struct.COGS, data_struct.PEGRatioTTM, data_struct.DividendPayoutTTM,
		)
		if err_exec != nil {
			log.Println("Error executing stored procedure:", err_exec)
		}
	}
	db.Close()

}

func get_quaterly_ratios(log_file *os.File) {
	var unfolded_quaterly stuct_quaterly_ratios
	db, err_db_open := sql.Open("mysql", "admin:saumitrasuparn@tcp(alerttradedb.czqug0e2in8p.ap-south-1.rds.amazonaws.com:3306)/alert_trade_db")
	log.SetOutput(log_file)
	if err_db_open != nil {
		log.Println(err_db_open)
	}
	rows, err := db.Query("call stp_get_all_cocodes_cmots()")
	if err != nil {
		log.Println(err)
	}
	list_cocodes := []int{}
	for rows.Next() {
		var cocode int
		rows.Scan(&cocode)
		list_cocodes = append(list_cocodes, cocode)
	}
	for i := 0; i < len(list_cocodes); i++ {
		str_cocode := strconv.Itoa(list_cocodes[i])
		url := fmt.Sprintf("https://insbaapis.cmots.com/api/QuarterlyRatio/%s/S/", str_cocode)
		fmt.Println(url)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Println(err)
		}
		req.Header.Set("Authorization", "Bearer "+api_key)

		res, err := http.DefaultClient.Do(req)

		if err != nil {
			log.Println(err)
		}

		res_body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Println(err)
		}

		// fmt.Println(string(res_body))
		json.Unmarshal(res_body, &unfolded_quaterly)
		data_struct := unfolded_quaterly.Data[0]

		_, err_exec := db.Exec(
			"CALL stp_insert_or_update_quarterly_financials_cmots(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
			data_struct.CoCode, data_struct.Qtrend, data_struct.Mcap, data_struct.Ev, data_struct.Pe,
			data_struct.Pbv, data_struct.Eps, data_struct.Bookvalue, data_struct.Ebit, data_struct.Ebitda,
			data_struct.EvSales, data_struct.EvEbitda, data_struct.Netincomemargin, data_struct.Grossincomemargin,
			data_struct.Ebitdamargin, data_struct.Epsdiluted, data_struct.Netsales, data_struct.Netprofit, data_struct.Cogs,
		)
		if err_exec != nil {
			log.Println("Error executing stored procedure:", err_exec)
		}

	}
	db.Close()
}
func get_quaterly_ratios2(log_file *os.File) {
	var unfolded_quaterly stuct_quaterly_ratios
	db, err_db_open := sql.Open("mysql", "admin:saumitrasuparn@tcp(alerttradedb.czqug0e2in8p.ap-south-1.rds.amazonaws.com:3306)/alert_trade_db")
	log.SetOutput(log_file)
	if err_db_open != nil {
		log.Println(err_db_open)
	}
	rows, err := db.Query("call stp_get_all_cocodes_cmots()")
	if err != nil {
		log.Println(err)
	}
	list_cocodes := []int{}
	for rows.Next() {
		var cocode int
		rows.Scan(&cocode)
		list_cocodes = append(list_cocodes, cocode)
	}
	for i := 0; i < len(list_cocodes); i++ {
		str_cocode := strconv.Itoa(list_cocodes[i])
		url := fmt.Sprintf("https://insbaapis.cmots.com/api/QuarterlyRatio/%s/S/", str_cocode)
		fmt.Println(url)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Println(err)
		}
		req.Header.Set("Authorization", "Bearer "+api_key)

		res, err := http.DefaultClient.Do(req)

		if err != nil {
			log.Println(err)
		}

		res_body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Println(err)
		}

		// fmt.Println(string(res_body))
		json.Unmarshal(res_body, &unfolded_quaterly)
		data_struct := unfolded_quaterly.Data[1]

		_, err_exec := db.Exec(
			"CALL stp_insert_or_update_quarterly_financials_cmots_2(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
			data_struct.CoCode, data_struct.Qtrend, data_struct.Mcap, data_struct.Ev, data_struct.Pe,
			data_struct.Pbv, data_struct.Eps, data_struct.Bookvalue, data_struct.Ebit, data_struct.Ebitda,
			data_struct.EvSales, data_struct.EvEbitda, data_struct.Netincomemargin, data_struct.Grossincomemargin,
			data_struct.Ebitdamargin, data_struct.Epsdiluted, data_struct.Netsales, data_struct.Netprofit, data_struct.Cogs,
		)
		if err_exec != nil {
			log.Println("Error executing stored procedure:", err_exec)
		}

	}
	db.Close()
}

func get_quaterly_ratios3(log_file *os.File) {
	var unfolded_quaterly stuct_quaterly_ratios
	db, err_db_open := sql.Open("mysql", "admin:saumitrasuparn@tcp(alerttradedb.czqug0e2in8p.ap-south-1.rds.amazonaws.com:3306)/alert_trade_db")
	log.SetOutput(log_file)
	if err_db_open != nil {
		log.Println(err_db_open)
	}
	rows, err := db.Query("call stp_get_all_cocodes_cmots()")
	if err != nil {
		log.Println(err)
	}
	list_cocodes := []int{}
	for rows.Next() {
		var cocode int
		rows.Scan(&cocode)
		list_cocodes = append(list_cocodes, cocode)
	}
	for i := 0; i < len(list_cocodes); i++ {
		str_cocode := strconv.Itoa(list_cocodes[i])
		url := fmt.Sprintf("https://insbaapis.cmots.com/api/QuarterlyRatio/%s/S/", str_cocode)
		fmt.Println(url)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Println(err)
		}
		req.Header.Set("Authorization", "Bearer "+api_key)

		res, err := http.DefaultClient.Do(req)

		if err != nil {
			log.Println(err)
		}

		res_body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Println(err)
		}

		// fmt.Println(string(res_body))
		json.Unmarshal(res_body, &unfolded_quaterly)
		data_struct := unfolded_quaterly.Data[2]

		_, err_exec := db.Exec(
			"CALL stp_insert_or_update_quarterly_financials_cmots_3(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
			data_struct.CoCode, data_struct.Qtrend, data_struct.Mcap, data_struct.Ev, data_struct.Pe,
			data_struct.Pbv, data_struct.Eps, data_struct.Bookvalue, data_struct.Ebit, data_struct.Ebitda,
			data_struct.EvSales, data_struct.EvEbitda, data_struct.Netincomemargin, data_struct.Grossincomemargin,
			data_struct.Ebitdamargin, data_struct.Epsdiluted, data_struct.Netsales, data_struct.Netprofit, data_struct.Cogs,
		)
		if err_exec != nil {
			log.Println("Error executing stored procedure:", err_exec)
		}

	}
	db.Close()
}

func get_cash_fow_ratio(log_file *os.File) {
	var cash_ratio_struct cashflow_ratio
	db, err_db_open := sql.Open("mysql", "admin:saumitrasuparn@tcp(alerttradedb.czqug0e2in8p.ap-south-1.rds.amazonaws.com:3306)/alert_trade_db")
	log.SetOutput(log_file)
	if err_db_open != nil {
		log.Println(err_db_open)
	}
	rows, err := db.Query("call stp_get_all_cocodes_cmots()")
	if err != nil {
		log.Println(err)
	}
	list_cocodes := []int{}
	for rows.Next() {
		var cocode int
		rows.Scan(&cocode)
		list_cocodes = append(list_cocodes, cocode)
	}
	for i := 0; i < len(list_cocodes); i++ {
		str_cocode := strconv.Itoa(list_cocodes[i])
		url := fmt.Sprintf("https://insbaapis.cmots.com/api/CashFlowRatios/%s/S/", str_cocode)
		fmt.Println(url)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Println(err)
		}
		req.Header.Set("Authorization", "Bearer "+api_key)

		res, err := http.DefaultClient.Do(req)

		if err != nil {
			log.Println(err)
		}

		res_body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Println(err)
		}

		// fmt.Println(string(res_body))
		json.Unmarshal(res_body, &cash_ratio_struct)
		// fmt.Println(cash_ratio_struct)
		if len(cash_ratio_struct.Data) == 0 {
			error_message := fmt.Sprintf("cashratio_data not found for: %s", str_cocode)
			log.Println(error_message)
			continue
		}
		data_struct := cash_ratio_struct.Data[0]

		_, err_exec := db.Exec(
			"CALL stp_insert_or_update_cash_flow_data_cmots(?,?,?,?,?,?,?,?)",
			data_struct.CoCode, data_struct.Yrc, data_struct.Cashflowpershare,
			data_struct.Pricetocashflowratio, data_struct.Freecashflowpershare,
			data_struct.Pricetofreecashflow, data_struct.Freecashflowyield,
			data_struct.Salestocashflowratio,
		)
		if err_exec != nil {
			log.Println("Error executing stored procedure:", err_exec)
		}

	}
	db.Close()
}
