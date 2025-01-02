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
    ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
    Completed bool `json:"completed"`
    Task string `json:"task"`
}

var collection *mongo.Collection

func main() {
    if err := godotenv.Load(); err != nil {
        log.Println("Couldn't load env")
    }

    uri := os.Getenv("MONGODB_URI")

    if uri == "" {
        log.Fatal("Couldn't find MONGO_URI")
    }


    clientOptions := options.Client().ApplyURI(uri)

    client, err:= mongo.Connect(context.Background(), clientOptions)

    if err != nil {
        panic(err)
    }

    defer client.Disconnect(context.Background())

    if err := client.Ping(context.Background(), nil); err != nil {
        panic(err)
    }

    fmt.Println("Success: Connected to MongoDB Atlas")

    collection = client.Database("golang_db").Collection("todos")

    app := fiber.New()

    app.Get("/api/todo", getTodos)
    app.Post("/api/todo", createTodo)
    app.Patch("/api/todo/:id", updateTodo)
    app.Delete("/api/todo/:id", deleteTodo)


    port := os.Getenv("PORT")

    if port == "" {
        port = "5000"
    }

    log.Fatal(app.Listen(":" + port))
}

func getTodos (c *fiber.Ctx) error {
    var todos []Todo

    cursor, err := collection.Find(context.Background(), bson.M{})

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

func createTodo (c *fiber.Ctx) error {
    todo := new(Todo)

    if err := c.BodyParser(todo); err != nil {
        return err
    }

    if todo.Task == "" {
        return c.Status(400).JSON(fiber.Map{"error": "Todo task cannot be empty"})
    }

    insertResult, err := collection.InsertOne(context.Background(), todo)

    if err != nil {
        return err
    }

    todo.ID = insertResult.InsertedID.(primitive.ObjectID)

    return c.Status(201).JSON(todo)
}
func updateTodo (c *fiber.Ctx) error {
    id := c.Params("id")

    objectID, err := primitive.ObjectIDFromHex(id)

    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "invalid Todo id"})
    }

    filter := bson.M{"_id": objectID}

    todo := new(Todo)
    err = collection.FindOne(context.Background(), filter).Decode(&todo)

    if err != nil {
        return err
    }

    update := bson.M{"$set": bson.M{"completed": !todo.Completed}}

    _, err = collection.UpdateOne(context.Background(), filter, update)
    if err != nil {
        return err
    }

    return c.Status(200).JSON(fiber.Map{"success": "true"})

}
func deleteTodo (c *fiber.Ctx) error {
    id := c.Params("id")

    objectID, err := primitive.ObjectIDFromHex(id)

    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error":"invalid Todo id"})
    }

    filter := bson.M{"_id": objectID}

    _, err = collection.DeleteOne(context.Background(), filter)

    if err != nil {
        return err
    }

    return c.Status(200).JSON(fiber.Map{"success": "true"})

}