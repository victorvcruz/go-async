package main

import (
	"fmt"
	"sync"
	"time"
)

var amount int
var count int

func printGo(wg *sync.WaitGroup, mu *sync.Mutex, value int) {
	mu.Lock()
	defer func() {
		count++
		if count == 2 {
			time.Sleep(5 * time.Second)
			count = 0
		}
		wg.Done()
		mu.Unlock()
	}()

	for i := 0; i < 3; i++ {
		amount += value
		fmt.Printf("%d - %d\n", value, amount)
	}
}

func main() {
	fmt.Println("START")
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := 0; i < 7; i++ {
		wg.Add(1)
		go printGo(&wg, &mu, 10*i)
	}

	wg.Wait()
	fmt.Println("FINISH")

}
