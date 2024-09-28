package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Todo struct {
    ID        int    `json:"id" bson:"_id"`
    Completed bool   `json:"completed"`
    Body      string `json:"body"`
}

var collection *mongo.Collection

func main() {
    fmt.Println("Server is running...")
    app := fiber.New()

    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file:", err)
    }

    PORT := os.Getenv("PORT")
    MONGODB_URI := os.Getenv("MONGODB_URI")
    clientOptions := options.Client().ApplyURI(MONGODB_URI)
    client, err := mongo.Connect(context.Background(), clientOptions)

    if err != nil {
        log.Fatal(err)
    }

    err = client.Ping(context.Background(), nil)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Connected to MongoDB")

    todos := []Todo{}

    // Read all todos
    app.Get("/api/todos", func(c *fiber.Ctx) error {
        return c.Status(200).JSON(todos)
    })

    // Create a todo
    app.Post("/api/todos", func(c *fiber.Ctx) error {
        todo := &Todo{}

        if err := c.BodyParser(todo); err != nil {
            return err
        }

        if len(todo.Body) > 0 && todo.Body == "" {
            return c.Status(400).JSON(fiber.Map{"error": "Todo body is required"})
        }

        todo.ID = len(todos) + 1
        todos = append(todos, *todo)

        return c.Status(201).JSON(todo)
    })

    // Update a Todo (Partial data: Completed)
    app.Patch("/api/todos/:id", func(c *fiber.Ctx) error {
        id := c.Params("id")

        for i, todo := range todos {
            if fmt.Sprint(todo.ID) == id {
                todos[i].Completed = !todos[i].Completed
                return c.Status(200).JSON(todos[i])
            }
        }
        return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
    })

    // Delete a Todo
    app.Delete("/api/todos/:id", func(c *fiber.Ctx) error {
        id := c.Params("id")
        
        for i, todo := range todos {
            if fmt.Sprint(todo.ID) == id {
                todos = append(todos[:i], todos[i+1:]...)
                return c.Status(200).JSON(fiber.Map{"success": true})
            }
        }
        return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
    })

    log.Fatal(app.Listen(":" + PORT))
}
