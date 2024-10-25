package main

import "github.com/markcheno/go-talib"

type indicator struct {
	signature string
	function  func([]float64, []int) [][]float64
}

func map_creator() map[string]indicator {
	Master_object_indicator := make(map[string]indicator)
	Master_object_indicator["ma"] = indicator{signature: "ma", function: ma}
	Master_object_indicator["rsi"] = indicator{signature: "rsi", function: rsi}
	Master_object_indicator["macd"] = indicator{signature: "macd", function: macd}
	return Master_object_indicator
}

var ma = func(arr []float64, args []int) [][]float64 {
	ma_calc := talib.Ma(arr, args[0], talib.SMA)
	var list_returns [][]float64
	list_returns = append(list_returns, ma_calc)
	return list_returns
}

var rsi = func(arr []float64, arrgs []int) [][]float64 {
	ma_calc := talib.Rsi(arr, arrgs[0])
	var list_returns [][]float64
	list_returns = append(list_returns, ma_calc)
	return list_returns
}

var macd = func(arr []float64, arrgs []int) [][]float64 {
	ma_calc, ma_calc1, ma_calc2 := talib.Macd(arr, arrgs[0], arrgs[1], arrgs[2])
	var list_returns [][]float64
	list_returns = append(list_returns, ma_calc, ma_calc1, ma_calc2)
	return list_returns
}
