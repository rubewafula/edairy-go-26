package middleware

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rubewafula/edairy-go-26/internal/db"
)

type contextKey string

const UserContextKey contextKey = "auth_user"

type AuthUser struct {
	UserID      uint64
	Permissions []string
	Roles       []string
}

func AuthMiddleware(jwtSecret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.AbortWithStatusJSON(401, gin.H{"Error": "missing token"})
			return
		}

		tokenStr = strings.Replace(tokenStr, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"Error": "unauthorized"})
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		// user_id
		userID, ok := claims["user_id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"Error": "invalid user id"})
			return
		}
		c.Set("user_id", uint64(userID))

		// Fetch roles for the user using raw SQL
		var roles []string
		roleQuery := `
			SELECT r.name 
			FROM roles r 
			INNER JOIN user_roles ur ON ur.role_id = r.id 
			WHERE ur.user_id = ? AND r.deleted_at IS NULL`
		if err := db.WithContext(c.Request.Context()).Raw(roleQuery, uint64(userID)).Scan(&roles).Error; err != nil {
			c.AbortWithStatusJSON(500, gin.H{"Error": "internal server error: session data fetch failed"})
			return
		}

		// Fetch unique permissions (combined from roles and direct assignments) using raw SQL
		var permissions []string
		permQuery := `
			SELECT DISTINCT p.name
			FROM permissions p

			LEFT JOIN role_permissions rp ON rp.permission_id = p.id
			LEFT JOIN user_roles ur ON ur.role_id = rp.role_id AND ur.user_id = ?

			LEFT JOIN user_permissions up ON up.permission_id = p.id AND up.user_id = ?

			WHERE p.deleted_at IS NULL
			AND (ur.user_id IS NOT NULL OR up.user_id IS NOT NULL)
			`
		if err := db.WithContext(c.Request.Context()).Raw(permQuery, uint64(userID), uint64(userID)).Scan(&permissions).Error; err != nil {
			c.AbortWithStatusJSON(500, gin.H{"Error": "internal server error: session data fetch failed"})
			return
		}

		c.Set("permissions", permissions)
		c.Set("roles", roles)

		authUser := AuthUser{
			UserID:      uint64(userID),
			Permissions: permissions,
			Roles:       roles,
		}

		ctx := context.WithValue(c.Request.Context(), UserContextKey, authUser)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func extractResource(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")

	resources := make([]string, 0)

	for _, p := range parts {

		if p == "api" || p == "" {
			continue
		}
		if _, err := strconv.Atoi(p); err == nil {
			continue
		}

		if strings.HasPrefix(p, ":") {
			continue
		}

		resources = append(resources, p)
	}

	return strings.Join(resources, ".")
}

func RequireAutoPermission() gin.HandlerFunc {
	return func(c *gin.Context) {

		method := c.Request.Method
		path := c.FullPath()
		log.Printf("[RBAC] Incoming request => Method: %s Path: %s", method, path)

		resource := extractResource(path)

		var required string

		switch {
		case strings.HasSuffix(path, "/import"):
			required = resource

		case strings.HasSuffix(path, "/export"):
			required = resource

		default:
			var action string

			switch method {
			case "GET":
				action = "view"
			case "POST":
				action = "create"
			case "PUT", "PATCH":
				action = "update"
			case "DELETE":
				action = "delete"
			default:
				c.AbortWithStatusJSON(405, gin.H{"error": "method not allowed"})
				return
			}

			required = resource + "." + action
		}

		val, exists := c.Get("permissions")

		if !exists {
			log.Printf("[RBAC] No permissions found in context")

			c.AbortWithStatusJSON(403, gin.H{"Error": "no permissions"})
			return
		}

		perms := val.([]string)
		//log.Printf("[RBAC] User permissions: %+v", perms)

		for _, p := range perms {
			if p == required {
				c.Next()
				return
			}
		}
		log.Printf("[RBAC] Forbidden: %s", required)

		c.AbortWithStatusJSON(403, gin.H{
			"Error": "forbidden: " + required,
		})
	}
}
