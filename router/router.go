package router

import "github.com/gofiber/fiber/v3"

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	api.Post("/register", RegisterHandler)
	api.Post("/login", LoginHandler)

	admin := api.Group("/users", AuthMiddleware, AdminMiddleware)
	admin.Get("/", ListUsersHandler)
	admin.Patch("/:id/approve", ApproveUserHandler)

	// Event routes
	api.Post("/events", AuthMiddleware, AdminMiddleware, CreateEventHandler)
	api.Get("/events", AuthMiddleware, ListEventsHandler)
	api.Get("/events/:id", AuthMiddleware, GetEventHandler)
	api.Patch("/events/:id", AuthMiddleware, AdminMiddleware, UpdateEventHandler)
	api.Delete("/events/:id", AuthMiddleware, AdminMiddleware, DeleteEventHandler)

	// RSVP routes
	api.Post("/events/:id/responses", AuthMiddleware, RespondToEventHandler)
	api.Patch("/events/:id/responses/me", AuthMiddleware, UpdateMyResponseHandler)
	api.Get("/events/:id/responses/me", AuthMiddleware, GetMyResponseHandler)
	api.Get("/events/:id/responses", AuthMiddleware, AdminMiddleware, ListEventResponsesHandler)
}
