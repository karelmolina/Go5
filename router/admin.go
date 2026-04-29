package router

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/karelmolina/play5/database"
	"github.com/karelmolina/play5/internal/utils"
	"github.com/karelmolina/play5/model"
)

type ApproveRequest struct {
	Approved bool `json:"approved" validate:"required"`
}

func ListUsersHandler(c fiber.Ctx) error {
	var users []model.User
	query := database.DB

	isApproved := c.Query("isApproved")
	if isApproved != "" {
		if isApproved == "true" {
			query = query.Where("is_approved = ?", true)
		} else if isApproved == "false" {
			query = query.Where("is_approved = ?", false)
		}
	}

	if err := query.Find(&users).Error; err != nil {
		return utils.SendError(c, "INTERNAL_ERROR", fiber.StatusInternalServerError)
	}

	return c.JSON(users)
}

func ApproveUserHandler(c fiber.Ctx) error {
	var req ApproveRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.SendError(c, "VALIDATION_ERROR", fiber.StatusBadRequest)
	}

	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.SendError(c, "USER_NOT_FOUND", fiber.StatusNotFound)
	}

	var user model.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		return utils.SendError(c, "USER_NOT_FOUND", fiber.StatusNotFound)
	}

	claims := c.Locals("claims").(*utils.Claims)

	if req.Approved {
		now := time.Now()
		user.IsApproved = true
		user.ApprovedAt = &now
		user.ApprovedBy = &claims.Sub
	} else {
		user.IsApproved = false
		user.ApprovedAt = nil
		user.ApprovedBy = nil
	}

	if err := database.DB.Save(&user).Error; err != nil {
		return utils.SendError(c, "INTERNAL_ERROR", fiber.StatusInternalServerError)
	}

	return c.JSON(user)
}
