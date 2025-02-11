package handlers

import (
	"fmt"
	"go_base/internal/models"
	"go_base/internal/services"
	"go_base/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var request models.UserLoginRequest

	if err := c.BodyParser(&request); err != nil {
		return c.JSON(utils.ResponseApi(fiber.StatusBadRequest, "Bad Request", err.Error()))
	}

	validator := utils.NewValidator()
	errs := validator.Validate(request)
	if errs != nil {
		errorMsg := []string{}
		for _, err := range errs {
			errorMsg = append(errorMsg, fmt.Sprintf("%s: %s", err.FailedField, err.Tag))
		}

		return c.JSON(utils.ResponseApi(fiber.StatusBadRequest, "Invalid Input", errorMsg))
	}

	token, err := h.userService.Login(request)
	if err != nil {
		return c.JSON(utils.ResponseApi(fiber.StatusNotFound, "User Not Found", err.Error()))
	}

	return c.JSON(utils.ResponseApi(fiber.StatusOK, "Login Success", token))
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	var request models.UserGetRequest

	if err := c.ParamsParser(&request); err != nil {
		return c.JSON(utils.ResponseApi(fiber.StatusBadRequest, "Bad Request", err.Error()))
	}

	user, err := h.userService.GetUserByID(request)
	if err != nil {
		return c.JSON(utils.ResponseApi(fiber.StatusNotFound, "User Not Found", err.Error()))
	}

	return c.JSON(utils.ResponseApi(fiber.StatusOK, "Get User Success", user))
}
