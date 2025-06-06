package services

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/newssourcecrawler/realtorinstall/api/models"
	"github.com/newssourcecrawler/realtorinstall/api/repos"
)

// JWTClaims holds custom fields plus RegisteredClaims.
type JWTClaims struct {
	UserID      int64    `json:"user_id"`
	TenantID    string   `json:"tenant_id"`
	Permissions []string `json:"perms"`
	jwt.RegisteredClaims
}

type AuthService struct {
	userRepo  repos.UserRepo
	jwtSecret []byte
	ttl       time.Duration
}

func NewAuthService(userRepo repos.UserRepo, jwtSecret string, tokenTTL time.Duration) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: []byte(jwtSecret),
		ttl:       tokenTTL,
	}
}

// Register a new user (bcrypt hash); ErrUserAlreadyExists if username taken.
func (s *AuthService) Register(
	ctx context.Context,
	tenantID string,
	currentUser string,
	u models.User,
	rawPassword string,
) (int64, error) {
	if u.UserName == "" || u.FirstName == "" || u.LastName == "" || u.Role == "" {
		return 0, errors.New("missing required fields")
	}

	if existing, _ := s.userRepo.GetByUsername(ctx, tenantID, u.UserName); existing != nil {
		return 0, repos.ErrUserAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	u.PasswordHash = string(hash)
	now := time.Now().UTC()
	u.TenantID = tenantID
	u.CreatedAt = now
	u.LastModified = now
	u.CreatedBy = currentUser
	u.ModifiedBy = currentUser
	u.Deleted = false

	return s.userRepo.Create(ctx, &u)
}

// Login authenticates username/password, loads permissions, issues JWT.
func (s *AuthService) Login(
	ctx context.Context,
	tenantID string,
	username string,
	password string,
) (string, error) {
	user, err := s.userRepo.GetByUsername(ctx, tenantID, username)
	if err != nil {
		return "", repos.ErrInvalidCredentials
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return "", repos.ErrInvalidCredentials
	}

	// 1) Load this user’s permissions from userRepo
	perms, err := s.userRepo.ListPermissionsForUser(ctx, tenantID, user.ID)
	if err != nil {
		return "", err
	}

	now := time.Now().UTC()
	claims := JWTClaims{
		UserID:      user.ID,
		TenantID:    user.TenantID,
		Permissions: perms, // e.g. ["create_sale","delete_user","view_commissions"]
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.ttl)),
			Issuer:    "realtor-installment-app",
			Subject:   username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}
	return signed, nil
}

// ParseToken verifies a JWT string and returns its claims.
func (s *AuthService) ParseToken(tokenStr string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, repos.ErrInvalidTokenClaims
	}
	return claims, nil
}
