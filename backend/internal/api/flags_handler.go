package api

import (
	"context"
	"strconv"
	"time"

	"github.com/RXNova/ToggleFlow/internal/db"
	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
)

type FlagEnvState struct {
	EnvironmentID   int64  `json:"environment_id"`
	EnvironmentName string `json:"environment_name"`
	EnvironmentSlug string `json:"environment_slug"`
	Enabled         bool   `json:"enabled"`
}

type FlagResponse struct {
	db.Flag
	Environments []FlagEnvState `json:"environments"`
}

func (h *handler) ListFlags(c *fiber.Ctx) error {
	pid, err := strconv.ParseInt(c.Params("pid"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid project id"})
	}
	ctx := context.Background()

	var flags []db.Flag
	if err := h.db.NewSelect().Model(&flags).Where("project_id = ?", pid).OrderExpr("created_at ASC").Scan(ctx); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch flags"})
	}

	var envs []db.Environment
	if err := h.db.NewSelect().Model(&envs).Where("project_id = ?", pid).OrderExpr("created_at ASC").Scan(ctx); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch environments"})
	}

	result := make([]FlagResponse, len(flags))

	if len(flags) > 0 {
		flagIDs := make([]int64, len(flags))
		for i, f := range flags {
			flagIDs[i] = f.ID
		}

		var flagEnvs []db.FlagEnvironment
		h.db.NewSelect().Model(&flagEnvs).Where("flag_id IN (?)", bun.In(flagIDs)).Scan(ctx)

		type feKey struct{ flagID, envID int64 }
		feMap := make(map[feKey]bool)
		for _, fe := range flagEnvs {
			feMap[feKey{fe.FlagID, fe.EnvironmentID}] = fe.Enabled
		}

		for i, flag := range flags {
			states := make([]FlagEnvState, len(envs))
			for j, env := range envs {
				states[j] = FlagEnvState{
					EnvironmentID:   env.ID,
					EnvironmentName: env.Name,
					EnvironmentSlug: env.Slug,
					Enabled:         feMap[feKey{flag.ID, env.ID}],
				}
			}
			result[i] = FlagResponse{Flag: flag, Environments: states}
		}
	}

	return c.JSON(result)
}

type createFlagRequest struct {
	Name        string `json:"name"`
	Key         string `json:"key"`
	Description string `json:"description"`
}

func (h *handler) CreateFlag(c *fiber.Ctx) error {
	pid, err := strconv.ParseInt(c.Params("pid"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid project id"})
	}

	var req createFlagRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "name is required"})
	}

	key := req.Key
	if key == "" {
		key = slugify(req.Name)
	}

	ctx := context.Background()

	var existing db.Flag
	if err := h.db.NewSelect().Model(&existing).Where("project_id = ? AND key = ?", pid, key).Scan(ctx); err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "a flag with that key already exists"})
	}

	flag := &db.Flag{
		ProjectID:   pid,
		Key:         key,
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if _, err := h.db.NewInsert().Model(flag).Exec(ctx); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to create flag"})
	}

	// Seed flag_environment rows for every existing environment
	var envs []db.Environment
	if err := h.db.NewSelect().Model(&envs).Where("project_id = ?", pid).Scan(ctx); err == nil {
		for _, env := range envs {
			fe := &db.FlagEnvironment{FlagID: flag.ID, EnvironmentID: env.ID}
			h.db.NewInsert().Model(fe).Exec(ctx)
		}
	}

	return c.Status(fiber.StatusCreated).JSON(flag)
}

type toggleFlagRequest struct {
	EnvironmentID int64 `json:"environment_id"`
	Enabled       bool  `json:"enabled"`
}

func (h *handler) UpdateFlag(c *fiber.Ctx) error {
	pid, err := strconv.ParseInt(c.Params("pid"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid project id"})
	}

	ctx := context.Background()
	var flag db.Flag
	if err := h.db.NewSelect().Model(&flag).Where("project_id = ? AND key = ?", pid, c.Params("key")).Scan(ctx); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "flag not found"})
	}

	var req toggleFlagRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}

	var fe db.FlagEnvironment
	if err := h.db.NewSelect().Model(&fe).Where("flag_id = ? AND environment_id = ?", flag.ID, req.EnvironmentID).Scan(ctx); err != nil {
		fe = db.FlagEnvironment{FlagID: flag.ID, EnvironmentID: req.EnvironmentID, Enabled: req.Enabled}
		h.db.NewInsert().Model(&fe).Exec(ctx)
	} else {
		fe.Enabled = req.Enabled
		h.db.NewUpdate().Model(&fe).Column("enabled").Where("id = ?", fe.ID).Exec(ctx)
	}

	return c.JSON(fiber.Map{"ok": true})
}

func (h *handler) GetFlag(c *fiber.Ctx) error {
	pid, err := strconv.ParseInt(c.Params("pid"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid project id"})
	}

	var flag db.Flag
	if err := h.db.NewSelect().Model(&flag).Where("project_id = ? AND key = ?", pid, c.Params("key")).Scan(context.Background()); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "flag not found"})
	}
	return c.JSON(flag)
}

func (h *handler) DeleteFlag(c *fiber.Ctx) error {
	pid, err := strconv.ParseInt(c.Params("pid"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid project id"})
	}

	ctx := context.Background()
	var flag db.Flag
	if err := h.db.NewSelect().Model(&flag).Where("project_id = ? AND key = ?", pid, c.Params("key")).Scan(ctx); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "flag not found"})
	}

	h.db.NewDelete().TableExpr("flag_environments").Where("flag_id = ?", flag.ID).Exec(ctx)

	if _, err := h.db.NewDelete().Model(&flag).Where("id = ?", flag.ID).Exec(ctx); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to delete flag"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
