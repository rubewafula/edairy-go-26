package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/initializers"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/routes"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const defaultEmail = "superuser@edairy.local"

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "seed superuser:", err)
		os.Exit(1)
	}
}

func run() error {
	_ = godotenv.Load()
	if os.Getenv("DB_NAME") != "edairy" || os.Getenv("DB_HOST") != "localhost" {
		return errors.New("refusing to seed: expected the local edairy database (DB_HOST=localhost, DB_NAME=edairy)")
	}

	email := strings.ToLower(strings.TrimSpace(os.Getenv("SUPERUSER_EMAIL")))
	if email == "" {
		email = defaultEmail
	}
	if !strings.Contains(email, "@") {
		return errors.New("SUPERUSER_EMAIL must be a valid email address")
	}

	db.ConnectToDatabase()
	userMigrator := db.DB.Migrator()
	for _, column := range []string{"IsVerified", "VerificationToken", "ResetToken", "ResetTokenExpiry"} {
		if !userMigrator.HasColumn(&models.User{}, column) {
			if err := userMigrator.AddColumn(&models.User{}, column); err != nil {
				return fmt.Errorf("add users.%s: %w", column, err)
			}
		}
	}
	if err := db.DB.AutoMigrate(
		&models.Role{},
		&models.Permission{},
		&models.UserRole{},
		&models.UserPermission{},
		&models.RolePermission{},
	); err != nil {
		return fmt.Errorf("migrate auth tables: %w", err)
	}

	var existing models.User
	existingAccount := false
	if err := db.DB.Where("email = ?", email).First(&existing).Error; err == nil {
		if os.Getenv("SUPERUSER_RESET") != "true" || existing.Name != "eDairy Super Admin" {
			return fmt.Errorf("account %q already exists; refusing to reset its password", email)
		}
		existingAccount = true
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("check account: %w", err)
	}

	password, err := generatePassword()
	if err != nil {
		return err
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
	permissions, err := protectedRoutePermissions()
	if err != nil {
		return err
	}
	if len(permissions) == 0 {
		return errors.New("no protected API route permissions were found")
	}

	user := models.User{
		Name:       "eDairy Super Admin",
		Email:      email,
		Password:   string(passwordHash),
		IsVerified: true,
	}
	role := models.Role{Name: "Super Admin", GuardName: "web"}
	gin.SetMode(gin.ReleaseMode)
	initializers.InitTimezone()

	if err := db.DB.Transaction(func(tx *gorm.DB) error {
		if existingAccount {
			existing.Password = user.Password
			existing.IsVerified = true
			if err := tx.Save(&existing).Error; err != nil {
				return err
			}
			user = existing
		} else if err := tx.Create(&user).Error; err != nil {
			return err
		}
		if err := tx.Where("name = ?", role.Name).FirstOrCreate(&role).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM user_roles WHERE user_id = ?", user.ID).Error; err != nil {
			return err
		}
		if err := tx.Model(&user).Association("Roles").Append(&role); err != nil {
			return err
		}
		permissionRecords := make([]models.Permission, 0, len(permissions))
		for _, name := range permissions {
			permission := models.Permission{}
			if err := tx.Where("name = ?", name).First(&permission).Error; errors.Is(err, gorm.ErrRecordNotFound) {
				permission = models.Permission{Name: name, GuardName: "web"}
				if err := tx.Create(&permission).Error; err != nil {
					return err
				}
			} else if err != nil {
				return err
			}
			permissionRecords = append(permissionRecords, permission)
		}
		if err := tx.Exec("DELETE FROM role_permissions WHERE role_id = ?", role.ID).Error; err != nil {
			return err
		}
		if err := tx.Model(&role).Association("Permissions").Append(&permissionRecords); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return fmt.Errorf("create superuser and grants: %w", err)
	}

	login, err := services.NewAuthService().Login(dtos.LoginRequest{Email: email, Password: password})
	if err != nil {
		return fmt.Errorf("verify seeded login: %w", err)
	}
	if len(login.Roles) == 0 || len(login.Permissions) != len(permissions) {
		return fmt.Errorf("verify seeded grants: expected %d permissions, received %d", len(permissions), len(login.Permissions))
	}

	fmt.Printf("Email: %s\nPassword: %s\nRole: Super Admin\nPermissions: %d protected API route permissions\n", email, password, len(login.Permissions))
	return nil
}

func generatePassword() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generate password: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func protectedRoutePermissions() ([]string, error) {
	methodActions := map[string]string{
		http.MethodGet:    "view",
		http.MethodPost:   "create",
		http.MethodPut:    "update",
		http.MethodPatch:  "update",
		http.MethodDelete: "delete",
	}
	publicPaths := map[string]bool{
		"/api/signup": true, "/api/verify": true, "/api/login": true,
		"/api/forgot-password": true, "/api/reset-password": true,
		"/api/withdrawal-callback": true, "/api/mpesa-b2c-result": true,
		"/api/mpesa-b2c-timeout": true, "/api/airtel-disbursement-callback": true,
		"/api/jenga-mobile-callback": true,
	}

	paths := routes.SetupRouter().Routes()
	permissionSet := make(map[string]struct{})
	for _, route := range paths {
		if !strings.HasPrefix(route.Path, "/api/") || publicPaths[route.Path] {
			continue
		}
		action, ok := methodActions[route.Method]
		if !ok {
			continue
		}
		resourceParts := make([]string, 0)
		for _, part := range strings.Split(strings.Trim(route.Path, "/"), "/") {
			if part == "api" || part == "" || strings.HasPrefix(part, ":") || isNumeric(part) {
				continue
			}
			resourceParts = append(resourceParts, part)
		}
		resource := strings.Join(resourceParts, ".")
		if strings.HasSuffix(route.Path, "/import") || strings.HasSuffix(route.Path, "/export") {
			permissionSet[resource] = struct{}{}
		} else {
			permissionSet[resource+"."+action] = struct{}{}
		}
	}

	permissions := make([]string, 0, len(permissionSet))
	for permission := range permissionSet {
		permissions = append(permissions, permission)
	}
	sort.Strings(permissions)
	return permissions, nil
}

func isNumeric(value string) bool {
	_, err := strconv.Atoi(value)
	return err == nil
}
