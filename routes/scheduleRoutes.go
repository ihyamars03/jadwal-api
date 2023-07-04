package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ihyamars03/jadwal-api/controller"
)

func ScheduleRoutes(app *fiber.App) {
	app.Get("/schedule", controller.GetSchedule)
	app.Post("/schedule", controller.CreateSchedule)
	app.Patch("/schedule", controller.Updateschedule)
	app.Delete("/schedule", controller.Deleteschedule)
}
