package main

import (
	"os"

	"github.com/gofiber/fiber"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) {
		c.Send("Cool, Worldd!")
	})

	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "3010"
	}

	app.Listen(port)
}
