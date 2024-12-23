package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/jaunewick/go-react-task-app/routes"
    "github.com/joho/godotenv"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

func main() {
    // Load .env file in development
    if os.Getenv("ENV") != "production" {
        err := godotenv.Load()
        if err != nil {
            log.Fatal("Error loading .env file:", err)
        }
    }

    // Connect to MongoDB
    MONGODB_URI := os.Getenv("MONGODB_URI")
    clientOptions := options.Client().ApplyURI(MONGODB_URI)
    client, err := mongo.Connect(context.Background(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }

    // Disconnect when shutdown server
    defer client.Disconnect(context.Background())

    // Ping the server
    err = client.Ping(context.Background(), nil)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Connected to MongoDB")

    // Set collection to routes
    collection = client.Database("golang_db").Collection("todos")
    routes.SetCollection(collection)

    // Create a new Fiber instance
    fmt.Println("Server is running...")
    app := fiber.New()

    // CORS
    // app.Use(cors.New(cors.Config{
    //     // Frontend URL
    //     AllowOrigins: "http://localhost:5173",
    //     AllowHeaders: "Origin, Content-Type, Accept",
    // }))

    // Define routes
    app.Get("/api/todos", routes.GetTodos)
    app.Post("/api/todos", routes.CreateTodos)
    app.Patch("/api/todos/:id", routes.UpdateTodos)
    app.Delete("/api/todos/:id", routes.DeleteTodos)

    // Listen on server
    PORT := os.Getenv("PORT")
    if PORT == "" {
        PORT = "4000"
    }

    if os.Getenv("ENV") == "production" {
        app.Static("/", "./frontend/dist")
    }
    log.Fatal(app.Listen(":" + PORT))
}
