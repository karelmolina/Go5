package router

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/karelmolina/go5/database"
	"github.com/karelmolina/go5/internal/utils"
	"github.com/karelmolina/go5/model"
)

type CreateEventRequest struct {
	Location      string    `json:"location" validate:"required,max=200"`
	Date          time.Time `json:"date" validate:"required"`
	Time          string    `json:"time" validate:"required,max=10"`
	Description   string    `json:"description" validate:"max=1000"`
	MaxAssistants int       `json:"maxAssistants" validate:"gte=0"`
}

type UpdateEventRequest struct {
	Location      *string    `json:"location,omitempty" validate:"omitempty,max=200"`
	Date          *time.Time `json:"date,omitempty" validate:"omitempty"`
	Time          *string    `json:"time,omitempty" validate:"omitempty,max=10"`
	Description   *string    `json:"description,omitempty" validate:"omitempty,max=1000"`
	MaxAssistants *int       `json:"maxAssistants,omitempty" validate:"omitempty,gte=0"`
}

type EventWithCounts struct {
	model.Event
	GoingCount    int64 `json:"goingCount"`
	PendingCount  int64 `json:"pendingCount"`
	RejectedCount int64 `json:"rejectedCount"`
}

func CreateEventHandler(c fiber.Ctx) error {
	var req CreateEventRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.SendError(c, "VALIDATION_ERROR", fiber.StatusBadRequest)
	}

	if err := utils.ValidateStruct(req); err != nil {
		return utils.SendError(c, "VALIDATION_ERROR", fiber.StatusBadRequest)
	}

	event := model.Event{
		Location:      req.Location,
		Date:          req.Date,
		Time:          req.Time,
		Description:   req.Description,
		MaxAssistants: req.MaxAssistants,
	}

	if err := database.DB.Create(&event).Error; err != nil {
		return utils.SendError(c, "INTERNAL_ERROR", fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusCreated).JSON(event)
}

func ListEventsHandler(c fiber.Ctx) error {
	var events []model.Event
	if err := database.DB.Order("date ASC").Find(&events).Error; err != nil {
		return utils.SendError(c, "INTERNAL_ERROR", fiber.StatusInternalServerError)
	}

	result := make([]EventWithCounts, 0, len(events))
	for _, event := range events {
		counts := getEventResponseCounts(event.ID)
		result = append(result, EventWithCounts{
			Event:         event,
			GoingCount:    counts.Going,
			PendingCount:  counts.Pending,
			RejectedCount: counts.Rejected,
		})
	}

	return c.JSON(result)
}

func GetEventHandler(c fiber.Ctx) error {
	eventID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.SendError(c, "EVENT_NOT_FOUND", fiber.StatusNotFound)
	}

	var event model.Event
	if err := database.DB.First(&event, "id = ?", eventID).Error; err != nil {
		return utils.SendError(c, "EVENT_NOT_FOUND", fiber.StatusNotFound)
	}

	counts := getEventResponseCounts(event.ID)
	return c.JSON(EventWithCounts{
		Event:         event,
		GoingCount:    counts.Going,
		PendingCount:  counts.Pending,
		RejectedCount: counts.Rejected,
	})
}

func UpdateEventHandler(c fiber.Ctx) error {
	var req UpdateEventRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.SendError(c, "VALIDATION_ERROR", fiber.StatusBadRequest)
	}

	if err := utils.ValidateStruct(req); err != nil {
		if verrs, ok := err.(validator.ValidationErrors); ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    "VALIDATION_ERROR",
					"message": utils.ValidationErrorsToMap(verrs),
				},
			})
		}
		return utils.SendError(c, "VALIDATION_ERROR", fiber.StatusBadRequest)
	}

	eventID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.SendError(c, "EVENT_NOT_FOUND", fiber.StatusNotFound)
	}

	var event model.Event
	if err := database.DB.First(&event, "id = ?", eventID).Error; err != nil {
		return utils.SendError(c, "EVENT_NOT_FOUND", fiber.StatusNotFound)
	}

	if req.Location != nil {
		event.Location = *req.Location
	}
	if req.Date != nil {
		event.Date = *req.Date
	}
	if req.Time != nil {
		event.Time = *req.Time
	}
	if req.Description != nil {
		event.Description = *req.Description
	}
	if req.MaxAssistants != nil {
		event.MaxAssistants = *req.MaxAssistants
	}

	if err := database.DB.Save(&event).Error; err != nil {
		return utils.SendError(c, "INTERNAL_ERROR", fiber.StatusInternalServerError)
	}

	counts := getEventResponseCounts(event.ID)
	return c.JSON(EventWithCounts{
		Event:         event,
		GoingCount:    counts.Going,
		PendingCount:  counts.Pending,
		RejectedCount: counts.Rejected,
	})
}

func DeleteEventHandler(c fiber.Ctx) error {
	eventID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.SendError(c, "EVENT_NOT_FOUND", fiber.StatusNotFound)
	}

	var event model.Event
	if err := database.DB.First(&event, "id = ?", eventID).Error; err != nil {
		return utils.SendError(c, "EVENT_NOT_FOUND", fiber.StatusNotFound)
	}

	if err := database.DB.Delete(&event).Error; err != nil {
		return utils.SendError(c, "INTERNAL_ERROR", fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

type responseCounts struct {
	Going    int64
	Pending  int64
	Rejected int64
}

func getEventResponseCounts(eventID uuid.UUID) responseCounts {
	var counts responseCounts
	database.DB.Model(&model.EventResponse{}).
		Where("event_id = ?", eventID).
		Select(
			"COALESCE(SUM(CASE WHEN status = ? THEN 1 ELSE 0 END), 0) as going, "+
			"COALESCE(SUM(CASE WHEN status = ? THEN 1 ELSE 0 END), 0) as pending, "+
			"COALESCE(SUM(CASE WHEN status = ? THEN 1 ELSE 0 END), 0) as rejected",
			model.StatusGoing, model.StatusPending, model.StatusRejected,
		).
		Scan(&counts)
	return counts
}
