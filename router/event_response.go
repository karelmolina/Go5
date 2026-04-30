package router

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/karelmolina/go5/database"
	"github.com/karelmolina/go5/internal/utils"
	"github.com/karelmolina/go5/model"
)

type RespondRequest struct {
	Status model.ResponseStatus `json:"status" validate:"required,oneof=pending going rejected"`
}

func RespondToEventHandler(c fiber.Ctx) error {
	var req RespondRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.SendError(c, "VALIDATION_ERROR", fiber.StatusBadRequest)
	}

	if err := utils.ValidateStruct(req); err != nil {
		return utils.SendError(c, "VALIDATION_ERROR", fiber.StatusBadRequest)
	}

	eventID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.SendError(c, "EVENT_NOT_FOUND", fiber.StatusNotFound)
	}

	claims := c.Locals("claims").(*utils.Claims)
	userID := claims.Sub

	var event model.Event
	if err := database.DB.First(&event, "id = ?", eventID).Error; err != nil {
		return utils.SendError(c, "EVENT_NOT_FOUND", fiber.StatusNotFound)
	}

	var existingCount int64
	database.DB.Model(&model.EventResponse{}).
		Where("event_id = ? AND user_id = ?", eventID, userID).
		Count(&existingCount)
	if existingCount > 0 {
		return utils.SendError(c, "RSVP_EXISTS", fiber.StatusConflict)
	}

	if req.Status == model.StatusGoing && event.MaxAssistants > 0 {
		var goingCount int64
		database.DB.Model(&model.EventResponse{}).
			Where("event_id = ? AND status = ?", eventID, model.StatusGoing).
			Count(&goingCount)
		if goingCount >= int64(event.MaxAssistants) {
			return utils.SendError(c, "EVENT_FULL", fiber.StatusConflict)
		}
	}

	response := model.EventResponse{
		EventID: eventID,
		UserID:  userID,
		Status:  req.Status,
	}

	if err := database.DB.Create(&response).Error; err != nil {
		return utils.SendError(c, "INTERNAL_ERROR", fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

func UpdateMyResponseHandler(c fiber.Ctx) error {
	var req RespondRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.SendError(c, "VALIDATION_ERROR", fiber.StatusBadRequest)
	}

	if err := utils.ValidateStruct(req); err != nil {
		return utils.SendError(c, "VALIDATION_ERROR", fiber.StatusBadRequest)
	}

	eventID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.SendError(c, "EVENT_NOT_FOUND", fiber.StatusNotFound)
	}

	claims := c.Locals("claims").(*utils.Claims)
	userID := claims.Sub

	var event model.Event
	if err := database.DB.First(&event, "id = ?", eventID).Error; err != nil {
		return utils.SendError(c, "EVENT_NOT_FOUND", fiber.StatusNotFound)
	}

	var response model.EventResponse
	if err := database.DB.Where("event_id = ? AND user_id = ?", eventID, userID).First(&response).Error; err != nil {
		return utils.SendError(c, "RSVP_NOT_FOUND", fiber.StatusNotFound)
	}

	if req.Status == model.StatusGoing && event.MaxAssistants > 0 && response.Status != model.StatusGoing {
		var goingCount int64
		database.DB.Model(&model.EventResponse{}).
			Where("event_id = ? AND status = ?", eventID, model.StatusGoing).
			Count(&goingCount)
		if goingCount >= int64(event.MaxAssistants) {
			return utils.SendError(c, "EVENT_FULL", fiber.StatusConflict)
		}
	}

	response.Status = req.Status
	if err := database.DB.Save(&response).Error; err != nil {
		return utils.SendError(c, "INTERNAL_ERROR", fiber.StatusInternalServerError)
	}

	return c.JSON(response)
}

func ListEventResponsesHandler(c fiber.Ctx) error {
	eventID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.SendError(c, "EVENT_NOT_FOUND", fiber.StatusNotFound)
	}

	var event model.Event
	if err := database.DB.First(&event, "id = ?", eventID).Error; err != nil {
		return utils.SendError(c, "EVENT_NOT_FOUND", fiber.StatusNotFound)
	}

	var responses []model.EventResponse
	if err := database.DB.Where("event_id = ?", eventID).Find(&responses).Error; err != nil {
		return utils.SendError(c, "INTERNAL_ERROR", fiber.StatusInternalServerError)
	}

	return c.JSON(responses)
}

func GetMyResponseHandler(c fiber.Ctx) error {
	eventID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.SendError(c, "EVENT_NOT_FOUND", fiber.StatusNotFound)
	}

	claims := c.Locals("claims").(*utils.Claims)
	userID := claims.Sub

	var event model.Event
	if err := database.DB.First(&event, "id = ?", eventID).Error; err != nil {
		return utils.SendError(c, "EVENT_NOT_FOUND", fiber.StatusNotFound)
	}

	var response model.EventResponse
	if err := database.DB.Where("event_id = ? AND user_id = ?", eventID, userID).First(&response).Error; err != nil {
		return utils.SendError(c, "RSVP_NOT_FOUND", fiber.StatusNotFound)
	}

	return c.JSON(response)
}
