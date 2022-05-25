package main

// add engine
import (
	"log"
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
	// unittest.TestTrimTimeStamp()
	// ------------------------------------------------------------
	logFile, _ := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer logFile.Close()
	log.SetOutput(logFile)

	viewEngine := html.New("./views", ".html")
	// Create a new Fiber template with template engine
	app := fiber.New(fiber.Config{
		Views: viewEngine,
	})
	// Init controller + Connect database and cache
	controller := new(controller.Controller)
	controller.Init()
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
	defer func() {
		err := recover()
		if err != nil {
			log.Println(err)
		}
		controller.Close()
	}()
}
