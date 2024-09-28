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
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file:", err)
    }

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
    
    collection = client.Database("golang_db").Collection("todos")

    fmt.Println("Server is running...")
    app := fiber.New()

    app.Get("/api/todos", getTodos)
    app.Post("/api/todos", createTodos)
    app.Patch("/api/todos/:id", updateTodos)
    app.Delete("/api/todos/:id", deleteTodos)

    PORT := os.Getenv("PORT")
    if PORT == "" {
        PORT = "5000"
    }
    log.Fatal(app.Listen(":" + PORT))
}

// TODO : Get all Todos
func getTodos(c *fiber.Ctx) error {
    return nil
}

// TODO : Create a todo
func createTodos(c *fiber.Ctx) error {
    return nil
}

// TODO : Update a todo by id (Complete property)
func updateTodos(c *fiber.Ctx) error {
    return nil
}

// TODO : Delete a todo by id
func deleteTodos(c *fiber.Ctx) error {
    return nil
}