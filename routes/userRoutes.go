package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ihyamars03/jadwal-api/controller"
)

func UserRoutes(app *fiber.App) {
	app.Post("/checkin", controller.UserCheckin)
}
