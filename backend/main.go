package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/gofiber/fiber/v2"
    "github.com/joho/godotenv"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type Todo struct {
    ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"` //omitempty to prevent 000...000 id
    Completed bool               `json:"completed"`
    Body      string             `json:"body"`
}

var collection *mongo.Collection

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file:", err)
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
        PORT = "4000"
    }
    log.Fatal(app.Listen(":" + PORT))
}

// TODO : Get all Todos
func getTodos(c *fiber.Ctx) error {
    var todos []Todo

    cursor, err := collection.Find(context.Background(), bson.M{}) // No Filters {}
    if err != nil {
        return err
    }

    defer cursor.Close(context.Background())

    for cursor.Next(context.Background()) {
        var todo Todo
        if err := cursor.Decode(&todo); err != nil {
            return err
        }
        todos = append(todos, todo)
    }

    return c.JSON(todos)
}

// TODO : Create a todo
func createTodos(c *fiber.Ctx) error {
    todo := new(Todo)

    if err := c.BodyParser(todo); err != nil {
        return err
    }

    if len(todo.Body) > 0 && todo.Body == "" {
        return c.Status(404).JSON(fiber.Map{"error": "Todo body cannot be empty"})
    }

    insertResult, err := collection.InsertOne(context.Background(), todo)
    if err != nil {
        return err
    }

    todo.ID = insertResult.InsertedID.(primitive.ObjectID)

    return c.Status(201).JSON(todo)
}

// TODO : Update a todo by id (Complete property)
func updateTodos(c *fiber.Ctx) error {
    id := c.Params("id")
    ObjectID, err := primitive.ObjectIDFromHex(id)

    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Invalid todo ID"})
    }

    filter := bson.M{"_id": ObjectID}
    update := bson.M{"$set": bson.M{"completed": true}}
    _, err = collection.UpdateOne(context.Background(), filter, update)

    if err != nil {
        return err
    }

    return c.Status(200).JSON(fiber.Map{"success": true})
}

// TODO : Delete a todo by id
func deleteTodos(c *fiber.Ctx) error {
    id := c.Params("id")
    objectID, err := primitive.ObjectIDFromHex(id)

    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Invalid todo ID"})
    }

    filter := bson.M{"_id": objectID}
    _, err = collection.DeleteOne(context.Background(), filter)

    if err != nil {
        return err
    }

    return c.Status(200).JSON(fiber.Map{"success": true})
}
