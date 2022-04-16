package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func SetHeaders(c *fiber.Ctx) error {
	c.Type("json", "utf-8")
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Headers", "*")
	c.Set("Access-Control-Max-Age", "1728000")

	return c.Next()
}

func HandleHealth(c *fiber.Ctx) error {
	defer calc()()

	return c.SendStatus(204)
}

func HandleNext(c *fiber.Ctx) error {
	defer calc()()

	res, status := FetchNext(c.Params("id"))
	return c.Status(status).SendString(res)
}

func HandleBrowse(c *fiber.Ctx) error {
	defer calc()()

	id := c.Params("id")

	switch id[:2] {
	case "UC":
		res, status := FetchArtist(id)
		return c.Status(status).SendString(res)
	default:
		return c.SendString("{error: \"Invalid URL\"}")
	}
}

func HandleArtist(c *fiber.Ctx) error {
	defer calc()()
	res, status := FetchArtist(c.Params("id"))
	return c.Status(status).SendString(res)
}

func main() {
	app := fiber.New()

	app.Use(SetHeaders)
	app.Use(recover.New())

	app.Get("/healthz", HandleHealth)
	app.Get("/next/:id", HandleNext)
	app.Get("/browse/:id", HandleBrowse)
	app.Get("/channel/:id", HandleArtist)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	app.Listen(":" + port)
}
