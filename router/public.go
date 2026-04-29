package router

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/karelmolina/play5/database"
	"github.com/karelmolina/play5/internal/utils"
	"github.com/karelmolina/play5/model"
)

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=1,max=30"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func RegisterHandler(c fiber.Ctx) error {
	var req RegisterRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.SendError(c, "VALIDATION_ERROR", fiber.StatusBadRequest)
	}

	if err := utils.ValidateStruct(req); err != nil {
		if verrs, ok := err.(validator.ValidationErrors); ok {
			for _, e := range verrs {
				if e.Field() == "Password" && e.Tag() == "min" {
					return utils.SendError(c, "PASSWORD_TOO_SHORT", fiber.StatusBadRequest)
				}
			}
		}
		return utils.SendError(c, "VALIDATION_ERROR", fiber.StatusBadRequest)
	}

	if utils.IsUsernameTaken(req.Username) {
		return utils.SendError(c, "USERNAME_TAKEN", fiber.StatusConflict)
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return utils.SendError(c, "VALIDATION_ERROR", fiber.StatusBadRequest)
	}

	user := model.User{
		Username:     req.Username,
		PasswordHash: hash,
		Role:         model.RolePlayer,
		IsApproved:   false,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return utils.SendError(c, "INTERNAL_ERROR", fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

func LoginHandler(c fiber.Ctx) error {
	var req LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.SendError(c, "VALIDATION_ERROR", fiber.StatusBadRequest)
	}

	if err := utils.ValidateStruct(req); err != nil {
		return utils.SendError(c, "VALIDATION_ERROR", fiber.StatusBadRequest)
	}

	var user model.User
	if err := database.DB.Where("LOWER(username) = LOWER(?)", req.Username).First(&user).Error; err != nil {
		return utils.SendError(c, "INVALID_CREDENTIALS", fiber.StatusUnauthorized)
	}

	if !utils.CheckPassword(user.PasswordHash, req.Password) {
		return utils.SendError(c, "INVALID_CREDENTIALS", fiber.StatusUnauthorized)
	}

	token, err := utils.GenerateToken(user)
	if err != nil {
		return utils.SendError(c, "VALIDATION_ERROR", fiber.StatusBadRequest)
	}

	return c.JSON(fiber.Map{
		"token": token,
		"user":  user,
	})
}
