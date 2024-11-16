package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/repository/redis"
	"github.com/ShavelSoSmetanoi/messenger-backend/pkg"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// EmailValidate structure for registration request
type EmailValidate struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// EmailValidator sends a code to the email and sets a timeout
func EmailValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req EmailValidate

		// Extract JSON data from the request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request format.",
			})
			c.Abort()
			return
		}

		email := req.Email

		ctx := context.Background()
		if _, err := redis.Rdb.Get(ctx, email).Result(); err == nil {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Please wait before requesting another code.",
			})
			c.Abort()
			return
		}

		code := "12345"

		err := sendCodeToEmail(email, code)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to send email.",
			})
			c.Abort()
			return
		}

		uuid := pkg.GenerateUniqueID()

		userData := map[string]string{
			"username": req.Username,
			"email":    req.Email,
			"password": req.Password,
			"code":     code,
		}

		userDataJSON, err := json.Marshal(userData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to process registration data.",
			})
			c.Abort()
			return
		}

		// Save the user data in Redis with the generated UUID as the key
		err = redis.Rdb.Set(ctx, uuid, userDataJSON, 5*time.Minute).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to save verification code.",
			})
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"UUID": uuid,
		})
		c.Abort()
	}
}

type VerifyCodeRequest struct {
	Code string `json:"code" binding:"required"` // Verification code entered by the user
	UUID string `json:"uuid" binding:"required"` // User's UUID
}

// VerifyCode checks the verification code sent by the user
func VerifyCode() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req VerifyCodeRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request format.",
			})
			c.Abort()
			return
		}

		ctx := context.Background()

		userDataJSON, err := redis.Rdb.Get(ctx, req.UUID).Result()
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid UUID or verification code.",
			})
			c.Abort()
			return
		}

		var userData map[string]string
		if err := json.Unmarshal([]byte(userDataJSON), &userData); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to parse user data.",
			})
			c.Abort()
			return
		}

		if userData["code"] != req.Code {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid verification code.",
			})
			c.Abort()
			return
		}

		c.Set("userData", userData)

		c.Next()
	}
}

// sendCodeToEmail function for sending verification code to email
func sendCodeToEmail(email, code string) error {
	fmt.Printf("Sent verification code %s to %s\n", code, email)
	return nil
}
