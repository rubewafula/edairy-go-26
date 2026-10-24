package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// queryIntParam returns the first positive integer among keys. List endpoints
// historically mixed "page"/"Page" and "limit"/"Limit"; the mobile app sends
// lowercase, so accept either.
func queryIntParam(ctx *gin.Context, def int, keys ...string) int {
	for _, key := range keys {
		raw := ctx.Query(key)
		if raw == "" {
			continue
		}
		if value, err := strconv.Atoi(raw); err == nil && value > 0 {
			return value
		}
	}
	return def
}
