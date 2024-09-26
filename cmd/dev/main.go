package main

import (
	"github.com/KingSupermarket/router"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	router.NewRouter(app)
	app.Listen(":3251")
}
