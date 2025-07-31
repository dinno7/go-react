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
	Body      string `json:"body"`
	Completed bool   `json:"completed"`
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Err in loading Envs")
	}

	app := fiber.New()
	api := app.Group("/api")
	apiV1 := api.Group("/v1")

	todos := []Todo{}

	todosApi := apiV1.Group("/todos")

	todosApi.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"message": "Here are", "meta": todos})
	})

	todosApi.Post("/", func(c *fiber.Ctx) error {
		todo := &Todo{}

		if err := c.BodyParser(todo); err != nil {
			fmt.Println("💀 err > ", err)
			return c.Status(fiber.StatusBadRequest).
				JSON(fiber.Map{"error": "Please provide valid todo object"})
		}

		if len(todo.Body) <= 0 {
			return c.Status(fiber.StatusBadRequest).
				JSON(fiber.Map{"error": "The body of todo can not be empty"})
		}

		todo.ID = len(todos) + 1
		todo.Completed = false

		todos = append(todos, *todo)

		return c.Status(200).
			JSON(fiber.Map{"message": "New todo added", "meta": map[string]any{"todo": todo}})
	})

	todosApi.Patch("/:id", func(c *fiber.Ctx) error {
		p := c.Params("id", "0")
		if p == "" || p == "0" {
			return c.Status(fiber.StatusBadRequest).
				JSON(fiber.Map{"error": "Please provide valid id"})
		}

		foundTodo := &Todo{}
		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == p {
				todos[i].Completed = !todo.Completed
				*foundTodo = todo
				break
			}
		}

		return c.Status(200).
			JSON(fiber.Map{"message": "The todo toggle successfully", "meta": map[string]any{
				"todo": foundTodo,
			}})
	})

	todosApi.Delete("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id", "0")
		if id == "" || id == "0" {
			return c.Status(fiber.StatusBadRequest).
				JSON(fiber.Map{"error": "Please provide valid id"})
		}

		deleted := false
		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				todos = append(todos[:i], todos[i+1:]...)
				deleted = true
			}
		}

		if !deleted {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"ok":    false,
				"error": fmt.Sprintf("The todo with id %s is not exists", id),
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"ok":      true,
			"message": fmt.Sprintf("The todo with %s has been deleted successfully", id),
		})
	})

	PORT := os.Getenv("PORT")
	log.Fatal(app.Listen(":" + PORT))
}
