package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type Todo struct {
    ID        int    `json:"id"`
    Completed bool   `json:"completed"`
    Task      string `json:"task"`
}

func main() {
    err := godotenv.Load(".env"); if err != nil {
        log.Fatal("error loading .env file")
    }

    PORT := os.Getenv("PORT")
	app := fiber.New()

    todoList := []Todo{}

    app.Get("/", func(c *fiber.Ctx) error {
        return c.Status(200).JSON(fiber.Map{"message": "Hello, World!"})
    })

    app.Post("/api/todo", func(c *fiber.Ctx) error {
        todo := &Todo{}

        if err := c.BodyParser(todo); err !=nil {
            return err
        }

        if todo.Task == "" {
            return c.Status(400).JSON(fiber.Map{"error": "Task is required"})
        }

        todo.ID = len(todoList) + 1
        todoList = append(todoList, *todo)

        return c.Status(201).JSON(todo)
    })

    //return all todos
    app.Get("/api/todo", func(c *fiber.Ctx) error {
        return c.Status(200).JSON(todoList)
    })

    //update a todo
    app.Patch("/api/todo/:id", func (c *fiber.Ctx) error {
        id := c.Params("id")
        todo := &Todo{}

        if err := c.BodyParser(todo); err != nil {
            return err
        }

        for i, todo := range todoList {
            if fmt.Sprint(todo.ID) == id {
                todoList[i].Completed = !todoList[i].Completed
                return c.Status(200).JSON(todoList[i])
            }
        }

        return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
    })

    //delete a todo
    app.Delete("/api/todo/:id", func (c *fiber.Ctx) error {
        id := c.Params("id")

        for i, todo := range todoList {
            if fmt.Sprint(todo.ID) == id {
                todoList = append(todoList[:i], todoList[i+1:]...)
                return c.Status(404).JSON(fiber.Map{"success": true})
            }
        }
        return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
    })

    log.Fatal(app.Listen(":" + PORT))



}