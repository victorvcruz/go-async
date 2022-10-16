package main

import (
	"github.com/gofiber/fiber/v2"
	"log"
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

	app.Listen(":3000")
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
