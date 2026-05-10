package services

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) hashPassword(p string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(p), 14)
	return string(b), err
}

func (s *AuthService) checkPassword(hash, p string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(p)) == nil
}

func (s *AuthService) generateToken(
	userID uint64,
	email string,
	roles []string,
	permissions []string,
) (string, error) {

	expiryHours, _ := strconv.Atoi(os.Getenv("JWT_EXPIRY_HOURS"))
	if expiryHours <= 0 {
		expiryHours = 24
	}

	claims := jwt.MapClaims{
		"user_id":     userID,
		"email":       email,
		"roles":       roles,
		"permissions": permissions,
		"exp":         time.Now().Add(time.Duration(expiryHours) * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	return token.SignedString(jwtSecret)
}

func (s *AuthService) Signup(req dtos.SignupRequest) error {
	hash, err := s.hashPassword(req.Password)
	if err != nil {
		return err
	}

	user := models.User{
		Name:              req.Name,
		Email:             req.Email,
		Password:          hash,
		IsVerified:        false,
		VerificationToken: uuid.NewString(),
	}

	return db.DB.Create(&user).Error
}

func (s *AuthService) VerifyAccount(token string) error {
	var user models.User

	if err := db.DB.Where("verification_token = ?", token).First(&user).Error; err != nil {
		return err
	}

	user.IsVerified = true
	user.VerificationToken = ""

	return db.DB.Save(&user).Error
}

func (s *AuthService) Login(req dtos.LoginRequest) (*dtos.LoginResponse, error) {
	var user models.User

	if err := db.DB.
		Preload("Roles.Permissions").
		Preload("Permissions").
		Where("email = ?", req.Email).
		First(&user).Error; err != nil {
		return nil, err
	}

	if !user.IsVerified {
		return nil, gorm.ErrRecordNotFound
	}

	if !s.checkPassword(user.Password, req.Password) {
		return nil, errors.New("invalid email or password")
	}

	roles := make([]string, 0)
	permMap := make(map[string]bool)

	for _, role := range user.Roles {
		roles = append(roles, role.Name)

		for _, perm := range role.Permissions {
			permMap[perm.Name] = true
		}
	}
	for _, perm := range user.Permissions {
		permMap[perm.Name] = true
	}

	permissions := make([]string, 0, len(permMap))
	for p := range permMap {
		permissions = append(permissions, p)
	}

	token, err := s.generateToken(user.ID, user.Email, roles, permissions)

	if err != nil {
		return nil, err
	}

	return &dtos.LoginResponse{
		Token:       token,
		UserID:      user.ID,
		Email:       user.Email,
		Roles:       roles,
		Permissions: permissions,
	}, nil
}

func (s *AuthService) ForgotPassword(email string) error {
	var user models.User

	if err := db.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil // do not leak existence
	}

	token := uuid.NewString()
	exp := time.Now().Add(1 * time.Hour)

	user.ResetToken = token
	user.ResetTokenExpiry = &exp

	return db.DB.Save(&user).Error
}

func (s *AuthService) ResetPassword(token, password string) error {
	var user models.User

	if err := db.DB.Where("reset_token = ?", token).First(&user).Error; err != nil {
		return err
	}

	if user.ResetTokenExpiry == nil || user.ResetTokenExpiry.Before(time.Now()) {
		return errors.New("reset token has expired or is invalid")
	}

	hash, err := s.hashPassword(password)
	if err != nil {
		return err
	}

	user.Password = hash
	user.ResetToken = ""
	user.ResetTokenExpiry = nil

	return db.DB.Save(&user).Error
}

func (s *AuthService) ChangePassword(userID uint64, oldPass, newPass string) error {
	var user models.User

	if err := db.DB.First(&user, userID).Error; err != nil {
		return err
	}

	if !s.checkPassword(user.Password, oldPass) {
		return errors.New("incorrect old password")
	}

	hash, err := s.hashPassword(newPass)
	if err != nil {
		return err
	}
	user.Password = hash

	return db.DB.Save(&user).Error
}
