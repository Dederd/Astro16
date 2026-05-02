package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"bouquet-app/internal/models"
)

// ────────────────────────────────────────────────────────────
// Groq API structs
// ────────────────────────────────────────────────────────────

type groqRequest struct {
	Model       string        `json:"model"`
	MaxTokens   int           `json:"max_tokens"`
	Temperature float64       `json:"temperature"`
	Messages    []groqMessage `json:"messages"`
}

type groqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type groqResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

// ────────────────────────────────────────────────────────────
// callGroq — single unified Groq caller
// ────────────────────────────────────────────────────────────

func callGroq(systemPrompt, userPrompt string) (string, error) {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("GROQ_API_KEY tidak diisi di .env")
	}

	payload := groqRequest{
		Model:       "llama-3.1-8b-instant",
		MaxTokens:   2000,
		Temperature: 0.7,
		Messages: []groqMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("gagal marshal request Groq: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("gagal buat HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("gagal panggil Groq API: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("gagal baca response body: %w", err)
	}

	log.Printf("[Groq] HTTP status: %d", resp.StatusCode)
	if resp.StatusCode != 200 {
		log.Printf("[Groq] Error body: %s", string(respBody))
		return "", fmt.Errorf("Groq API error (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	var groqResp groqResponse
	if err := json.Unmarshal(respBody, &groqResp); err != nil {
		return "", fmt.Errorf("gagal parse response Groq: %w", err)
	}
	if groqResp.Error != nil {
		return "", fmt.Errorf("Groq error: %s", groqResp.Error.Message)
	}
	if len(groqResp.Choices) == 0 {
		return "", fmt.Errorf("response Groq kosong (tidak ada choices)")
	}

	text := groqResp.Choices[0].Message.Content
	log.Printf("[Groq] Raw response (200 chars): %.200s", text)
	return text, nil
}

// ────────────────────────────────────────────────────────────
// extractJSON — bersihkan markdown fence dari respons
// ────────────────────────────────────────────────────────────

func extractJSON(raw string) string {
	raw = strings.TrimSpace(raw)

	re := regexp.MustCompile(`(?s)` + "```" + `(?:json)?\s*([\s\S]+?)\s*` + "```")
	if matches := re.FindStringSubmatch(raw); len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	start := strings.Index(raw, "{")
	end := strings.LastIndex(raw, "}")
	if start != -1 && end != -1 && end > start {
		return strings.TrimSpace(raw[start : end+1])
	}
	return raw
}

// ────────────────────────────────────────────────────────────
// Agent 1 — Verifikasi pilihan momen & rekomendasikan bunga
// ────────────────────────────────────────────────────────────

func Agent1VerifySelection(req models.AgentVerifyRequest, allFlowers []models.Flower) (*models.AgentVerifyResponse, error) {
	var flowerList strings.Builder
	for _, f := range allFlowers {
		status := "READY"
		if !f.IsAvailable {
			status = "HABIS"
		}
		flowerList.WriteString(fmt.Sprintf(
			"- ID: %s | %s (%s) | Rp%d/tangkai | Status: %s | Cocok untuk: %s | Makna: %s\n",
			f.ID, f.NameID, f.Name, f.Price, status, strings.Join(f.Occasions, ","), f.Meaning,
		))
	}

	system := `Kamu adalah asisten ahli florist profesional dari toko bunga premium Indonesia.
Tugasmu adalah membantu customer memilih bunga yang tepat untuk bouquet mereka.
Selalu jawab dalam Bahasa Indonesia yang ramah dan profesional.
PENTING: Jawab HANYA dengan JSON murni tanpa markdown, tanpa penjelasan tambahan.`

	userPrompt := fmt.Sprintf(`Customer ingin membuat bouquet untuk acara: %s

Daftar bunga yang tersedia:
%s

Jawab HANYA dengan JSON ini (tanpa markdown, tanpa teks lain):
{"message":"pesan sambutan 2-3 kalimat","recommendations":["flower_id_1","flower_id_2"],"tips":"tips memilih bunga 2-3 kalimat"}`,
		req.BouquetType, flowerList.String())

	response, err := callGroq(system, userPrompt)
	if err != nil {
		log.Printf("[Agent1] callGroq error: %v — menggunakan fallback", err)
		return agent1Fallback(req, allFlowers), nil
	}

	cleaned := extractJSON(response)
	log.Printf("[Agent1] Cleaned JSON: %.300s", cleaned)

	var result models.AgentVerifyResponse
	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		log.Printf("[Agent1] JSON parse error: %v — menggunakan fallback", err)
		return agent1Fallback(req, allFlowers), nil
	}
	return &result, nil
}

func agent1Fallback(req models.AgentVerifyRequest, allFlowers []models.Flower) *models.AgentVerifyResponse {
	var recommended []string
	for _, f := range allFlowers {
		if f.IsAvailable {
			for _, occ := range f.Occasions {
				if occ == req.BouquetTypeID {
					recommended = append(recommended, f.ID)
					break
				}
			}
		}
	}
	return &models.AgentVerifyResponse{
		Message:         fmt.Sprintf("Pilihan bouquet %s Anda sudah tepat! Kami punya berbagai bunga cantik untuk acara ini.", req.BouquetType),
		Recommendations: recommended,
		Tips:            "Pilih bunga sesuai warna dan makna yang paling mewakili momen spesialmu.",
	}
}

// ────────────────────────────────────────────────────────────
// Agent 2 — Generate desain bouquet
// SINKRON: stem_count di output AI = total tangkai yg dipilih user
// ────────────────────────────────────────────────────────────

func Agent2GenerateBouquet(req models.GenerateBouquetRequest) (*models.GenerateBouquetResponse, error) {
	allFlowers := GetAllFlowers()
	flowerMap := make(map[string]models.Flower)
	for _, f := range allFlowers {
		flowerMap[f.ID] = f
	}

	// Hitung total tangkai & harga dari pilihan user
	var flowerDetails strings.Builder
	var totalFlowerPrice int64
	totalStemCount := 0

	for _, sf := range req.SelectedFlowers {
		if f, ok := flowerMap[sf.FlowerID]; ok {
			flowerDetails.WriteString(fmt.Sprintf("- %s: %d tangkai @ Rp%d\n", f.NameID, sf.Quantity, f.Price))
			totalFlowerPrice += f.Price * int64(sf.Quantity)
			totalStemCount += sf.Quantity
		}
	}

	// Gunakan total dari request jika dikirim (lebih akurat)
	if req.TotalStemCount > 0 {
		totalStemCount = req.TotalStemCount
	}

	// Harga paket:
	// Mini  = harga bunga yg dipilih + packaging Rp30.000
	// Premium = (harga bunga × 2) + packaging Rp75.000
	priceSmall := totalFlowerPrice + 30000
	priceLarge := totalFlowerPrice*2 + 75000

	// Stem count:
	// Mini  = total tangkai user
	// Premium = total tangkai user × 2 (lebih lebat)
	stemSmall := totalStemCount
	stemLarge := totalStemCount * 2

	bouquetTypes := GetAllBouquetTypes()
	var bouquetTypeName string
	for _, bt := range bouquetTypes {
		if bt.ID == req.BouquetTypeID {
			bouquetTypeName = bt.Name
			break
		}
	}
	if bouquetTypeName == "" {
		bouquetTypeName = req.BouquetTypeID
	}

	system := `Kamu adalah desainer bouquet bunga profesional dari Indonesia.
PENTING: Jawab HANYA dengan JSON murni. Jangan gunakan markdown, jangan ada teks sebelum atau sesudah JSON.
Stem count harus SAMA PERSIS dengan nilai yang diberikan di prompt — jangan mengubah angka tersebut.`

	userPrompt := fmt.Sprintf(`Buat 3 desain bouquet berbeda untuk acara %s dengan bunga pilihan customer:
%s

Total tangkai yang dipilih customer: %d tangkai
Harga paket Mini: Rp%d (untuk %d tangkai — JANGAN UBAH angka ini)
Harga paket Premium: Rp%d (untuk %d tangkai — JANGAN UBAH angka ini)

Jawab HANYA dengan JSON ini (tanpa markdown):
{
  "message": "pesan inspiratif 2-3 kalimat dalam Bahasa Indonesia",
  "designs": [
    {
      "id": "design_1",
      "name": "nama desain Bahasa Indonesia",
      "description": "deskripsi singkat 1-2 kalimat Bahasa Indonesia",
      "style": "Romantic",
      "image_prompt": "detailed english prompt for bouquet image generation",
      "small": {"label": "Mini Bouquet", "price": %d, "description": "deskripsi mini Bahasa Indonesia", "stem_count": %d},
      "large": {"label": "Premium Bouquet", "price": %d, "description": "deskripsi premium Bahasa Indonesia", "stem_count": %d}
    },
    {
      "id": "design_2",
      "name": "nama desain Bahasa Indonesia",
      "description": "deskripsi singkat",
      "style": "Elegant",
      "image_prompt": "detailed english prompt",
      "small": {"label": "Mini Bouquet", "price": %d, "description": "deskripsi mini", "stem_count": %d},
      "large": {"label": "Premium Bouquet", "price": %d, "description": "deskripsi premium", "stem_count": %d}
    },
    {
      "id": "design_3",
      "name": "nama desain Bahasa Indonesia",
      "description": "deskripsi singkat",
      "style": "Modern",
      "image_prompt": "detailed english prompt",
      "small": {"label": "Mini Bouquet", "price": %d, "description": "deskripsi mini", "stem_count": %d},
      "large": {"label": "Premium Bouquet", "price": %d, "description": "deskripsi premium", "stem_count": %d}
    }
  ]
}`,
		bouquetTypeName, flowerDetails.String(), totalStemCount,
		priceSmall, stemSmall,
		priceLarge, stemLarge,
		// design_1
		priceSmall, stemSmall, priceLarge, stemLarge,
		// design_2
		priceSmall, stemSmall, priceLarge, stemLarge,
		// design_3
		priceSmall, stemSmall, priceLarge, stemLarge,
	)

	response, err := callGroq(system, userPrompt)
	if err != nil {
		log.Printf("[Agent2] callGroq error: %v — menggunakan fallback", err)
		return agent2Fallback(totalFlowerPrice, stemSmall, stemLarge), nil
	}

	cleaned := extractJSON(response)
	log.Printf("[Agent2] Cleaned JSON: %.500s", cleaned)

	var result models.GenerateBouquetResponse
	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		log.Printf("[Agent2] JSON parse error: %v — menggunakan fallback", err)
		return agent2Fallback(totalFlowerPrice, stemSmall, stemLarge), nil
	}

	// Pastikan stem_count & price konsisten dengan pilihan user (override jika AI mengubah)
	for i := range result.Designs {
		result.Designs[i].SmallSize.Price = priceSmall
		result.Designs[i].SmallSize.StemCount = stemSmall
		result.Designs[i].LargeSize.Price = priceLarge
		result.Designs[i].LargeSize.StemCount = stemLarge
	}

	return &result, nil
}

func agent2Fallback(totalFlowerPrice int64, stemSmall, stemLarge int) *models.GenerateBouquetResponse {
	small := totalFlowerPrice + 30000
	large := totalFlowerPrice*2 + 75000
	return &models.GenerateBouquetResponse{
		Message: "Kombinasi bunga yang indah! Berikut tiga konsep desain bouquet spesial untukmu.",
		Designs: []models.BouquetDesign{
			{
				ID: "design_1", Name: "Taman Romantis", Style: "Romantic",
				Description: "Desain klasik penuh kelembutan dengan dominasi warna hangat.",
				ImagePrompt: "romantic bouquet with mixed flowers, soft pink and red tones, elegant white wrap",
				SmallSize:   models.SizeVariant{Label: "Mini Bouquet", Price: small, Description: "Pas untuk hadiah personal", StemCount: stemSmall},
				LargeSize:   models.SizeVariant{Label: "Premium Bouquet", Price: large, Description: "Tampilan mewah dan mengesankan", StemCount: stemLarge},
			},
			{
				ID: "design_2", Name: "Pesona Elegan", Style: "Elegant",
				Description: "Desain modern dengan sentuhan mewah dan keanggunan tinggi.",
				ImagePrompt: "elegant bouquet with premium flowers, white and gold tones, luxury ribbon",
				SmallSize:   models.SizeVariant{Label: "Mini Bouquet", Price: small, Description: "Elegan dalam ukuran compact", StemCount: stemSmall},
				LargeSize:   models.SizeVariant{Label: "Premium Bouquet", Price: large, Description: "Kehadiran yang memukau", StemCount: stemLarge},
			},
			{
				ID: "design_3", Name: "Ceria Natural", Style: "Natural",
				Description: "Tampilan segar alami dengan warna-warna cerah yang menyenangkan.",
				ImagePrompt: "natural fresh bouquet with bright colorful flowers, green leaves, rustic twine",
				SmallSize:   models.SizeVariant{Label: "Mini Bouquet", Price: small, Description: "Cerah dan menyegarkan", StemCount: stemSmall},
				LargeSize:   models.SizeVariant{Label: "Premium Bouquet", Price: large, Description: "Kebun bunga dalam genggaman", StemCount: stemLarge},
			},
		},
	}
}
