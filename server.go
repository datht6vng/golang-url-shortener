package main

// add engine
import (
	"fmt"
	"server_go/controller"
	"server_go/limiter"
	"server_go/model"
	"server_go/unittest"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

func main() {
	// test
	unittest.TestShortener()
	viewEngine := html.New("./views", ".html")
	// Create a new Fiber template with template engine
	app := fiber.New(fiber.Config{
		Views: viewEngine,
	})
	model.Connect()
	// Limiter
	app.Use(limiter.CreateLimiter())
	app.Get("/", controller.GetIndexController)
	app.Post("/gen-url", controller.PostGenUrlController)
	app.Get("/init-db", controller.InitDBController)
	app.Get("/:url", controller.GetUrlController)
	app.Listen("localhost:8080")
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
		model.Connection.Close()
	}()
}
