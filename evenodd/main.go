package main

import "fmt"

func main() {
	// will research a way to do it without typing all elements
	ints := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for _, n := range ints {
		if n%2 == 0 {
			fmt.Println(n, "is even")
		} else {
			fmt.Println(n, "is odd")
		}
	}
}
