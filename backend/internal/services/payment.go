package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"bouquet-app/internal/database"
	"bouquet-app/internal/models"
)

// GenerateOrderID membuat ID order unik
func GenerateOrderID() string {
	return fmt.Sprintf("BQ-%d", time.Now().UnixMilli())
}

// SaveOrder menyimpan order baru ke PostgreSQL
func SaveOrder(order *models.Order) error {
	record := &models.OrderDB{
		ID:               order.ID,
		CustomerName:     order.CustomerName,
		CustomerEmail:    order.CustomerEmail,
		CustomerPhone:    order.CustomerPhone,
		BouquetTypeID:    order.BouquetTypeID,
		SelectedFlowers:  models.SelectedFlowersJSON(order.SelectedFlowers),
		DesignID:         order.DesignID,
		DesignName:       order.DesignName,
		Size:             order.Size,
		TotalAmount:      order.TotalAmount,
		Status:           order.Status,
		Notes:            order.Notes,
		ShippingAddress:  order.ShippingAddress,
		ShippingCity:     order.ShippingCity,
		ShippingPostcode: order.ShippingPostcode,
		ShippingPhone:    order.ShippingPhone,
		CourierService:   order.CourierService,
		OrderSource:      order.OrderSource,
		CatalogItemID:    order.CatalogItemID,
		UserID:           order.UserID,
		FlowerCost:       order.FlowerCost,
		MakingFee:        order.MakingFee,
		AIFee:            order.AIFee,
		ShippingCost:     order.ShippingCost,
		CreatedAt:        order.CreatedAt,
	}
	result := database.DB.Create(record)
	if result.Error != nil {
		log.Printf("[SaveOrder] DB error: %v", result.Error)
	}
	return result.Error
}

// GetOrderByID mengambil order berdasarkan ID dari PostgreSQL
func GetOrderByID(id string) (*models.Order, bool) {
	var record models.OrderDB
	if err := database.DB.First(&record, "id = ?", id).Error; err != nil {
		log.Printf("[GetOrderByID] tidak ditemukan: %s, error: %v", id, err)
		return nil, false
	}
	return record.ToOrder(), true
}

// GetAllOrders untuk admin
func GetAllOrders() []models.OrderDB {
	var records []models.OrderDB
	database.DB.Order("created_at DESC").Find(&records)
	return records
}

