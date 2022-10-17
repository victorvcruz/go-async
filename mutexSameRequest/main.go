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

type Mutex struct {
	mu map[string]*sync.Mutex
}

func main() {

	mutexMap := make(map[string]*sync.Mutex)

	mutex := Mutex{mu: mutexMap}

	app := fiber.New()

	app.Get("/", mutex.Request)

	go app.Listen(":3000")

	go Response()

	time.Sleep(10 * time.Minute)
}

func (m *Mutex) Request(c *fiber.Ctx) error {
	cpf := c.Query("cpf")

	mu := m.lock(cpf)
	defer func() {
		mu.Unlock()
		delete(m.mu, cpf)
		log.Println("UNLOCK - " + cpf)
	}()

	time.Sleep(3 * time.Second)

	log.Println("REQUEST OK - " + cpf)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "ok"})
}

func (m *Mutex) lock(_cpf string) *sync.Mutex {
	for cpf, mu := range m.mu {
		if cpf == _cpf {
			mu.Lock()
			return mu
		}
	}

	m.mu[_cpf] = &sync.Mutex{}
	m.mu[_cpf].Lock()
	return m.mu[_cpf]
}

func Response() string {
	var wg sync.WaitGroup

	cpf1 := "?cpf=70099063190"
	cpf2 := "?cpf=70099062119"
	cpf3 := "?cpf=70099060411"
	cpf4 := "?cpf=70099061133"
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
