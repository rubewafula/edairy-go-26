package controllers

import (
	"net/http"

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
