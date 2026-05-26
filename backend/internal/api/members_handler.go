package api

import (
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"toggleflow/internal/auth"
	"toggleflow/internal/db"
)

type MemberResponse struct {
	UserID   int64     `bun:"user_id"   json:"user_id"`
	Name     string    `bun:"name"      json:"name"`
	Email    string    `bun:"email"     json:"email"`
	Role     db.Role   `bun:"role"      json:"role"`
	JoinedAt time.Time `bun:"joined_at" json:"joined_at"`
}

func (h *handler) ListMembers(c *fiber.Ctx) error {
	pid, err := strconv.ParseInt(c.Params("pid"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid project id"})
	}
	if err := h.checkProjectAccess(c, pid); err != nil {
		return err
	}

	ctx := context.Background()
	members := make([]MemberResponse, 0)
	if err := h.db.NewSelect().
		TableExpr("project_members AS pm").
		ColumnExpr("pm.user_id, u.name, u.email, u.role, pm.created_at AS joined_at").
		Join("JOIN users AS u ON u.id = pm.user_id").
		Where("pm.project_id = ?", pid).
		OrderExpr("pm.created_at ASC").
		Scan(ctx, &members); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch members"})
	}

	return c.JSON(members)
}

type addMemberRequest struct {
	UserID int64 `json:"user_id"`
}

func (h *handler) AddMember(c *fiber.Ctx) error {
	pid, err := strconv.ParseInt(c.Params("pid"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid project id"})
	}
	if err := h.checkProjectAccess(c, pid); err != nil {
		return err
	}

	var req addMemberRequest
	if err := c.BodyParser(&req); err != nil || req.UserID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user_id is required"})
	}

	ctx := context.Background()
	var user db.User
	if err := h.db.NewSelect().Model(&user).Where("id = ?", req.UserID).Scan(ctx); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}

	now := time.Now()
	member := &db.ProjectMember{ProjectID: pid, UserID: req.UserID, CreatedAt: now}
	if _, err := h.db.NewInsert().Model(member).On("CONFLICT DO NOTHING").Exec(ctx); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to add member"})
	}

	return c.Status(fiber.StatusCreated).JSON(MemberResponse{
		UserID: user.ID, Name: user.Name, Email: user.Email, Role: user.Role, JoinedAt: now,
	})
}

func (h *handler) RemoveMember(c *fiber.Ctx) error {
	pid, err := strconv.ParseInt(c.Params("pid"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid project id"})
	}
	uid, err := strconv.ParseInt(c.Params("uid"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	claims := auth.GetClaims(c)
	if uid == claims.UserID && db.RoleRank(claims.Role) < db.RoleRank(db.RoleAdmin) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "you cannot remove yourself from a project"})
	}
	if err := h.checkProjectAccess(c, pid); err != nil {
		return err
	}

	if _, err := h.db.NewDelete().TableExpr("project_members").
		Where("project_id = ? AND user_id = ?", pid, uid).
		Exec(context.Background()); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to remove member"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
