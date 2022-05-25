package main

// add engine
import (
	"log"
	"os"
	"os/signal"
	"server_go/controller"
	"server_go/limiter"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html"
)

func main() {
	logger, _ := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	log.SetOutput(logger)
	log.Println("Server start!")

	viewEngine := html.New("./views", ".html")
	// Init controller
	controller := new(controller.Controller)
	controller.Init()

	// Catch Ctr + C
	go func() {
		signalChannel := make(chan os.Signal)
		signal.Notify(signalChannel, os.Interrupt)
		<-signalChannel
		log.Println("Server end!")
		controller.Close()
		logger.Close()
		os.Exit(0)
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
