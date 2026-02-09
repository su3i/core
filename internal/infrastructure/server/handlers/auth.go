package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/darksuei/suei-intelligence/internal/application/account"
	accountService "github.com/darksuei/suei-intelligence/internal/application/account"
	"github.com/darksuei/suei-intelligence/internal/application/authentication"
	"github.com/darksuei/suei-intelligence/internal/application/mfa"
	"github.com/darksuei/suei-intelligence/internal/config"
	authenticationDomain "github.com/darksuei/suei-intelligence/internal/domain/authentication"
	"github.com/darksuei/suei-intelligence/internal/infrastructure/cache"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
)

func Login(c *gin.Context) {
	// Parse the request body
	var req struct {
		Email string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Invalid request: Missing required fields.",
		})
		return
	}

	var databaseCfg config.DatabaseConfig
	if err := envconfig.Process("", &databaseCfg); err != nil {
		log.Fatalf("Failed to load database config: %v", err)
	}

	_account, err := account.RetrieveAccountWithPassword(req.Email, req.Password, &databaseCfg)

	if err != nil || _account == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	challengeID := uuid.New().String()

	err = cache.GetCache().Set(fmt.Sprintf("challenge-id-%s", challengeID), req.Email, time.Hour)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"mfa_required": true,
		"challenge_id": challengeID,
	})
	return
}

func MFA(c *gin.Context) {
	// Parse the request body
	var req struct {
		ChallengeID string `json:"challenge_id" binding:"required"`
		Code string `json:"code" binding:"required"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Invalid request: Missing required fields.",
		})
		return
	}

	challengeKey := fmt.Sprintf("challenge-id-%s", req.ChallengeID)

	email, err := cache.GetCache().Get(challengeKey)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Please restart login flow.",
		})
		return
	}

	var commonCfg config.CommonConfig
	if err := envconfig.Process("", &commonCfg); err != nil {
		log.Fatalf("Failed to load common config: %v", err)
	}

	var databaseCfg config.DatabaseConfig
	if err := envconfig.Process("", &databaseCfg); err != nil {
		log.Fatalf("Failed to load database config: %v", err)
	}

	// Retrieve account
	_account, err := accountService.RetrieveAccount(email, &databaseCfg)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if _account == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid account.",
		})
		return
	}

	codeUint64, err := strconv.ParseUint(req.Code, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid mfa code format"})
		return
	}

	code := uint32(codeUint64)

	// Confirm and enable MFA
	isCodeValid := mfa.VerifyTOTP(_account.MFASecret, code, time.Now())

	if !isCodeValid {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid TOTP code.",
		})
		return
	}

	defer func() {
		_ = cache.GetCache().Delete(challengeKey)
	}()

	auth, err := authentication.LoginWithoutPassword(email, &commonCfg, &databaseCfg)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"access_token": auth.AccessToken,
		"refresh_token": auth.RefreshToken,
	})
	return
}

func RevokeToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "refresh_token is required",
		})
		return
	}

	err := cache.GetCache().Delete(fmt.Sprintf("refresh-token-%s", authenticationDomain.HashRefreshToken(req.RefreshToken)))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
	return
}

func RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "refresh_token is required",
		})
		return
	}

	var commonCfg config.CommonConfig
	if err := envconfig.Process("", &commonCfg); err != nil {
		log.Fatalf("Failed to load common config: %v", err)
	}

	var databaseCfg config.DatabaseConfig
	if err := envconfig.Process("", &databaseCfg); err != nil {
		log.Fatalf("Failed to load database config: %v", err)
	}

	authTokens, err := authentication.Refresh(req.RefreshToken, &commonCfg, &databaseCfg)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid or expired refresh token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"access_token":  authTokens.AccessToken,
		"refresh_token": authTokens.RefreshToken,
	})
}
