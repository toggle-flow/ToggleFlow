package api

import (
	"toggleflow/internal/auth"
	"toggleflow/internal/db"

	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
)

func Register(app *fiber.App, database *bun.DB) {
	h := newHandler(database)

	// Public — no auth required
	app.Get("/api/setup/status", h.SetupStatus)
	app.Post("/api/setup", h.Setup)
	app.Post("/api/auth/login", h.Login)
	app.Post("/api/auth/activate", h.ActivateAccount)
	app.Get("/api/auth/invite/:uuid", h.GetInviteInfo)
	app.Post("/api/auth/reset", h.ResetPassword)
	app.Get("/api/auth/reset/:uuid", h.GetResetInfo)

	// Authenticated — all roles
	protected := app.Group("/api", auth.Require)
	protected.Get("/auth/me", h.Me)
	protected.Patch("/auth/profile", h.UpdateProfile)

	// User management — Admin and above
	users := protected.Group("/users", auth.RequireRole(db.RoleAdmin))
	users.Get("/", h.ListUsers)
	users.Post("/", h.CreateUser)
	users.Patch("/:id", h.UpdateUser)
	users.Post("/:id/reinvite", h.ReinviteUser)

	// Delete user and reset link — Superuser only
	protected.Delete("/users/:id", auth.RequireRole(db.RoleSuperuser), h.DeleteUser)
	protected.Post("/users/:id/reset-link", auth.RequireRole(db.RoleSuperuser), h.GenerateResetLink)

	// Projects — Owner and above
	projects := protected.Group("/projects", auth.RequireRole(db.RoleOwner))
	projects.Post("/", h.CreateProject)
	projects.Get("/", h.ListProjects)
	projects.Patch("/:id", h.UpdateProject)
	projects.Delete("/:id", h.DeleteProject)
	projects.Post("/:pid/environments", h.CreateEnvironment)
	projects.Get("/:pid/environments", h.ListEnvironments)
	projects.Patch("/:pid/environments/:eid", h.UpdateEnvironment)
	projects.Delete("/:pid/environments/:eid", h.DeleteEnvironment)

	// Flags — Editor and above to write, Viewer to read
	protected.Get("/projects/:pid/flags", h.ListFlags)
	protected.Get("/projects/:pid/flags/:key", h.GetFlag)
	protected.Get("/projects/:pid/audit", h.ListAudit)

	flagsWrite := protected.Group("/projects/:pid/flags", auth.RequireRole(db.RoleEditor))
	flagsWrite.Post("/", h.CreateFlag)
	flagsWrite.Patch("/:key", h.UpdateFlag)
	flagsWrite.Patch("/:key/env", h.ToggleFlagEnv)
	flagsWrite.Delete("/:key", h.DeleteFlag)

	// SDK — authenticated by sdk_key query param, not JWT
	sdk := app.Group("/sdk")
	sdk.Get("/flags", h.SDKGetFlags)
	sdk.Post("/evaluate", h.SDKEvaluate)
	sdk.Get("/stream", h.SDKStream)
}
