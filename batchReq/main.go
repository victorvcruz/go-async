package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type req struct {
	valor string
}

func process(batch []req) {
	fmt.Println("processando batch com valores: ", batch)
	//for _, value := range batch {
	//	fmt.Println("processando valor: ", value)
	//}
	time.Sleep(3 * time.Second)
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

type res struct {
	in chan req
}

func main() {
	app := fiber.New()
	response := &res{in: make(chan req)}
	back := processingBatch(response.in, 2)
	done := processBatch(back)

	app.Get("/", response.Request)

	go app.Listen(":3000")
	Response()

	time.Sleep(10 * time.Minute)
	close(response.in)
	<-done
}

func (r *res) Request(c *fiber.Ctx) error {
	r.toQueue(c.Query("value"))
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "ok"})
}

func (r *res) toQueue(value string) {
	r.in <- req{valor: value}
}

func Response() string {
	var wg sync.WaitGroup

	cpf1 := "?value=aa"
	cpf2 := "?value=bb"
	cpf3 := "?value=cc"
	cpf4 := "?value=dd"
	start := time.Now()
	for i := 0; i < 60; i++ {
		var cpf string
		wg.Add(1)
		cpf = cpf1
		if i > 15 && i <= 30 {
			cpf = cpf2
		} else if i > 30 && i <= 45 {
			cpf = cpf3
		} else if i > 45 && i <= 60 {
			cpf = cpf4
		}
		go func() {
			defer wg.Done()
			client := &http.Client{}

			req, err := http.NewRequest("GET", "http://localhost:3000"+cpf, nil)
			if err != nil {
				log.Println(err)
			}

			resp, err := client.Do(req)
			if err != nil {
			}

			_, err = io.ReadAll(resp.Body)
			if err != nil {
				log.Println(err)
			}
		}()
	}

	wg.Wait()
	fmt.Println("START - " + start.String())
	fmt.Println("END - " + time.Now().String())
	fmt.Println("SINCE - " + time.Since(start).String())

	return ""
}
