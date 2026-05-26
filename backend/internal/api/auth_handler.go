package api

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	"toggleflow/internal/auth"
	"toggleflow/internal/db"
)

// SetupStatus returns whether the system has been initialized.
// The frontend calls this on load to decide whether to show setup or login.
func (h *handler) SetupStatus(c *fiber.Ctx) error {
	count, err := h.db.NewSelect().Model((*db.User)(nil)).Count(context.Background())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "db error"})
	}
	return c.JSON(fiber.Map{"initialized": count > 0})
}

type setupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Locale   string `json:"locale"`
}

// Setup creates the first superuser. Fails if any user already exists.
func (h *handler) Setup(c *fiber.Ctx) error {
	count, err := h.db.NewSelect().Model((*db.User)(nil)).Count(context.Background())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "db error"})
	}
	if count > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "already initialized"})
	}

	var req setupRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	if req.Name == "" || req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "name, email and password are required"})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to hash password"})
	}

	locale := req.Locale
	if locale != "en" && locale != "de" {
		locale = "en"
	}

	now := time.Now()
	user := &db.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hash),
		Role:         db.RoleSuperuser,
		Locale:       locale,
		ActivatedAt:  &now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if _, err := h.db.NewInsert().Model(user).Exec(context.Background()); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to create user"})
	}

	token, err := auth.GenerateToken(user.ID, user.Role)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to generate token"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"token": token, "user": user})
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login validates credentials and returns a JWT token.
func (h *handler) Login(c *fiber.Ctx) error {
	var req loginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}

	var user db.User
	err := h.db.NewSelect().Model(&user).Where("email = ?", req.Email).Scan(context.Background())
	if err != nil {
		// Return same error for wrong email and wrong password to avoid user enumeration
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}

	if user.WelcomeToken != "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "account not yet activated — check your welcome link"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}

	token, err := auth.GenerateToken(user.ID, user.Role)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to generate token"})
	}

	return c.JSON(fiber.Map{"token": token, "user": user})
}

type updateProfileRequest struct {
	Name   string `json:"name"`
	Locale string `json:"locale"`
}

// UpdateProfile lets a user update their own name and locale preference.
func (h *handler) UpdateProfile(c *fiber.Ctx) error {
	claims := auth.GetClaims(c)

	var req updateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}

	var user db.User
	if err := h.db.NewSelect().Model(&user).Where("id = ?", claims.UserID).Scan(context.Background()); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Locale == "en" || req.Locale == "de" {
		user.Locale = req.Locale
	}
	user.UpdatedAt = time.Now()

	if _, err := h.db.NewUpdate().Model(&user).WherePK().Exec(context.Background()); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to update profile"})
	}

	return c.JSON(user)
}

type activateRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

// ActivateAccount validates the welcome token and sets the user's password.
// After activation the token is cleared and the user can log in normally.
func (h *handler) ActivateAccount(c *fiber.Ctx) error {
	var req activateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	if req.Token == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "token and password are required"})
	}
	if len(req.Password) < 8 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "password must be at least 8 characters"})
	}

	sum := sha256.Sum256([]byte(req.Token))
	tokenHash := hex.EncodeToString(sum[:])

	var user db.User
	err := h.db.NewSelect().Model(&user).
		Where("welcome_token = ? AND welcome_token_expires_at > ?", tokenHash, time.Now()).
		Scan(context.Background())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid or expired token"})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to hash password"})
	}

	now := time.Now()
	user.PasswordHash = string(hash)
	user.WelcomeToken = ""
	user.WelcomeTokenExpiresAt = nil
	user.ActivatedAt = &now
	user.UpdatedAt = now

	if _, err := h.db.NewUpdate().Model(&user).WherePK().Exec(context.Background()); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to activate account"})
	}

	jwtToken, err := auth.GenerateToken(user.ID, user.Role)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to generate token"})
	}

	return c.JSON(fiber.Map{"token": jwtToken, "user": user})
}

// GetResetInfo returns name and email for a user with a pending reset token.
func (h *handler) GetResetInfo(c *fiber.Ctx) error {
	id := c.Params("uuid")
	var user db.User
	err := h.db.NewSelect().Model(&user).
		Where("uuid = ? AND reset_token != ''", id).
		Scan(context.Background())
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "reset link not found"})
	}
	return c.JSON(fiber.Map{"name": user.Name, "email": user.Email})
}

type resetPasswordRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

// ResetPassword validates the reset token and sets a new password.
func (h *handler) ResetPassword(c *fiber.Ctx) error {
	var req resetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	if req.Token == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "token and password are required"})
	}
	if len(req.Password) < 8 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "password must be at least 8 characters"})
	}

	sum := sha256.Sum256([]byte(req.Token))
	tokenHash := hex.EncodeToString(sum[:])

	var user db.User
	err := h.db.NewSelect().Model(&user).
		Where("reset_token = ? AND reset_token_expires_at > ?", tokenHash, time.Now()).
		Scan(context.Background())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid or expired token"})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to hash password"})
	}

	user.PasswordHash = string(hash)
	user.ResetToken = ""
	user.ResetTokenExpiresAt = nil
	user.UpdatedAt = time.Now()

	if _, err := h.db.NewUpdate().Model(&user).WherePK().Exec(context.Background()); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to reset password"})
	}

	jwtToken, err := auth.GenerateToken(user.ID, user.Role)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to generate token"})
	}

	return c.JSON(fiber.Map{"token": jwtToken, "user": user})
}

// Me returns the currently authenticated user.
func (h *handler) Me(c *fiber.Ctx) error {
	claims := auth.GetClaims(c)

	var user db.User
	if err := h.db.NewSelect().Model(&user).Where("id = ?", claims.UserID).Scan(context.Background()); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}

	return c.JSON(user)
}
