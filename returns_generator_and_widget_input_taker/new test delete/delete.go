package main

import "fmt"

func main() {
	final_timestamps := []int{1, 2, 78, 4, 4, 80, 65, 7, 8}
	final_timestamps = append(final_timestamps[:2], final_timestamps[2+1:]...)
	fmt.Println(final_timestamps)
	InsertionSort(final_timestamps)
	fmt.Println(final_timestamps)
	x := 10
	sdd(x)
	fmt.Println(x)

}
func InsertionSort(arr []int) {
	for i := 1; i < len(arr); i++ {
		key := arr[i]
		j := i - 1

		// Move elements of arr[0..i-1], that are greater than key,
		// to one position ahead of their current position
		for j >= 0 && arr[j] > key {
			arr[j+1] = arr[j]
			j = j - 1
		}
		arr[j+1] = key
	}
}

func sdd(a int) {
	a = a + 1
}


