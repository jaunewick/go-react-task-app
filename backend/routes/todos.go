package routes

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Define model
type Todo struct {
    ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"` //omitempty to prevent 000...000 id
    Completed bool               `json:"completed"`
    Body      string             `json:"body"`
}

// Define collection
var collection *mongo.Collection

// Set collection from main.go
func SetCollection(col *mongo.Collection) {
    collection = col
}

// Get all todos
func GetTodos(c *fiber.Ctx) error {
    var todos []Todo

    cursor, err := collection.Find(context.Background(), bson.M{})
    if err != nil {
        return err
    }

    // Close the cursor once finished to free up resources
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

// Create a todo
func CreateTodos(c *fiber.Ctx) error {
    todo := new(Todo)

    if err := c.BodyParser(todo); err != nil {
        return err
    }

    // Validate body field is not empty
    if strings.TrimSpace(todo.Body) == "" {
        return c.Status(400).JSON(fiber.Map{"error": "Todo body cannot be empty"})
    }

    insertResult, err := collection.InsertOne(context.Background(), todo)
    if err != nil {
        return err
    }
    todo.ID = insertResult.InsertedID.(primitive.ObjectID)

    return c.Status(201).JSON(todo)
}

// Update a todo by id
func UpdateTodos(c *fiber.Ctx) error {
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

// Delete a todo by id
func DeleteTodos(c *fiber.Ctx) error {
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
