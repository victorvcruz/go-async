package main

import (
	"fmt"
	"time"
)

type req struct {
	valor int
}

func process(batch []req) {
	fmt.Println("processando batch com valores: ", batch)
	for _, value := range batch {
		fmt.Println("processando valor: ", value)
		time.Sleep(2 * time.Second)
	}
}

func processBatch(in <-chan []req) chan bool {
	done := make(chan bool)

	go func() {
		for batch := range in {
			process(batch)
		}
		done <- true
	}()
	return done
}

func processingBatch(in <-chan req, lengthBatch int) chan []req {
	back := make(chan []req)
	go func() {
		buf := make([]req, 0, lengthBatch)
		valueZero := req{}
		var closed bool
		for !closed {
			var toOut bool

			select {
			case r := <-in:
				if r == valueZero {
					closed = true
					continue
				}
				buf = append(buf, r)
				toOut = len(buf) == lengthBatch
			}
			if toOut {
				back <- buf
				buf = make([]req, 0, lengthBatch)
			}
		}
		if len(buf) > 0 {
			back <- buf
		}
		close(back)
	}()
	return back
}

func main() {
	in := make(chan req)
	back := processingBatch(in, 3)
	done := processBatch(back)
	in <- req{valor: 1}
	in <- req{valor: 2}
	in <- req{valor: 3}
	in <- req{valor: 4}
	in <- req{valor: 5}
	in <- req{valor: 6}
	in <- req{valor: 7}
	in <- req{valor: 8}
	in <- req{valor: 9}
	in <- req{valor: 10}

	close(in)
	<-done
}
