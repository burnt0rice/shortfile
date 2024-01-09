package main

import (
	"flag"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"go-shortfile/router"
	"go-shortfile/utils"
	"log"
	"net/http"
	"strconv"
)

func main() {
	utils.InitLogger()

	app := fiber.New()
	app.Use(cors.New())

	var port int
	flag.IntVar(&port, "port", 3000, "The port to listen on")
	flag.Parse()

	//Api Endpoints
	router.AddRouters(app)

	app.Use("/", filesystem.New(filesystem.Config{
		Root: http.Dir("./frontend/dist"),
	}))

	utils.Log.Info("Starting fiber on port: " + strconv.Itoa(port))
	err := app.Listen(":" + strconv.Itoa(port))
	if err != nil {
		log.Fatalln(err.Error())
	}
}
