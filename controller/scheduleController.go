package controller

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/ihyamars03/jadwal-api/initialize"
	"github.com/ihyamars03/jadwal-api/models"
	"github.com/ihyamars03/jadwal-api/utils"
	"gorm.io/gorm"
)

type ScheduleResponse struct {
	Status  string                       `json:"status"`
	Message string                       `json:"message"`
	Data    map[string][]models.Schedule `json:"data"`
}

type SuccessResponseStruct struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type UpdateBody struct {
	Title string
	Day   string
}

func SuccessResponse(Data interface{}) SuccessResponseStruct {
	return SuccessResponseStruct{
		Status:  "Success",
		Message: "Success",
		Data:    Data,
	}
}

func GetSchedule(c *fiber.Ctx) error {
	query := c.Queries()

	var schedules []models.Schedule
	var user models.Users

	if query["email"] == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "Bad Request",
			"message": "Email is required",
		})
	}

	if !utils.EmailRegex.Match([]byte(query["email"])) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "Bad Request",
			"message": "Invalid email",
		})
	}

	result := initialize.DB.Where("email = ?", query["email"]).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "Not Found",
				"message": "Email is not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "Internal Server Error",
			"message": "Failed to get user",
		})
	}

	if query["day"] == "" {
		schedules = getSchedulesFromDB(user.Id)
		Data := map[string][]models.Schedule{
			"monday":    filterSchedulesByDay(schedules, "monday"),
			"tuesday":   filterSchedulesByDay(schedules, "tuesday"),
			"wednesday": filterSchedulesByDay(schedules, "wednesday"),
			"thursday":  filterSchedulesByDay(schedules, "thursday"),
			"friday":    filterSchedulesByDay(schedules, "friday"),
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "Success",
			"message": "Success",
			"data":    Data,
		})

	} else {
		validDays := []string{"monday", "tuesday", "wednesday", "thursday", "friday"}
		if !contains(validDays, query["day"]) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "Bad Request",
				"message": "Day is invalid",
			})
		}

		resultSchedule := initialize.DB.Where("user_id = ? AND day = ?", user.Id, query["day"]).Find(&schedules)
		if resultSchedule.Error != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "Not Found",
				"message": "Schedule not found",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "Success",
			"message": "Success",
			"data":    schedules,
		})
	}
}

func CreateSchedule(c *fiber.Ctx) error {
	query := c.Queries()

	var user models.Users
	var schedule models.Schedule

	if query["email"] == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "Bad Request",
			"message": "Email is required",
		})
	}

	if !utils.EmailRegex.Match([]byte(query["email"])) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "Bad Request",
			"message": "Invalid email",
		})
	}

	result := initialize.DB.Where("email = ?", query["email"]).First(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "Not Found",
			"message": "Email is not found",
		})
	}

	if err := c.BodyParser(&schedule); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "Bad Request",
			"message": "Invalid request body",
		})
	}

	if schedule.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "Bad Request",
			"message": "Title is required",
		})
	}

	if schedule.Day == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "Bad Request",
			"message": "Day is required",
		})
	}

	if !utils.IsValidDay(schedule.Day) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "Bad Request",
			"message": "Day is invalid",
		})
	}

	newSchedule := models.Schedule{
		UserId: user.Id,
		Title:  schedule.Title,
		Day:    schedule.Day,
	}

	// Insert the new schedule into the database
	result = initialize.DB.Create(&newSchedule)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "Internal Server Error",
			"message": "Failed to create schedule",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "Success",
		"message": "Success",
		"data":    newSchedule,
	})
}

func Updateschedule(c *fiber.Ctx) error {
	query := c.Queries()

	var schedule models.Schedule
	var user models.Users
	var body UpdateBody

	if query["email"] == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "Bad Request",
			"message": "Email is required",
		})
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if !utils.EmailRegex.Match([]byte(query["email"])) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "Bad Request",
			"message": "Invalid email",
		})
	}

	if query["id"] == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "Bad Request",
			"message": "ID is required",
		})
	}

	id, err := strconv.ParseUint(query["id"], 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	result := initialize.DB.Where("email = ?", query["email"]).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "Not Found",
				"message": "Email is not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "Internal Server Error",
			"message": "Failed to get user",
		})
	}

	if body.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "Bad Request",
			"message": "Title is required",
		})
	}

	if schedule.UserId != user.Id {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "Forbidden",
			"message": "Access denied!",
		})
	}

	result = initialize.DB.Where("id = ? AND user_id = ?", id, user.Id).First(&schedule)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "Not Found",
				"message": "Schedule with ID " + query["id"] + " Not Found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "Internal Server Error",
			"message": "Failed to get schedule",
		})
	}

	result = initialize.DB.Model(&schedule).Updates(body)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "Internal Server Error",
			"message": "Failed to update schedule",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "Success",
		"message": "Success",
		"data":    schedule,
	})
}

func Deleteschedule(c *fiber.Ctx) error {
	query := c.Queries()

	var schedule models.Schedule
	var user models.Users

	if query["email"] == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "Bad Request",
			"message": "Email is required",
		})
	}

	if !utils.EmailRegex.Match([]byte(query["email"])) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "Bad Request",
			"message": "Invalid email",
		})
	}

	if query["id"] == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "Bad Request",
			"message": "Id is required",
		})
	}

	id, err := strconv.ParseUint(query["id"], 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	result := initialize.DB.Where("email = ?", query["email"]).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "Not Found",
				"message": "Email is not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "Internal Server Error",
			"message": "Failed to get user",
		})
	}

	if schedule.UserId != user.Id {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "Forbidden",
			"message": "Access denied!",
		})
	}

	result = initialize.DB.Where("id = ? AND user_id = ?", id, user.Id).First(&schedule)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "Not Found",
				"message": "Schedule with ID " + query["id"] + " Not Found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "Internal Server Error",
			"message": "Failed to get schedule",
		})
	}

	result = initialize.DB.Delete(&schedule)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "Internal Server Error",
			"message": "Failed to delete schedule",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "Success",
		"message": "Success",
		"data":    fiber.Map{},
	})
}

func getSchedulesFromDB(userID uint) []models.Schedule {
	var schedules []models.Schedule
	result := initialize.DB.Where("user_id = ?", userID).Find(&schedules)
	if result.Error != nil {
		return []models.Schedule{}
	}
	return schedules
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func filterSchedulesByDay(schedules []models.Schedule, day string) []models.Schedule {
	filteredSchedules := make([]models.Schedule, 0)
	for _, schedule := range schedules {
		if schedule.Day == day {
			filteredSchedules = append(filteredSchedules, schedule)
		}
	}
	return filteredSchedules
}