// UpdateOrderStatus mengupdate status dan payment_id di PostgreSQL
func UpdateOrderStatus(id, status, paymentID string) error {
	updates := map[string]interface{}{"status": status}
	if paymentID != "" {
		updates["payment_id"] = paymentID
	}
	result := database.DB.Model(&models.OrderDB{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		log.Printf("[UpdateOrderStatus] error: %v", result.Error)
	}
	return result.Error
}

// UpdateOrderShipping update info pengiriman dari admin
func UpdateOrderShipping(id string, req models.UpdateOrderStatusRequest) error {
	updates := map[string]interface{}{
		"status": req.Status,
	}
	if req.TrackingNumber != "" {
		updates["tracking_number"] = req.TrackingNumber
	}
	if req.ShippingStatus != "" {
		updates["shipping_status"] = req.ShippingStatus
	}
	if req.CourierService != "" {
		updates["courier_service"] = req.CourierService
	}
	result := database.DB.Model(&models.OrderDB{}).Where("id = ?", id).Updates(updates)
	return result.Error
}

// ────────────────────────────────────────────────────────────
// Generate session (rate limiting)
// ────────────────────────────────────────────────────────────

func GetOrCreateSession(sessionID string) (*models.GenerateSessionDB, error) {
	var session models.GenerateSessionDB
	result := database.DB.First(&session, "session_id = ?", sessionID)
	if result.Error != nil {
		// Buat baru
		session = models.GenerateSessionDB{SessionID: sessionID, GenerateCount: 0, IsPaid: false}
		if err := database.DB.Create(&session).Error; err != nil {
			return nil, err
		}
	}
	return &session, nil
}

func IncrementGenerateCount(sessionID string) error {
	return database.DB.Model(&models.GenerateSessionDB{}).
		Where("session_id = ?", sessionID).
		UpdateColumn("generate_count", database.DB.Raw("generate_count + 1")).Error
}

// ────────────────────────────────────────────────────────────
// Midtrans Snap
// ────────────────────────────────────────────────────────────

func MidtransCreateToken(order *models.Order) (*models.PaymentTokenResponse, error) {
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	if serverKey == "" {
		return nil, fmt.Errorf("MIDTRANS_SERVER_KEY tidak diisi di .env")
	}

	isProduction := os.Getenv("MIDTRANS_IS_PRODUCTION") == "true"
	baseURL := "https://app.sandbox.midtrans.com/snap/v1/transactions"
	if isProduction {
		baseURL = "https://app.midtrans.com/snap/v1/transactions"
	}

	designLabel := order.DesignName
	if designLabel == "" {
		designLabel = order.CatalogItemID
	}

	payload := map[string]interface{}{
		"transaction_details": map[string]interface{}{
			"order_id":     order.ID,
			"gross_amount": order.TotalAmount,
		},
		"customer_details": map[string]interface{}{
			"first_name": order.CustomerName,
			"email":      order.CustomerEmail,
			"phone":      order.CustomerPhone,
			"shipping_address": map[string]interface{}{
				"first_name": order.CustomerName,
				"phone":      order.ShippingPhone,
				"address":    order.ShippingAddress,
				"city":       order.ShippingCity,
				"postal_code": order.ShippingPostcode,
			},
		},
		"item_details": []map[string]interface{}{
			{
				"id":       order.DesignID,
				"price":    order.TotalAmount,
				"quantity": 1,
				"name":     fmt.Sprintf("Bouquet - %s (%s)", designLabel, order.Size),
			},
		},
		"callbacks": map[string]interface{}{
			"finish": os.Getenv("FRONTEND_URL") + "/payment/finish",
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("gagal marshal payload Midtrans: %w", err)
	}

	log.Printf("[Midtrans] Payload: %s", string(body))

	req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("gagal buat HTTP request Midtrans: %w", err)
	}

	auth := base64.StdEncoding.EncodeToString([]byte(serverKey + ":"))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gagal panggil Midtrans API: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("gagal baca response Midtrans: %w", err)
	}

	log.Printf("[Midtrans] HTTP status: %d, Response: %s", resp.StatusCode, string(respBody))

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return nil, fmt.Errorf("Midtrans API error (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	var midtransResp struct {
		Token       string   `json:"token"`
		RedirectURL string   `json:"redirect_url"`
		ErrorMessages []string `json:"error_messages"`
	}
	if err := json.Unmarshal(respBody, &midtransResp); err != nil {
		return nil, fmt.Errorf("gagal parse response Midtrans: %w — body: %s", err, string(respBody))
	}

	if len(midtransResp.ErrorMessages) > 0 {
		return nil, fmt.Errorf("Midtrans validation error: %v", midtransResp.ErrorMessages)
	}

	if midtransResp.Token == "" {
		return nil, fmt.Errorf("token Midtrans kosong — response: %s", string(respBody))
	}

	return &models.PaymentTokenResponse{
		Token:       midtransResp.Token,
		RedirectURL: midtransResp.RedirectURL,
	}, nil
}

// ────────────────────────────────────────────────────────────
// Tracking JNE/JNT (BinderByte API — gratis)
// ────────────────────────────────────────────────────────────

func GetTrackingInfo(courierService, trackingNumber string) (interface{}, error) {
	apiKey := os.Getenv("BINDERBYTE_API_KEY")
	if apiKey == "" {
		// Return mock jika API key belum diset
		return map[string]interface{}{
			"summary": map[string]interface{}{
				"courier":     courierService,
				"awb":         trackingNumber,
				"status":      "IN_TRANSIT",
				"description": "Paket sedang dalam perjalanan",
			},
			"history": []map[string]interface{}{
				{"date": time.Now().Format("2006-01-02 15:04:05"), "desc": "Paket diterima oleh kurir", "location": "Jakarta"},
			},
		}, nil
	}

	url := fmt.Sprintf("https://api.binderbyte.com/v1/track?api_key=%s&courier=%s&awb=%s",
		apiKey, courierService, trackingNumber)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("gagal panggil tracking API: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result interface{}
	json.Unmarshal(body, &result)
	return result, nil
}
