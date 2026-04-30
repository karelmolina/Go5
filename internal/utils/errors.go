package utils

import (
	"github.com/gofiber/fiber/v3"
)

func SendError(c fiber.Ctx, code string, status int) error {
	return c.Status(status).JSON(fiber.Map{
		"error": fiber.Map{
			"code":    code,
			"message": errorMessage(code, c.Get("Accept-Language")),
		},
	})
}

func errorMessage(code, lang string) string {
	messages := map[string]map[string]string{
		"USERNAME_TAKEN": {
			"es": "El nombre de usuario ya está en uso",
			"en": "Username is already taken",
		},
		"PASSWORD_TOO_SHORT": {
			"es": "La contraseña debe tener al menos 6 caracteres",
			"en": "Password must be at least 6 characters",
		},
		"VALIDATION_ERROR": {
			"es": "Error de validación",
			"en": "Validation error",
		},
		"INVALID_CREDENTIALS": {
			"es": "Credenciales inválidas",
			"en": "Invalid credentials",
		},
		"TOKEN_MISSING": {
			"es": "Token de autenticación faltante",
			"en": "Authentication token is missing",
		},
		"TOKEN_INVALID": {
			"es": "Token de autenticación inválido",
			"en": "Authentication token is invalid",
		},
		"TOKEN_EXPIRED": {
			"es": "Token de autenticación expirado",
			"en": "Authentication token has expired",
		},
		"FORBIDDEN": {
			"es": "No tienes permiso para realizar esta acción",
			"en": "You do not have permission to perform this action",
		},
		"USER_NOT_FOUND": {
			"es": "Usuario no encontrado",
			"en": "User not found",
		},
		"EVENT_NOT_FOUND": {
			"es": "Evento no encontrado",
			"en": "Event not found",
		},
		"EVENT_FULL": {
			"es": "El evento ha alcanzado el número máximo de asistentes",
			"en": "Event has reached the maximum number of attendees",
		},
		"RSVP_EXISTS": {
			"es": "Ya has respondido a este evento",
			"en": "You have already responded to this event",
		},
		"RSVP_NOT_FOUND": {
			"es": "No se encontró tu respuesta para este evento",
			"en": "Your response for this event was not found",
		},
	}

	if m, ok := messages[code]; ok {
		if lang == "es" {
			return m["es"]
		}
		return m["en"]
	}
	return code
}
