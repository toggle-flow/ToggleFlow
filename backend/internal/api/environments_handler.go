package api

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"toggleflow/internal/db"
)

func (h *handler) ListEnvironments(c *fiber.Ctx) error {
	pid, err := strconv.ParseInt(c.Params("pid"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid project id"})
	}
	if err := h.checkProjectAccess(c, pid); err != nil {
		return err
	}

	pq := parsePage(c)
	ctx := context.Background()

	q := h.db.NewSelect().Model((*db.Environment)(nil)).Where("project_id = ?", pid).OrderExpr("created_at ASC")
	if pq.Search != "" {
		q = q.Where("lower(name) LIKE lower(?)", "%"+pq.Search+"%")
	}

	total, err := q.Count(ctx)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to count environments"})
	}

	envs := make([]db.Environment, 0)
	if err := q.Limit(pq.Limit).Offset(pq.Offset).Scan(ctx, &envs); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch environments"})
	}

	return c.JSON(Page[db.Environment]{Data: envs, Total: total, Limit: pq.Limit, Offset: pq.Offset})
}

type createEnvironmentRequest struct {
	Name        string `json:"name"`
	Key         string `json:"key"`
	Description string `json:"description"`
	Protected   bool   `json:"protected"`
}

func (h *handler) CreateEnvironment(c *fiber.Ctx) error {
	pid, err := strconv.ParseInt(c.Params("pid"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid project id"})
	}
	if err := h.checkProjectAccess(c, pid); err != nil {
		return err
	}

	var req createEnvironmentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "name is required"})
	}

	rawSDKKey, keyHash, keyPrefix, err := generateRawKey("sdk_")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to generate sdk key"})
	}

	key := req.Key
	if key == "" {
		key = slugify(req.Name)
	}

	ctx := context.Background()
	env := &db.Environment{
		ProjectID:   pid,
		Name:        req.Name,
		Key:         key,
		Description: req.Description,
		Protected:   req.Protected,
		SDKKey:      rawSDKKey, // kept for schema compat, not returned in API
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if _, err := h.db.NewInsert().Model(env).Exec(ctx); err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "an environment with that name already exists"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "failed to create environment"})
	}

	// Create the initial SDK key entry in sdk_keys table.
	sdkKeyRecord := &db.SDKKey{
		EnvironmentID: env.ID,
		Label:         "Default",
		KeyHash:       keyHash,
		KeyPrefix:     keyPrefix,
		CreatedAt:     time.Now(),
	}
	_, _ = h.db.NewInsert().Model(sdkKeyRecord).Exec(ctx)

	actor := h.actorName(c)
	go h.writeAudit(pid, actor, "env.created", env.Key, "",
		toJSON(map[string]any{"name": env.Name, "key": env.Key}))

	return c.Status(fiber.StatusCreated).JSON(env)
}

type updateEnvironmentRequest struct {
	Name        string `json:"name"`
	Key         string `json:"key"`
	Description string `json:"description"`
	Protected   bool   `json:"protected"`
}

func (h *handler) UpdateEnvironment(c *fiber.Ctx) error {
	pid, err := strconv.ParseInt(c.Params("pid"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid project id"})
	}
	if err := h.checkProjectAccess(c, pid); err != nil {
		return err
	}
	eid, err := strconv.ParseInt(c.Params("eid"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid environment id"})
	}

	ctx := context.Background()
	var env db.Environment
	if err := h.db.NewSelect().Model(&env).Where("id = ?", eid).Scan(ctx); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "environment not found"})
	}

	var req updateEnvironmentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "name is required"})
	}

	oldName, oldKey := env.Name, env.Key
	env.Name = req.Name
	env.Key = req.Key
	env.Description = req.Description
	env.Protected = req.Protected
	env.UpdatedAt = time.Now()

	if _, err := h.db.NewUpdate().Model(&env).Column("name", "key", "description", "protected", "updated_at").Where("id = ?", eid).Exec(ctx); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to update environment"})
	}

	actor := h.actorName(c)
	go h.writeAudit(env.ProjectID, actor, "env.updated", env.Key,
		toJSON(map[string]any{"name": oldName, "key": oldKey}),
		toJSON(map[string]any{"name": env.Name, "key": env.Key}))

	return c.JSON(env)
}

func (h *handler) DeleteEnvironment(c *fiber.Ctx) error {
	pid, err := strconv.ParseInt(c.Params("pid"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid project id"})
	}
	if err := h.checkProjectAccess(c, pid); err != nil {
		return err
	}
	eid, err := strconv.ParseInt(c.Params("eid"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid environment id"})
	}

	ctx := context.Background()

	var env db.Environment
	if err := h.db.NewSelect().Model(&env).Where("id = ?", eid).Scan(ctx); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "environment not found"})
	}

	_, _ = h.db.NewDelete().TableExpr("flag_environments").Where("environment_id = ?", eid).Exec(ctx)

	if _, err := h.db.NewDelete().TableExpr("environments").Where("id = ?", eid).Exec(ctx); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to delete environment"})
	}

	actor := h.actorName(c)
	go h.writeAudit(env.ProjectID, actor, "env.deleted", env.Key,
		toJSON(map[string]any{"name": env.Name, "key": env.Key}), "")

	return c.SendStatus(fiber.StatusNoContent)
}
