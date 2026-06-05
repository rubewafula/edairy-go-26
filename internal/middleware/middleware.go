package middleware

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

		// permissions
		permissions := []string{}
		if p, ok := claims["permissions"]; ok {
			if arr, ok := p.([]interface{}); ok {
				for _, v := range arr {
					if s, ok := v.(string); ok {
						permissions = append(permissions, s)
					}
				}
			}
		}
		c.Set("permissions", permissions)

		// roles (optional)
		roles := []string{}
		if r, ok := claims["roles"]; ok {
			if arr, ok := r.([]interface{}); ok {
				for _, v := range arr {
					if s, ok := v.(string); ok {
						roles = append(roles, s)
					}
				}
			}
		}
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
