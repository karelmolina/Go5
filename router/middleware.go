package router

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/karelmolina/go5/internal/utils"
)

func AuthMiddleware(c fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return utils.SendError(c, "TOKEN_MISSING", fiber.StatusUnauthorized)
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return utils.SendError(c, "TOKEN_INVALID", fiber.StatusUnauthorized)
	}

	claims, err := utils.ParseToken(parts[1])
	if err != nil {
		if err == jwt.ErrTokenExpired {
			return utils.SendError(c, "TOKEN_EXPIRED", fiber.StatusUnauthorized)
		}
		return utils.SendError(c, "TOKEN_INVALID", fiber.StatusUnauthorized)
	}

	c.Locals("claims", claims)
	return c.Next()
}

func AdminMiddleware(c fiber.Ctx) error {
	claims, ok := c.Locals("claims").(*utils.Claims)
	if !ok || claims == nil {
		return utils.SendError(c, "TOKEN_MISSING", fiber.StatusUnauthorized)
	}

	if claims.Role != "admin" {
		return utils.SendError(c, "FORBIDDEN", fiber.StatusForbidden)
	}

	return c.Next()
}
