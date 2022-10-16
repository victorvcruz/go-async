package main

import (
	"fmt"
)

func main() {
	out := make(chan int, 3)
	in := make(chan int, 3)
	pair := make(chan bool, 3)
	outPair := make(chan int, 3)

	go multiplyByThree(in, out, outPair)
	go multiplyByThree(in, out, outPair)
	go multiplyByThree(in, out, outPair)

	go isPair(outPair, pair)
	go isPair(outPair, pair)
	go isPair(outPair, pair)

	go func() {
		in <- 1
		in <- 2
		in <- 3
	}()

	fmt.Printf("Is pair? %t\n", <-pair)
	fmt.Printf("Is pair? %t\n", <-pair)
	fmt.Printf("Is pair? %t\n", <-pair)

	fmt.Println(<-out)
	fmt.Println(<-out)
	fmt.Println(<-out)

}

func multiplyByThree(in <-chan int, out chan<- int, outPair chan<- int) {
	fmt.Println("Initializing goroutine multiply...")
	for {
		num := <-in
		result := num * 3
		out <- result
		outPair <- result
	}
}

func isPair(out <-chan int, pair chan<- bool) {
	fmt.Println("Initializing goroutine pair...")
	for {
		num := <-out
		result := num % 2
		if result == 0 {
			pair <- true
		} else {
			pair <- false
		}
	}
}
