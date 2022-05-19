package main

// add engine
import (
	"fmt"
	"os"
	"server_go/controller"
	"server_go/limiter"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html"
)

func main() {
	// ------------------------------------------------------------
	// unittest.TestShortener()
	// unittest.TestValidUrl()
	// ------------------------------------------------------------
	viewEngine := html.New("./views", ".html")
	// Create a new Fiber template with template engine
	app := fiber.New(fiber.Config{
		Views: viewEngine,
	})
	controller.BaseUrl = os.Getenv("DOMAIN")
	if controller.BaseUrl == "" {
		controller.BaseUrl = "http://localhost:8080/"
	}
	// Connect database, cache
	controller.Model.Connect()
	controller.Cache.Connect()
	// Limiter

	app.Use(limiter.CreateLimiter())
	// Cors
	app.Use(cors.New(cors.Config{
		Next:             nil,
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "",
		AllowCredentials: false,
		ExposeHeaders:    "",
		MaxAge:           0,
	}))
	// Controller
	app.Get("/", controller.GetIndexController)
	app.Post("/gen-url", controller.PostGenUrlController)
	app.Get("/init-db", controller.InitDBController)
	app.Get("/reset-cache", controller.GetResetCache)
	app.Get("/reset-db", controller.GetResetDB)
	app.Get("/:url", controller.ValidateController, controller.GetUrlController)
	app.Listen("0.0.0.0:8080")
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
		controller.Model.Close()
	}()
}
