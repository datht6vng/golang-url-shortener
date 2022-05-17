package main

// add engine
import (
	"fmt"
	"server_go/controller"
	"server_go/limiter"
	"server_go/unittest"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

func main() {
	// ------------------------------------------------------------
	unittest.TestShortener()
	// unittest.TestValidUrl()
	// ------------------------------------------------------------

	viewEngine := html.New("./views", ".html")
	// Create a new Fiber template with template engine
	app := fiber.New(fiber.Config{
		Views: viewEngine,
	})
	// Connect database
	controller.Model.Connect()
	controller.Cache.Connect()
	// Limiter
	app.Use(limiter.CreateLimiter())
	// Controller
	app.Get("/", controller.GetIndexController)
	app.Post("/gen-url", controller.PostGenUrlController)
	app.Get("/init-db", controller.InitDBController)
	app.Get("/reset", controller.GetReset)
	app.Get("/:url", controller.GetUrlController)
	app.Listen("localhost:8080")
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
		controller.Model.Close()
	}()
}
