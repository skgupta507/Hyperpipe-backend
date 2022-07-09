package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"os"
)

func SetHeaders(c *fiber.Ctx) error {
	c.Type("json", "utf-8")
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Headers", "*")
	c.Set("Access-Control-Max-Age", "1728000")

	return c.Next()
}

func HandleHealth(c *fiber.Ctx) error {
	defer calc(c.OriginalURL())()

	fmt.Println("Health Check!!")

	return c.SendStatus(204)
}

func HandleExplore(c *fiber.Ctx) error {
	defer calc(c.OriginalURL())()

	res, status := FetchExplore()

	return c.Status(status).SendString(res)
}

func HandleNext(c *fiber.Ctx) error {
	defer calc(c.OriginalURL())()

	res, status := FetchNext(c.Params("id"))

	return c.Status(status).SendString(res)
}

func HandleGenres(c *fiber.Ctx) error {
	defer calc(c.OriginalURL())()

	res, status := FetchGenres()

	return c.Status(status).SendString(res)
}

func HandleGenre(c *fiber.Ctx) error {
	defer calc(c.OriginalURL())()

	res, status := FetchGenre(c.Params("id"))

	return c.Status(status).SendString(res)
}

func HandleBrowse(c *fiber.Ctx) error {
	defer calc(c.OriginalURL())()

	id := c.Params("id")

	if len(id) < 4 {
		return c.Status(500).SendString("{\"error\": \"Browse Id is too Short\"}")
	}

	switch {
	case id[:2] == "UC":
		res, status := FetchArtist(id)
		return c.Status(status).SendString(res)
	case id[:4] == "MPLY":
		res, status := FetchLyrics(id)
		return c.Status(status).SendString(res)
	case id[:4] == "MPRE":
		res, status := FetchAlbum(id)
		return c.Status(status).SendString(res)
	default:
		return c.Status(500).SendString("{\"error\": \"Invalid Browse URL\"}")
	}
}

func HandleArtist(c *fiber.Ctx) error {
	defer calc(c.OriginalURL())()

	res, status := FetchArtist(c.Params("id"))

	return c.Status(status).SendString(res)
}

func main() {
	app := fiber.New()

	app.Use(SetHeaders)
	app.Use(recover.New())

	app.Get("/healthz", HandleHealth)
	app.Get("/explore", HandleExplore)
	app.Get("/genres", HandleGenres)
	app.Get("/genres/:id", HandleGenre)
	app.Get("/next/:id", HandleNext)
	app.Get("/browse/:id", HandleBrowse)
	app.Get("/channel/:id", HandleArtist)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	app.Listen(":" + port)
}
