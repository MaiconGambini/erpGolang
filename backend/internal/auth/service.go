package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/MaiconGambini/erpGolang/backend/gen/db"
	"github.com/MaiconGambini/erpGolang/backend/internal/config"
	"github.com/MaiconGambini/erpGolang/backend/internal/platform/jwt"
	"github.com/MaiconGambini/erpGolang/backend/internal/shared/pgutil"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	queries    *db.Queries
	jwt        jwt.Service
	refreshTTL time.Duration
}

func NewService(pool *pgxpool.Pool, cfg config.Config) *Service {
	return &Service{
		queries:    db.New(pool),
		jwt:        jwt.NewService(cfg.JWTAccessSecret, cfg.JWTAccessTTL),
		refreshTTL: cfg.JWTRefreshTTL,
	}
}

type LoginInput struct {
	TenantSlug string
	Email      string
	Password   string
}

type LoginResult struct {
	AccessToken  string
	RefreshToken string
	UserID       string
	TenantID     string
	Role         string
	Name         string
	Email        string
}

type UserProfile struct {
	ID       string `json:"id"`
	TenantID string `json:"tenantId"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}

func (s *Service) Login(ctx context.Context, in LoginInput, bcryptCost int) (LoginResult, error) {
	tenant, err := s.queries.GetTenantBySlug(ctx, in.TenantSlug)
	if err != nil {
		return LoginResult{}, fmt.Errorf("invalid credentials")
	}
	user, err := s.queries.GetUserByEmail(ctx, in.Email)
	if err != nil {
		return LoginResult{}, fmt.Errorf("invalid credentials")
	}
	tenantID, _ := pgutil.PgToUUID(tenant.ID)
	userTenantID, _ := pgutil.PgToUUID(user.TenantID)
	if tenantID != userTenantID {
		return LoginResult{}, fmt.Errorf("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password)); err != nil {
		return LoginResult{}, fmt.Errorf("invalid credentials")
	}
	userID, _ := pgutil.PgToUUID(user.ID)
	access, err := s.jwt.GenerateAccessToken(userID.String(), tenantID.String(), user.Role)
	if err != nil {
		return LoginResult{}, err
	}
	refresh, hash, err := newRefreshToken()
	if err != nil {
		return LoginResult{}, err
	}
	_, err = s.queries.CreateAuthSession(ctx, db.CreateAuthSessionParams{
		TenantID:         pgutil.UUIDToPg(tenantID),
		UserID:           pgutil.UUIDToPg(userID),
		RefreshTokenHash: hash,
		ExpiresAt:        pgutil.Timestamptz(time.Now().Add(s.refreshTTL)),
	})
	if err != nil {
		return LoginResult{}, err
	}
	return LoginResult{
		AccessToken:  access,
		RefreshToken: refresh,
		UserID:       userID.String(),
		TenantID:     tenantID.String(),
		Role:         user.Role,
		Name:         user.Name,
		Email:        user.Email,
	}, nil
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (LoginResult, error) {
	hash := hashToken(refreshToken)
	session, err := s.queries.GetAuthSessionByHash(ctx, hash)
	if err != nil {
		return LoginResult{}, fmt.Errorf("invalid refresh token")
	}
	userID, _ := pgutil.PgToUUID(session.UserID)
	tenantID, _ := pgutil.PgToUUID(session.TenantID)
	user, err := s.queries.GetUserByID(ctx, db.GetUserByIDParams{
		ID:       session.UserID,
		TenantID: session.TenantID,
	})
	if err != nil {
		return LoginResult{}, fmt.Errorf("invalid refresh token")
	}
	_ = s.queries.RevokeAuthSession(ctx, session.ID)
	access, err := s.jwt.GenerateAccessToken(userID.String(), tenantID.String(), user.Role)
	if err != nil {
		return LoginResult{}, err
	}
	newRefresh, newHash, err := newRefreshToken()
	if err != nil {
		return LoginResult{}, err
	}
	_, err = s.queries.CreateAuthSession(ctx, db.CreateAuthSessionParams{
		TenantID:         session.TenantID,
		UserID:           session.UserID,
		RefreshTokenHash: newHash,
		ExpiresAt:        pgutil.Timestamptz(time.Now().Add(s.refreshTTL)),
	})
	if err != nil {
		return LoginResult{}, err
	}
	return LoginResult{
		AccessToken:  access,
		RefreshToken: newRefresh,
		UserID:       userID.String(),
		TenantID:     tenantID.String(),
		Role:         user.Role,
		Name:         user.Name,
		Email:        user.Email,
	}, nil
}

func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return nil
	}
	hash := hashToken(refreshToken)
	session, err := s.queries.GetAuthSessionByHash(ctx, hash)
	if err != nil {
		return nil
	}
	return s.queries.RevokeAuthSession(ctx, session.ID)
}

func (s *Service) Me(ctx context.Context, userID, tenantID string) (UserProfile, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return UserProfile{}, err
	}
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return UserProfile{}, err
	}
	user, err := s.queries.GetUserByID(ctx, db.GetUserByIDParams{
		ID:       pgutil.UUIDToPg(uid),
		TenantID: pgutil.UUIDToPg(tid),
	})
	if err != nil {
		return UserProfile{}, err
	}
	return UserProfile{
		ID:       userID,
		TenantID: tenantID,
		Email:    user.Email,
		Name:     user.Name,
		Role:     user.Role,
	}, nil
}

func (s *Service) JWT() jwt.Service {
	return s.jwt
}

func HashPassword(password string, cost int) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(b), err
}

func newRefreshToken() (token string, hash string, err error) {
	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return "", "", err
	}
	token = hex.EncodeToString(b)
	return token, hashToken(token), nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
