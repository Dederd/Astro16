package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"bouquet-app/internal/database"
	"bouquet-app/internal/models"
	"bouquet-app/internal/services"

	"github.com/gin-gonic/gin"
)

// ────────────────────────────────────────────────────────────
// Bouquet Types
// ────────────────────────────────────────────────────────────

// GetBouquetTypes godoc
// @Summary      Ambil semua tipe bouquet
// @Tags         bouquet
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /bouquet-types [get]
func GetBouquetTypes(c *gin.Context) {
	types := services.GetAllBouquetTypes()
	c.JSON(http.StatusOK, gin.H{"data": types})
}

// ────────────────────────────────────────────────────────────
// Flowers
// ────────────────────────────────────────────────────────────

// GetFlowers godoc
// @Summary      Ambil semua bunga dari database
// @Tags         flowers
// @Produce      json
// @Param        occasion  query  string  false  "Filter by occasion ID"
// @Success      200  {object}  map[string]interface{}
// @Router       /flowers [get]
func GetFlowers(c *gin.Context) {
	occasionID := c.Query("occasion")
	flowers := services.GetAllFlowers()

	if occasionID != "" {
		filtered := make([]models.Flower, 0)
		for _, f := range flowers {
			for _, occ := range f.Occasions {
				if occ == occasionID {
					filtered = append(filtered, f)
					break
				}
			}
		}
		c.JSON(http.StatusOK, gin.H{"data": filtered})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": flowers})
}

// ────────────────────────────────────────────────────────────
// Catalog Bouquet
// ────────────────────────────────────────────────────────────

// GetCatalog godoc
// @Summary      Ambil katalog bouquet pre-made
// @Tags         catalog
// @Produce      json
// @Param        occasion  query  string  false  "Filter by occasion"
// @Success      200  {object}  map[string]interface{}
// @Router       /catalog [get]
func GetCatalog(c *gin.Context) {
	catalogs := services.GetAllCatalog()
	occasionFilter := c.Query("occasion")

	if occasionFilter != "" {
		filtered := make([]models.CatalogBouquetDB, 0)
		for _, cat := range catalogs {
			if cat.Occasion == occasionFilter {
				filtered = append(filtered, cat)
			}
		}
		c.JSON(http.StatusOK, gin.H{"data": filtered})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": catalogs})
}

// ────────────────────────────────────────────────────────────
// AI Agents
// ────────────────────────────────────────────────────────────

// AgentVerifySelection godoc
// @Summary      Agent 1: Verifikasi pilihan momen & rekomendasikan bunga
// @Tags         agent
// @Accept       json
// @Produce      json
// @Param        body  body  models.AgentVerifyRequest  true  "Request body"
// @Success      200  {object}  map[string]interface{}
// @Router       /agent/verify-selection [post]
func AgentVerifySelection(c *gin.Context) {
	var req models.AgentVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	allFlowers := services.GetAllFlowers()
	result, err := services.Agent1VerifySelection(req, allFlowers)
	if err != nil {
		log.Printf("[AgentVerifySelection] error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

// AgentGenerateBouquet godoc
// @Summary      Agent 2: Generate desain bouquet berdasarkan bunga pilihan
// @Tags         agent
// @Accept       json
// @Produce      json
// @Param        X-Session-ID  header  string  true  "Session ID untuk rate limiting"
// @Param        body  body  models.GenerateBouquetRequest  true  "Request body"
// @Success      200  {object}  map[string]interface{}
// @Failure      429  {object}  map[string]interface{}  "Batas generate tercapai"
// @Router       /agent/generate-bouquet [post]
func AgentGenerateBouquet(c *gin.Context) {
	var req models.GenerateBouquetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(req.SelectedFlowers) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pilih minimal 1 jenis bunga"})
		return
	}

	// Rate limiting: gunakan user_id jika login, fallback ke session ID
	sessionID := c.GetHeader("X-Session-ID")
	if userID, exists := c.Get("user_id"); exists {
		if uid, ok := userID.(uint); ok && uid > 0 {
			sessionID = fmt.Sprintf("user_%d", uid)
		}
	}
	if sessionID == "" {
		sessionID = c.ClientIP()
	}

	const maxFree = 2
	const aiFeePerGenerate = int64(5000)

	session, err := services.GetOrCreateSession(sessionID)
	if err != nil {
		log.Printf("[AgentGenerateBouquet] session error: %v", err)
	} else {
		effectiveLimit := maxFree + session.ExtraQuota
		if session.GenerateCount >= effectiveLimit {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":          "Kuota generate habis. Beli 3 kuota tambahan seharga Rp5.000.",
				"generate_count": session.GenerateCount,
				"limit":          effectiveLimit,
				"is_limited":     true,
				"ai_fee":         aiFeePerGenerate,
			})
			return
		}
	}

	// Hitung total stem count
	totalStemCount := 0
	for _, sf := range req.SelectedFlowers {
		totalStemCount += sf.Quantity
	}
	req.TotalStemCount = totalStemCount

	result, err := services.Agent2GenerateBouquet(req)
	if err != nil {
		log.Printf("[AgentGenerateBouquet] error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal generate desain: " + err.Error()})
		return
	}

	if session != nil {
		services.IncrementGenerateCount(sessionID)
	}

	effectiveLimit := maxFree
	if session != nil {
		effectiveLimit = maxFree + session.ExtraQuota
	}

	// Apakah generate ini berbayar (di atas free tier)?
	isThisPaid := session != nil && session.GenerateCount >= maxFree
	aiFee := int64(0)
	if isThisPaid {
		aiFee = aiFeePerGenerate
	}

	c.JSON(http.StatusOK, gin.H{
		"data":             result,
		"generate_count":   session.GenerateCount + 1,
		"limit":            effectiveLimit,
		"ai_fee":           aiFee,
		"is_paid_generate": isThisPaid,
	})
}

// BuyGenerateQuota godoc
// @Summary      Buat token pembayaran Midtrans untuk beli 3 kuota generate (Rp5.000)
// @Tags         agent
// @Produce      json
// @Router       /agent/buy-quota [post]
func BuyGenerateQuota(c *gin.Context) {
	sessionID := c.GetHeader("X-Session-ID")
	var userIDVal *uint
	if userID, exists := c.Get("user_id"); exists {
		if uid, ok := userID.(uint); ok && uid > 0 {
			userIDVal = &uid
			sessionID = fmt.Sprintf("user_%d", uid)
		}
	}
	if sessionID == "" {
		sessionID = c.ClientIP()
	}

	// Buat order sementara untuk kuota generate (Midtrans membutuhkan order object)
	quotaOrderID := fmt.Sprintf("QUOTA-%s-%d", sessionID[:min(len(sessionID), 12)], time.Now().UnixMilli())
	quotaOrder := &models.Order{
		ID:            quotaOrderID,
		CustomerName:  "Quota Purchase",
		CustomerEmail: "quota@bloome.id",
		CustomerPhone: "08000000000",
		TotalAmount:   5000,
		Status:        "pending_quota",
		OrderSource:   "quota",
		DesignName:    "Paket 3 Generate",
		Size:          "paket",
		UserID:        userIDVal,
		CreatedAt:     time.Now(),
	}

	tokenResp, err := services.MidtransCreateQuotaToken(quotaOrder, sessionID)
	if err != nil {
		log.Printf("[BuyGenerateQuota] Midtrans error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Gagal membuat token pembayaran kuota",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"snap_token":  tokenResp.Token,
		"redirect_url": tokenResp.RedirectURL,
		"order_id":    quotaOrderID,
		"session_id":  sessionID,
		"amount":      5000,
	})
}

// ConfirmQuotaPayment godoc
// @Summary      Konfirmasi pembayaran kuota setelah Midtrans callback
// @Tags         agent
// @Accept       json
// @Produce      json
// @Router       /agent/confirm-quota [post]
func ConfirmQuotaPayment(c *gin.Context) {
	var req struct {
		OrderID   string `json:"order_id"`
		SessionID string `json:"session_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sessionID := req.SessionID
	if sessionID == "" {
		sessionID = c.GetHeader("X-Session-ID")
		if userID, exists := c.Get("user_id"); exists {
			if uid, ok := userID.(uint); ok && uid > 0 {
				sessionID = fmt.Sprintf("user_%d", uid)
			}
		}
	}

	session, err := services.BuyExtraQuota(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambah kuota: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "Berhasil membeli 3 kuota generate",
		"generate_count":  session.GenerateCount,
		"limit":           2 + session.ExtraQuota,
		"extra_quota":     session.ExtraQuota,
		"extra_quota_fee": session.ExtraQuotaFee,
		"is_limited":      session.GenerateCount >= (2 + session.ExtraQuota),
	})
}

// GetGenerateStatus godoc
// @Summary      Cek status generate session (sisa kuota)
// @Tags         agent
// @Produce      json
// @Param        X-Session-ID  header  string  true  "Session ID"
// @Success      200  {object}  map[string]interface{}
// @Router       /agent/generate-status [get]
func GetGenerateStatus(c *gin.Context) {
	sessionID := c.GetHeader("X-Session-ID")
	if userID, exists := c.Get("user_id"); exists {
		if uid, ok := userID.(uint); ok && uid > 0 {
			sessionID = fmt.Sprintf("user_%d", uid)
		}
	}
	if sessionID == "" {
		sessionID = c.ClientIP()
	}
	session, err := services.GetOrCreateSession(sessionID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"generate_count": 0, "limit": 2, "is_limited": false, "ai_fee": 5000})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"generate_count": session.GenerateCount,
		"limit":          2 + session.ExtraQuota,
		"extra_quota":    session.ExtraQuota,
		"is_limited":     session.GenerateCount >= (2 + session.ExtraQuota),
		"is_paid":        session.IsPaid,
		"ai_fee":         5000,
	})
}

// ────────────────────────────────────────────────────────────
// Orders
// ────────────────────────────────────────────────────────────

// CreateOrder godoc
// @Summary      Buat order baru
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        body  body  models.CreateOrderRequest  true  "Request body"
// @Success      201  {object}  map[string]interface{}
// @Router       /orders [post]
func CreateOrder(c *gin.Context) {
	var req models.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[CreateOrder] bind error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.OrderSource == "" {
		req.OrderSource = "ai_generated"
	}

	order := &models.Order{
		ID:               services.GenerateOrderID(),
		CustomerName:     req.CustomerName,
		CustomerEmail:    req.CustomerEmail,
		CustomerPhone:    req.CustomerPhone,
		BouquetTypeID:    req.BouquetTypeID,
		SelectedFlowers:  req.SelectedFlowers,
		DesignID:         req.DesignID,
		DesignName:       req.DesignName,
		Size:             req.Size,
		TotalAmount:      req.TotalAmount,
		Notes:            req.Notes,
		ShippingAddress:  req.ShippingAddress,
		ShippingCity:     req.ShippingCity,
		ShippingPostcode: req.ShippingPostcode,
		ShippingPhone:    req.ShippingPhone,
		CourierService:   req.CourierService,
		OrderSource:      req.OrderSource,
		CatalogItemID:    req.CatalogItemID,
		FlowerCost:       req.FlowerCost,
		MakingFee:        req.MakingFee,
		AIFee:            req.AIFee,
		ShippingCost:     req.ShippingCost,
		Status:           "pending",
		CreatedAt:        time.Now(),
	}

	// Link order ke akun user jika sudah login (set oleh OptionalAuthMiddleware)
	if userID, exists := c.Get("user_id"); exists {
		if uid, ok := userID.(uint); ok && uid > 0 {
			order.UserID = &uid
		}
	}

	if err := services.SaveOrder(order); err != nil {
		log.Printf("[CreateOrder] save error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan order: " + err.Error()})
		return
	}

	log.Printf("[CreateOrder] berhasil: %s, amount: %d", order.ID, order.TotalAmount)
	c.JSON(http.StatusCreated, gin.H{"data": order})
}

// CreatePaymentToken godoc
// @Summary      Buat token pembayaran Midtrans Snap
// @Tags         payment
// @Accept       json
// @Produce      json
// @Param        body  body  models.PaymentTokenRequest  true  "Request body"
// @Success      200  {object}  map[string]interface{}
// @Router       /payment/token [post]
func CreatePaymentToken(c *gin.Context) {
	var req models.PaymentTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, ok := services.GetOrderByID(req.OrderID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order tidak ditemukan: " + req.OrderID})
		return
	}

	log.Printf("[CreatePaymentToken] order: %s, amount: %d", order.ID, order.TotalAmount)

	tokenResp, err := services.MidtransCreateToken(order)
	if err != nil {
		log.Printf("[CreatePaymentToken] Midtrans error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Gagal membuat token pembayaran",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tokenResp})
}

// PaymentNotification godoc
// @Summary      Webhook notifikasi pembayaran dari Midtrans
// @Tags         payment
// @Accept       json
// @Produce      json
// @Router       /payment/notification [post]
func PaymentNotification(c *gin.Context) {
	var notification map[string]interface{}
	if err := c.ShouldBindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orderID, _ := notification["order_id"].(string)
	transactionStatus, _ := notification["transaction_status"].(string)
	transactionID, _ := notification["transaction_id"].(string)
	customField1, _ := notification["custom_field1"].(string) // session ID untuk quota orders

	log.Printf("[PaymentNotification] order: %s, status: %s, custom_field1: %s", orderID, transactionStatus, customField1)

	var status string
	switch transactionStatus {
	case "capture", "settlement":
		status = "paid"
	case "pending":
		status = "pending"
	case "deny", "expire", "cancel":
		status = "failed"
	default:
		status = "unknown"
	}

	// Check if this is a quota payment (order_id starts with "QUOTA-")
	if status == "paid" && customField1 != "" && len(orderID) > 5 && orderID[:6] == "QUOTA-" {
		log.Printf("[PaymentNotification] Processing quota payment for session: %s", customField1)
		// Add quota to session
		if session, err := services.BuyExtraQuota(customField1); err != nil {
			log.Printf("[PaymentNotification] Error adding quota: %v", err)
		} else {
			log.Printf("[PaymentNotification] Quota added successfully. Extra quota: %d", session.ExtraQuota)
		}
	} else {
		// Regular bouquet order
		if err := services.UpdateOrderStatus(orderID, status, transactionID); err != nil {
			log.Printf("[PaymentNotification] update error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update status order"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// GetOrder godoc
// @Summary      Ambil detail order berdasarkan ID
// @Tags         orders
// @Produce      json
// @Param        id  path  string  true  "Order ID"
// @Success      200  {object}  map[string]interface{}
// @Router       /orders/{id} [get]
func GetOrder(c *gin.Context) {
	id := c.Param("id")
	order, ok := services.GetOrderByID(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": order})
}

// GetTracking godoc
// @Summary      Ambil info tracking pengiriman
// @Tags         orders
// @Produce      json
// @Param        id  path  string  true  "Order ID"
// @Success      200  {object}  map[string]interface{}
// @Router       /orders/{id}/tracking [get]
func GetTracking(c *gin.Context) {
	id := c.Param("id")
	order, ok := services.GetOrderByID(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order tidak ditemukan"})
		return
	}

	trackingInfo := models.TrackingInfo{
		OrderID:        order.ID,
		TrackingNumber: order.TrackingNumber,
		CourierService: order.CourierService,
		ShippingStatus: order.ShippingStatus,
	}

	// Ambil data real dari kurir jika ada resi
	if order.TrackingNumber != "" {
		courierData, err := services.GetTrackingInfo(order.CourierService, order.TrackingNumber)
		if err != nil {
			log.Printf("[GetTracking] courier API error: %v", err)
		} else {
			trackingInfo.CourierData = courierData
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": trackingInfo})
}

// ────────────────────────────────────────────────────────────
// Admin Handlers
// ────────────────────────────────────────────────────────────

// AdminGetOrders godoc
// @Summary      [Admin] Ambil semua order
// @Tags         admin
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /admin/orders [get]
func AdminGetOrders(c *gin.Context) {
	orders := services.GetAllOrders()
	c.JSON(http.StatusOK, gin.H{"data": orders, "total": len(orders)})
}

// AdminUpdateOrder godoc
// @Summary      [Admin] Update status order & info pengiriman
// @Tags         admin
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        id    path  string                          true  "Order ID"
// @Param        body  body  models.UpdateOrderStatusRequest true  "Request body"
// @Success      200   {object}  map[string]interface{}
// @Router       /admin/orders/{id} [put]
func AdminUpdateOrder(c *gin.Context) {
	id := c.Param("id")
	var req models.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.UpdateOrderShipping(id, req); err != nil {
		log.Printf("[AdminUpdateOrder] error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order berhasil diupdate"})
}

// AdminGetFlowers godoc
// @Summary      [Admin] Ambil semua bunga
// @Tags         admin
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /admin/flowers [get]
func AdminGetFlowers(c *gin.Context) {
	flowers := services.GetAllFlowers()
	c.JSON(http.StatusOK, gin.H{"data": flowers})
}

// AdminUpdateFlower godoc
// @Summary      [Admin] Update data bunga (harga, stok, gambar, ketersediaan)
// @Tags         admin
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        id    path  string                     true  "Flower ID"
// @Param        body  body  models.UpdateFlowerRequest true  "Request body"
// @Success      200   {object}  map[string]interface{}
// @Router       /admin/flowers/{id} [put]
func AdminUpdateFlower(c *gin.Context) {
	id := c.Param("id")
	var req models.UpdateFlowerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if req.Price != nil {
		updates["price"] = *req.Price
	}
	if req.IsAvailable != nil {
		updates["is_available"] = *req.IsAvailable
	}
	if req.Stock != nil {
		updates["stock"] = *req.Stock
	}
	if req.ImageURL != nil {
		updates["image_url"] = *req.ImageURL
	}
	if req.Emoji != nil {
		updates["emoji"] = *req.Emoji
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tidak ada field yang diupdate"})
		return
	}

	result := database.DB.Model(&models.FlowerDB{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bunga tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bunga berhasil diupdate"})
}

// AdminCreateFlower godoc
// @Summary      [Admin] Tambah bunga baru
// @Tags         admin
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Success      201  {object}  map[string]interface{}
// @Router       /admin/flowers [post]
func AdminCreateFlower(c *gin.Context) {
	var flower models.FlowerDB
	if err := c.ShouldBindJSON(&flower); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if flower.ID == "" || flower.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID dan Nama wajib diisi"})
		return
	}
	if err := database.DB.Create(&flower).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Bunga berhasil ditambahkan", "data": flower})
}

// AdminDeleteFlower godoc
// @Summary      [Admin] Hapus bunga
// @Tags         admin
// @Security     ApiKeyAuth
// @Param        id  path  string  true  "Flower ID"
// @Success      200  {object}  map[string]interface{}
// @Router       /admin/flowers/{id} [delete]
func AdminDeleteFlower(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.FlowerDB{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Bunga berhasil dihapus"})
}

// AdminGetCatalog godoc
// @Summary      [Admin] Ambil semua katalog (termasuk yg tidak aktif)
// @Tags         admin
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /admin/catalog [get]
func AdminGetCatalog(c *gin.Context) {
	catalogs := services.GetAllCatalogAdmin()
	c.JSON(http.StatusOK, gin.H{"data": catalogs})
}

// AdminCreateCatalog godoc
// @Summary      [Admin] Tambah item katalog baru
// @Tags         admin
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        body  body  models.CreateCatalogRequest  true  "Request body"
// @Success      201   {object}  map[string]interface{}
// @Router       /admin/catalog [post]
func AdminCreateCatalog(c *gin.Context) {
	var req models.CreateCatalogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	record := models.CatalogBouquetDB{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		Style:       req.Style,
		Occasion:    req.Occasion,
		Price:       req.Price,
		StemCount:   req.StemCount,
		IsAvailable: req.IsAvailable,
		Stock:       req.Stock,
		SortOrder:   req.SortOrder,
	}

	if err := database.DB.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": record})
}

// AdminUpdateCatalog godoc
// @Summary      [Admin] Update item katalog
// @Tags         admin
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        id    path  string  true  "Catalog ID"
// @Param        body  body  models.CreateCatalogRequest  true  "Request body"
// @Success      200   {object}  map[string]interface{}
// @Router       /admin/catalog/{id} [put]
func AdminUpdateCatalog(c *gin.Context) {
	id := c.Param("id")
	var req models.CreateCatalogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{
		"name":         req.Name,
		"description":  req.Description,
		"image_url":    req.ImageURL,
		"style":        req.Style,
		"occasion":     req.Occasion,
		"price":        req.Price,
		"stem_count":   req.StemCount,
		"is_available": req.IsAvailable,
		"stock":        req.Stock,
		"sort_order":   req.SortOrder,
	}

	result := database.DB.Model(&models.CatalogBouquetDB{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Katalog berhasil diupdate"})
}

// AdminDeleteCatalog godoc
// @Summary      [Admin] Hapus item katalog
// @Tags         admin
// @Security     ApiKeyAuth
// @Produce      json
// @Param        id  path  string  true  "Catalog ID"
// @Success      200  {object}  map[string]interface{}
// @Router       /admin/catalog/{id} [delete]
func AdminDeleteCatalog(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.CatalogBouquetDB{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Katalog berhasil dihapus"})
}

// AdminGetStats godoc
// @Summary      [Admin] Statistik dashboard
// @Tags         admin
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /admin/stats [get]
func AdminGetStats(c *gin.Context) {
	var totalOrders, paidOrders, pendingOrders int64
	var totalRevenue struct{ Sum int64 }

	database.DB.Model(&models.OrderDB{}).Count(&totalOrders)
	database.DB.Model(&models.OrderDB{}).Where("status = 'paid'").Count(&paidOrders)
	database.DB.Model(&models.OrderDB{}).Where("status = 'pending'").Count(&pendingOrders)
	database.DB.Model(&models.OrderDB{}).Where("status = 'paid'").
		Select("COALESCE(SUM(total_amount), 0) as sum").Scan(&totalRevenue)

	c.JSON(http.StatusOK, gin.H{
		"total_orders":   totalOrders,
		"paid_orders":    paidOrders,
		"pending_orders": pendingOrders,
		"total_revenue":  totalRevenue.Sum,
	})
}

// AdminNotifyNewOrder godoc
// @Summary      [Admin] Terima notifikasi pesanan baru yang sudah dibayar
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        body  body  map[string]interface{}  true  "Order notification data"
// @Success      200  {object}  map[string]interface{}
// @Router       /admin/notify-new-order [post]
func AdminNotifyNewOrder(c *gin.Context) {
	var req struct {
		OrderID      string `json:"order_id"`
		CustomerName string `json:"customer_name"`
		CustomerPhone string `json:"customer_phone"`
		TotalAmount  int64  `json:"total_amount"`
		DesignName   string `json:"design_name"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Log notifikasi ke console/file
	log.Printf("[ADMIN NOTIF] New paid order: %s | Customer: %s | Phone: %s | Amount: %d | Design: %s\n",
		req.OrderID, req.CustomerName, req.CustomerPhone, req.TotalAmount, req.DesignName)

	// Di masa depan, bisa simpan ke database atau kirim ke notification system
	// Untuk sekarang, frontend akan menampilkan notifikasi berdasarkan status order

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Notifikasi diterima",
	})
}

// AdminMiddleware — simple API key check untuk endpoint admin
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-Admin-Key")
		adminKey := "admin-bouquet-2024" // ganti dengan env variable di production
		if key != adminKey {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}
