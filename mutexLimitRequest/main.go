package main

import (
	"github.com/gofiber/fiber/v2"
	"io"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type Mutex struct {
	mu           *sync.Mutex
	count        int64
	isProcessing int32
}

func main() {

	mutex := Mutex{mu: &sync.Mutex{}}

	app := fiber.New()

	app.Get("/", mutex.Request)

	go app.Listen(":3000")

	go Response()

	time.Sleep(10 * time.Minute)
}

func (m *Mutex) Request(c *fiber.Ctx) error {
	atomic.AddInt64(&m.count, 1)
	atomic.StoreInt32(&m.isProcessing, 0)

	if m.count >= 10 {

		m.mu.Lock()
		defer func() {
			// atomic.StoreInt64(&m.count, 0)
			// atomic.AddInt64(&m.count, -1)
			log.Printf("UNLOCK - %d\n", m.count)
			m.mu.Unlock()
		}()
	}

	log.Printf("COUNT : %d\n", m.count)
	time.Sleep(3 * time.Second)
	atomic.AddInt64(&m.count, -1)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "ok"})
}

func Response() string {

	var wg sync.WaitGroup

	for i := 0; i < 60; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := &http.Client{}

			req, err := http.NewRequest("GET", "http://localhost:3000", nil)
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
	log.Println("FINISH")
	return ""
}
