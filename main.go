package main

import (
	"fmt"
	"os"

	"codeberg.org/Hyperpipe/hyperpipe-backend/lib"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func SetHeaders(c *fiber.Ctx) error {
	c.Type("json", "utf-8")
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Headers", "*")
	c.Set("Access-Control-Max-Age", "300")

	return c.Next()
}

func HandleHealth(c *fiber.Ctx) error {
	fmt.Println("Health Check!!")

	return c.SendStatus(204)
}

func HandleExplore(c *fiber.Ctx) error {
	res, status := lib.GetExplore()

	return c.Status(status).SendString(res)
}

func HandleNext(c *fiber.Ctx) error {
	res, status := lib.GetNext(c.Params("id"))

	return c.Status(status).SendString(res)
}

func HandleGenres(c *fiber.Ctx) error {
	res, status := lib.GetGenres()

	return c.Status(status).SendString(res)
}

func HandleGenre(c *fiber.Ctx) error {
	res, status := lib.GetGenre(c.Params("id"))

	return c.Status(status).SendString(res)
}

func HandleCharts(c *fiber.Ctx) error {
	res, status := lib.GetCharts(c.Query("params"), c.Query("code"))

	return c.Status(status).SendString(res)
}

func HandleLyrics(c *fiber.Ctx) error {
	res, status := lib.GetLyrics(c.Params("id"))

	return c.Status(status).SendString(res)
}

func HandleBrowse(c *fiber.Ctx) error {

	// Will be removed, Do not use

	id := c.Params("id")

	if len(id) < 4 {
		return c.Status(500).SendString("{\"error\": \"browse id is too short\"}")
	}

	switch {
	case id[:2] == "UC":
		res, status := lib.GetArtist(id)
		return c.Status(status).SendString(res)
	case id[:4] == "MPLY":
		res, status := lib.GetLyrics(id)
		return c.Status(status).SendString(res)
	default:
		return c.Status(500).SendString("{\"error\": \"URL not supported\"}")
	}
}

func HandleArtist(c *fiber.Ctx) error {
	res, status := lib.GetArtist(c.Params("id"))

	return c.Status(status).SendString(res)
}

func main() {
	if os.Getenv("HYP_PROXY") == "" {
		fmt.Println("HYP_PROXY is empty!")
		os.Exit(1)
	}

	app := fiber.New()

	app.Use(SetHeaders)
	app.Use(recover.New())

	app.Get("/healthz", HandleHealth)
	app.Get("/explore", HandleExplore)
	app.Get("/genres", HandleGenres)
	app.Get("/genres/:id", HandleGenre)
	app.Get("/charts", HandleCharts)
	app.Get("/next/:id", HandleNext)
	app.Get("/lyrics/:id", HandleLyrics)
	app.Get("/browse/:id", HandleBrowse)
	app.Get("/channel/:id", HandleArtist)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	app.Listen(":" + port)
}
