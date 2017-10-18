package main

import "fmt"

func main() {
	// will research a way to do it without typing all elements
	ints := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for _, n := range ints {
		// figuring out a way to define out took some minutes
		// using else wouldn't work :)
		out := fmt.Sprintf("%d is odd", n)
		if n%2 == 0 {
			out = fmt.Sprintf("%d is even", n)
		}
		fmt.Println(out)
	}
}
