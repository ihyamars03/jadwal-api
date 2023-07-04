package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ihyamars03/jadwal-api/initialize"
	"github.com/ihyamars03/jadwal-api/models"
	"github.com/ihyamars03/jadwal-api/utils"
)

func UserCheckin(c *fiber.Ctx) error {

	var user models.Users
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	if user.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "Bad Request",
			"message": "Email is required",
		})
	}
	if !utils.EmailRegex.Match([]byte(user.Email)) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "Bad Request",
			"message": "Invalid email",
		})
	}
	result := initialize.DB.Where("email = ?", user.Email).First(&user)
	if result.Error != nil {
		initialize.DB.Create(&user)

	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Success",
		"data":    user,
	})
}
