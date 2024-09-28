package main

import (
    "fmt"
    "log"

    "github.com/gofiber/fiber/v2"
)

func main() {
    fmt.Println("Server is running...")
    app := fiber.New()

    app.Get("/", func(c *fiber.Ctx) error {
        return c.Status(200).JSON(fiber.Map{"data": "Hello World"})
    })
    log.Fatal(app.Listen(":4000"))
}