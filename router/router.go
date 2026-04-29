package router

import "github.com/gofiber/fiber/v3"

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	api.Post("/register", RegisterHandler)
	api.Post("/login", LoginHandler)

	admin := api.Group("/users", AuthMiddleware, AdminMiddleware)
	admin.Get("/", ListUsersHandler)
	admin.Patch("/:id/approve", ApproveUserHandler)
}
