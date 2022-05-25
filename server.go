package main

// add engine
import (
	"os"
	"server_go/controller"
	"server_go/limiter"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html"
)

func main() {
	// ------------------------------------------------------------
	// unittest.TestShortener()
	// unittest.TestValidUrl()
	// unittest.TestTrimTimeStamp()
	// ------------------------------------------------------------
	viewEngine := html.New("./views", ".html")
	// Init controller
	controller := new(controller.Controller)
	controller.Init()
	defer func() {
		controller.Close()
	}()
	// Create a new Fiber template with template engine
	app := fiber.New(fiber.Config{
		Views:        viewEngine,
		ErrorHandler: controller.ErrorController,
	})
	// Default error handler (catch all panic)
	app.Use(recover.New())
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
	app.Get("/reset-cache", controller.GetResetCache)
	app.Get("/reset-db", controller.GetResetDB)
	app.Get("/:url", controller.ValidateController, controller.GetUrlController)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	app.Listen(":" + port)
}
