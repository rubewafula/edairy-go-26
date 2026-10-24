package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/services"
)

// Endpoints a member (farmer) login may call. Everything else is refused for them.
// Key: METHOD + " " + gin FullPath.
var memberSelfServiceRoutes = map[string]bool{
	"GET /api/member-me":                         true,
	"GET /api/member-credit-limit":               true,
	"GET /api/milk-collection-overview":          true,
	"GET /api/loan-module/applications":          true,
	"GET /api/loan-module/mobile/quote":          true,
	"GET /api/loan-module/mobile/transactions":   true,
	"GET /api/loan-module/mobile/loan-requests":  true,
	"POST /api/loan-module/mobile/loan-requests": true,
	"GET /api/loan-module/mobile/loan-requests/:id": true,
	"POST /api/loan-module/mobile/loan-requests/:id/resend-otp": true,
	"POST /api/loan-module/mobile/loan-requests/:id/verify-otp": true,
}

// Query parameters that name a member; for a member login they are always forced to their own id.
var memberQueryParams = []string{"member_id", "memberId", "member", "customer_id"}

// MemberSelfServiceScope runs after AuthMiddleware. For a member login (users.member_id set) it:
//   - refuses any endpoint not in memberSelfServiceRoutes,
//   - forces every member id in the query string and JSON body to the member's own id,
//   - refuses loan requests (/:id) that belong to another member,
//   - marks the request so RequireAutoPermission lets it through (the Member role has no staff permissions).
//
// Staff logins are not affected.
func MemberSelfServiceScope() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, ok := c.Get("user_id")
		if !ok {
			c.Next()
			return
		}
		userID, _ := uid.(uint64)
		memberID := services.MemberIDForUser(userID)
		if memberID == 0 {
			c.Next()
			return
		}
		c.Set("member_id", memberID)

		key := c.Request.Method + " " + c.FullPath()
		if !memberSelfServiceRoutes[key] {
			log.Printf("[MemberScope] member %d (user %d) refused %s", memberID, userID, key)
			c.AbortWithStatusJSON(403, gin.H{"error": "This is not available on a member account."})
			return
		}

		own := strconv.FormatUint(memberID, 10)
		q := c.Request.URL.Query()
		for _, p := range memberQueryParams {
			q.Set(p, own)
		}
		c.Request.URL.RawQuery = q.Encode()

		if c.Param("id") != "" {
			appID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
			if appID == 0 || services.LoanApplicationMemberID(appID) != memberID {
				c.AbortWithStatusJSON(404, gin.H{"error": "loan request not found"})
				return
			}
		}

		if c.Request.Method == "POST" && c.Request.Body != nil {
			raw, _ := io.ReadAll(c.Request.Body)
			body := map[string]interface{}{}
			if len(bytes.TrimSpace(raw)) > 0 {
				if err := json.Unmarshal(raw, &body); err != nil {
					c.AbortWithStatusJSON(400, gin.H{"error": "invalid request body"})
					return
				}
			}
			if _, has := body["member_id"]; has || c.FullPath() == "/api/loan-module/mobile/loan-requests" {
				body["member_id"] = memberID
			}
			out, _ := json.Marshal(body)
			c.Request.Body = io.NopCloser(bytes.NewReader(out))
			c.Request.ContentLength = int64(len(out))
		}

		c.Set("member_self_service", true)
		c.Next()
	}
}
