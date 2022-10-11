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
	var AAAAAAAAAAAAAA sync.WaitGroup
	var mu sync.Mutex

	for i := 0; i < 7; i++ {
		AAAAAAAAAAAAAA.Add(1)
		go printGo(&AAAAAAAAAAAAAA, &mu, 10*i)
	}

	AAAAAAAAAAAAAA.Wait()
	fmt.Println("FINISH")

}
