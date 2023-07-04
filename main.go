package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ihyamars03/jadwal-api/initialize"
	"github.com/ihyamars03/jadwal-api/routes"
)

func init() {
	//initialize.LoadEnv()
	initialize.ConnectDB()
}

func main() {
	app := fiber.New()

	routes.ScheduleRoutes(app)
	routes.UserRoutes(app)

	app.Listen(":3030")
}
