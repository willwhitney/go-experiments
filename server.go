package main

import (
	"net/http"
	"fmt"
)

func factor(n int) (factors []int) {
	for i := 1; i < n/2 + 1; i++ {
		if n % i == 0 {
			factors = append(factors, i)
		}
	}
	return
}

func handler(writer http.ResponseWriter, request *http.Request) {
	ch := make(chan []int, 1000)
	for i := 1000000; i < 1001000; i++ {

		go func(n int, ch chan []int) {
			ch <- factor(n)
		}(i, ch)

	}

	received := 0
	for received < 1000 {
		fmt.Fprintf(writer, "Factors: %v \n", <-ch)
		received++
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
