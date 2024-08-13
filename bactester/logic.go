package main

func Crossover(indicators_and_close [][]float64, indexes []int) []int {
	indexes_at_position := []int{}
	indicator_c1 := indexes[0]
	indicator_c2 := indexes[1]
	indicator1 := indicators_and_close[indicator_c1]
	indicator2 := indicators_and_close[indicator_c2]

	for i := 1; i < len(indicator1); i++ {
		if indicator1[i-1] < indicator2[i-1] && indicator1[i] > indicator2[i] {
			indexes_at_position = append(indexes_at_position, i)
		}
	}
	return indexes_at_position
}

func Greater_than(indicators_and_price [][]float64, indexes []int) []int {
	indexes_at_position := []int{}
	first_cross_indicator := indicators_and_price[indexes[0]]
	second_cross_indicator_ := indicators_and_price[indexes[1]]
	for i := 0; i < len(second_cross_indicator_); i++ {
		if first_cross_indicator[i] > second_cross_indicator_[i] {
			indexes_at_position = append(indexes_at_position, i)
		}
	}
	return indexes_at_position
}
func Less_than(indicators_and_price [][]float64, indexes []int) []int {
	indexes_at_position := []int{}
	first_cross_indicator := indicators_and_price[indexes[0]]
	second_cross_indicator_ := indicators_and_price[indexes[1]]
	for i := 0; i < len(second_cross_indicator_); i++ {
		if first_cross_indicator[i] < second_cross_indicator_[i] {
			indexes_at_position = append(indexes_at_position, i)
		}
	}
	return indexes_at_position
}
