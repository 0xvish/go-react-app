package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

type Todo struct {
    ID        int    `json:"id"`
    Completed bool   `json:"completed"`
    Task      string `json:"task"`
}

func main() {
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

    log.Fatal(app.Listen(":3000"))



}