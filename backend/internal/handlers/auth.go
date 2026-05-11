package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"

	"bouquet-app/internal/database"
	"bouquet-app/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(getEnvAuth("JWT_SECRET", "bouquet-secret-key-2025"))

func getEnvAuth(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// Register godoc
// @Summary Register user baru
// @Tags auth
// @Accept json
// @Produce json
// @Param body body models.RegisterRequest true "Register request"
// @Success 201 {object} map[string]interface{}
// @Router /auth/register [post]
func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid: " + err.Error()})
		return
	}

	// Check email sudah ada
	var existing models.UserDB
	if result := database.DB.Where("email = ?", req.Email).First(&existing); result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email sudah terdaftar"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal enkripsi password"})
		return
	}

	user := models.UserDB{
		Name:         req.Name,
		Email:        req.Email,
		Phone:        req.Phone,
		PasswordHash: string(hash),
	}
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat akun"})
		return
	}

	token, err := generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": buildAuthResponse(user, token)})
}

// Login godoc
// @Summary Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param body body models.LoginRequest true "Login request"
// @Success 200 {object} map[string]interface{}
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	var user models.UserDB
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
		return
	}

	token, err := generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": buildAuthResponse(user, token)})
}

// GetMe godoc
// @Summary Ambil data user yang sedang login
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /auth/me [get]
func GetMe(c *gin.Context) {
	userID := c.GetUint("user_id")
	var user models.UserDB
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{
		"id": user.ID, "name": user.Name, "email": user.Email, "phone": user.Phone,
	}})
}

// GetUserOrders godoc
// @Summary Ambil semua order milik user yang login
// @Tags user
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /user/orders [get]
func GetUserOrders(c *gin.Context) {
	userID := c.GetUint("user_id")
	var orders []models.OrderDB
	if err := database.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pesanan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": orders})
}

// ForgotPassword godoc
// @Summary Kirim link reset password ke email
// @Tags auth
// @Accept json
// @Produce json
// @Param body body models.ForgotPasswordRequest true "Forgot password request"
// @Success 200 {object} map[string]interface{}
// @Router /auth/forgot-password [post]
func ForgotPassword(c *gin.Context) {
	var req models.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email tidak valid"})
		return
	}

	var user models.UserDB
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		// Selalu return 200 agar tidak mengekspos apakah email terdaftar
		c.JSON(http.StatusOK, gin.H{"message": "Jika email terdaftar, link reset password sudah dikirim"})
		return
	}

	// Hapus token lama yang belum dipakai
	database.DB.Where("user_id = ? AND used = false", user.ID).Delete(&models.PasswordResetTokenDB{})

	// Buat token acak
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
		return
	}
	tokenStr := hex.EncodeToString(tokenBytes)

	resetToken := models.PasswordResetTokenDB{
		UserID:    user.ID,
		Token:     tokenStr,
		ExpiresAt: time.Now().Add(1 * time.Hour),
		Used:      false,
	}
	if err := database.DB.Create(&resetToken).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan token"})
		return
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", frontendURL, tokenStr)

	go sendResetEmail(user.Email, user.Name, resetLink)

	c.JSON(http.StatusOK, gin.H{"message": "Jika email terdaftar, link reset password sudah dikirim"})
}

// ResetPassword godoc
// @Summary Reset password dengan token dari email
// @Tags auth
// @Accept json
// @Produce json
// @Param body body models.ResetPasswordRequest true "Reset password request"
// @Success 200 {object} map[string]interface{}
// @Router /auth/reset-password [post]
func ResetPassword(c *gin.Context) {
	var req models.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	var resetToken models.PasswordResetTokenDB
	if err := database.DB.Where("token = ? AND used = false", req.Token).First(&resetToken).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Link reset tidak valid atau sudah digunakan"})
		return
	}

	if time.Now().After(resetToken.ExpiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Link reset sudah expired. Silakan minta link baru"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal enkripsi password"})
		return
	}

	if err := database.DB.Model(&models.UserDB{}).Where("id = ?", resetToken.UserID).Update("password_hash", string(hash)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update password"})
		return
	}

	database.DB.Model(&resetToken).Update("used", true)

	c.JSON(http.StatusOK, gin.H{"message": "Password berhasil diubah. Silakan login dengan password baru"})
}

func sendResetEmail(toEmail, toName, resetLink string) {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	fromEmail := os.Getenv("SMTP_FROM")

	if smtpHost == "" || smtpUser == "" || smtpPass == "" {
		log.Printf("[sendResetEmail] SMTP tidak dikonfigurasi. Reset link: %s", resetLink)
		return
	}
	if smtpPort == "" {
		smtpPort = "587"
	}
	if fromEmail == "" {
		fromEmail = smtpUser
	}

	subject := "Reset Password Bloome"
	body := fmt.Sprintf(`Halo %s,

Kami menerima permintaan reset password untuk akun Bloome kamu.

Klik link berikut untuk membuat password baru:
%s

Link ini berlaku selama 1 jam. Jika kamu tidak meminta reset password, abaikan email ini.

Salam,
Tim Bloome 🌸`, toName, resetLink)

	msg := fmt.Sprintf("From: Bloome <%s>\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		fromEmail, toEmail, subject, body)

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	addr := smtpHost + ":" + smtpPort
	if err := smtp.SendMail(addr, auth, fromEmail, []string{toEmail}, []byte(msg)); err != nil {
		log.Printf("[sendResetEmail] Gagal kirim email ke %s: %v", toEmail, err)
	} else {
		log.Printf("[sendResetEmail] Email berhasil dikirim ke %s", toEmail)
	}
}

// ─── Middleware ──────────────────────────────────────────────

// AuthMiddleware validates JWT token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token diperlukan"})
			c.Abort()
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
			c.Abort()
			return
		}
		if idStr, ok := claims["user_id"].(string); ok {
			if id, err := strconv.ParseUint(idStr, 10, 64); err == nil {
				c.Set("user_id", uint(id))
			}
		} else if idFloat, ok := claims["user_id"].(float64); ok {
			c.Set("user_id", uint(idFloat))
		}
		c.Next()
	}
}

// OptionalAuthMiddleware — set user_id if token present, don't block if absent
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header != "" && strings.HasPrefix(header, "Bearer ") {
			tokenStr := strings.TrimPrefix(header, "Bearer ")
			claims := jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
				return jwtSecret, nil
			})
			if err == nil && token.Valid {
				if idFloat, ok := claims["user_id"].(float64); ok {
					c.Set("user_id", uint(idFloat))
				}
			}
		}
		c.Next()
	}
}

// ─── Helpers ────────────────────────────────────────────────

func generateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(30 * 24 * time.Hour).Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret)
}

func buildAuthResponse(user models.UserDB, token string) models.AuthResponse {
	var resp models.AuthResponse
	resp.Token = token
	resp.User.ID = user.ID
	resp.User.Name = user.Name
	resp.User.Email = user.Email
	resp.User.Phone = user.Phone
	return resp
}
