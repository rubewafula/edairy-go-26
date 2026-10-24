package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/services"
)

type AuthController struct {
	Service *services.AuthService
}

func NewAuthController() *AuthController {
	return &AuthController{
		Service: services.NewAuthService(),
	}
}

func (a *AuthController) Signup(c *gin.Context) {
	var req dtos.SignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := a.Service.Signup(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to signup"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "account created"})
}

func (a *AuthController) Verify(c *gin.Context) {
	token := c.Query("token")

	if err := a.Service.VerifyAccount(token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "verified"})
}

func (a *AuthController) Login(c *gin.Context) {
	var req dtos.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := a.Service.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid login"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (a *AuthController) ForgotPassword(c *gin.Context) {
	var req dtos.ForgotPasswordRequest

	c.ShouldBindJSON(&req)

	a.Service.ForgotPassword(req.Email)

	c.JSON(http.StatusOK, gin.H{"message": "if exists, reset sent"})
}
func (a *AuthController) ResetPassword(c *gin.Context) {
	var req dtos.ResetPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := a.Service.ResetPassword(req.Token, req.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or expired"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password reset"})
}

func (a *AuthController) ChangePassword(c *gin.Context) {
	var req dtos.ChangePasswordRequest

	userID := c.GetUint64("user_id")

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := a.Service.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

// MemberLogin POST /api/member-login {"phone": "07…|7…|254…", "id_number": "…"} — farmer sign-in.
func (a *AuthController) MemberLogin(c *gin.Context) {
	var req services.MemberLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Enter your phone number and ID number."})
		return
	}
	resp, err := services.MemberLogin(req, c.ClientIP())
	if err != nil {
		respondServiceError(c, err, nil)
		return
	}
	c.JSON(http.StatusOK, resp)
}

// MemberMe GET /api/member-me — the signed-in member's profile (member logins only).
func (a *AuthController) MemberMe(c *gin.Context) {
	memberID := c.GetUint64("member_id")
	if memberID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "not a member account"})
		return
	}
	profile, err := services.MemberSelfProfile(memberID)
	if err != nil {
		respondServiceError(c, err, nil)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": profile})
}

// MemberPinOTP POST /api/member-pin/otp {"phone","id_number"} — SMS code to set or reset the PIN.
func (a *AuthController) MemberPinOTP(c *gin.Context) {
	var req services.MemberPinOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Enter your phone number and ID number."})
		return
	}
	resp, err := services.RequestMemberPinOTP(req, c.ClientIP())
	if err != nil {
		respondServiceError(c, err, nil)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": resp, "message": "Code sent"})
}

// MemberPinReset POST /api/member-pin/reset {"phone","id_number","otp","pin"} — sets the PIN and signs in.
func (a *AuthController) MemberPinReset(c *gin.Context) {
	var req services.MemberPinResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Enter your phone number, ID number, SMS code and new PIN."})
		return
	}
	resp, err := services.ResetMemberPinWithOTP(req, c.ClientIP())
	if err != nil {
		respondServiceError(c, err, nil)
		return
	}
	c.JSON(http.StatusOK, resp)
}

// MemberPinChange POST /api/member-pin/change {"current_pin","new_pin"} — signed-in member.
func (a *AuthController) MemberPinChange(c *gin.Context) {
	memberID := c.GetUint64("member_id")
	if memberID == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "not a member account"})
		return
	}
	var req services.MemberPinChangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Enter your current PIN and the new PIN."})
		return
	}
	if err := services.ChangeMemberPin(memberID, req); err != nil {
		respondServiceError(c, err, nil)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "PIN changed"})
}

// ResetMemberPin POST /api/members/:id/reset-pin — office reset: clears the PIN; the member sets a
// new one with an SMS code.
func (a *AuthController) ResetMemberPin(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid member id"})
		return
	}
	if err := services.ResetMemberPinByOffice(id); err != nil {
		respondServiceError(c, err, nil)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "PIN cleared. The member sets a new PIN in the app with an SMS code."})
}

// MemberPinStatus GET /api/members/:id/pin-status — whether the member has set a login PIN.
func (a *AuthController) MemberPinStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"member_id": id, "has_pin": services.MemberHasPin(id)}})
}
