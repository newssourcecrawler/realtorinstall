package services

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/newssourcecrawler/realtorinstall/api/models"
	"github.com/newssourcecrawler/realtorinstall/api/repos"
)

// JWTClaims defines the custom fields stored in our JWT.
type JWTClaims struct {
	UserID   int64  `json:"user_id"`
	TenantID string `json:"tenant_id"`
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

// RegisterUser registers a brand‐new user. It hashes the provided rawPassword before saving.
// Returns ErrUserAlreadyExists if the username is already taken.
func (s *AuthService) RegisterUser(
	ctx context.Context,
	tenantID string,
	currentUser string,
	u models.User,
	rawPassword string,
) (int64, error) {
	// Mandatory fields
	if u.UserName == "" || u.FirstName == "" || u.LastName == "" || u.Role == "" {
		return 0, repos.ErrInvalidRegistration
	}

	// Check if username already exists
	if existing, _ := s.userRepo.GetByUsername(ctx, tenantID, u.UserName); existing != nil {
		return 0, repos.ErrUserAlreadyExists
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return 0, repos.ErrGenerateFromPassword
	}
	u.PasswordHash = string(hashed)

	now := time.Now().UTC()
	u.TenantID = tenantID
	u.CreatedAt = now
	u.LastModified = now
	u.CreatedBy = currentUser
	u.ModifiedBy = currentUser
	u.Deleted = false

	return s.userRepo.Create(ctx, &u)
}

// Login authenticates a user by username+password, and if successful returns a signed JWT string.
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

	// Compare password hashes
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return "", repos.ErrInvalidCredentials
	}

	// Build claims
	now := time.Now().UTC()
	claims := JWTClaims{
		UserID:   user.ID,
		TenantID: user.TenantID,
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
		return "", repos.ErrSignedString
	}
	return signed, nil
}

// ParseToken validates a JWT string, returns its custom claims if valid, or an error.
func (s *AuthService) ParseToken(tokenStr string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})
	if err != nil {
		return nil, repos.ErrParseWithClaims
	}
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, repos.ErrInvalidTokenClaims
	}
	return claims, nil
}
