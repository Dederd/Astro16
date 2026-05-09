package services

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"

	"bouquet-app/internal/database"
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
// Design Cache helpers
// ────────────────────────────────────────────────────────────

// buildCacheKey membuat hash SHA256 dari kombinasi bouquet_type + sorted flower IDs & qty
func buildCacheKey(bouquetTypeID string, flowers []models.SelectedFlower) string {
	type flowerEntry struct {
		ID  string `json:"id"`
		Qty int    `json:"qty"`
	}
	entries := make([]flowerEntry, 0, len(flowers))
	for _, f := range flowers {
		entries = append(entries, flowerEntry{ID: f.FlowerID, Qty: f.Quantity})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].ID < entries[j].ID })

	raw, _ := json.Marshal(map[string]interface{}{
		"type":    bouquetTypeID,
		"flowers": entries,
	})

	h := sha256.Sum256(raw)
	return fmt.Sprintf("%x", h)[:32]
}

// lookupDesignCache cek apakah kombinasi bunga sudah pernah di-generate
func lookupDesignCache(cacheKey string) (*models.GenerateBouquetResponse, bool) {
	var cache models.DesignCacheDB
	if err := database.DB.First(&cache, "cache_key = ?", cacheKey).Error; err != nil {
		return nil, false
	}

	var result models.GenerateBouquetResponse
	if err := json.Unmarshal([]byte(cache.DesignsJSON), &result.Designs); err != nil {
		log.Printf("[DesignCache] gagal parse designs: %v", err)
		return nil, false
	}
	result.Message = cache.Message

	// Increment hit count
	database.DB.Model(&models.DesignCacheDB{}).
		Where("cache_key = ?", cacheKey).
		UpdateColumn("hit_count", database.DB.Raw("hit_count + 1"))

	log.Printf("[DesignCache] HIT untuk key=%s (hit #%d+)", cacheKey, cache.HitCount+1)
	return &result, true
}

// saveDesignCache simpan hasil generate ke cache
func saveDesignCache(cacheKey, bouquetTypeID string, flowers []models.SelectedFlower, result *models.GenerateBouquetResponse) {
	designsJSON, err := json.Marshal(result.Designs)
	if err != nil {
		log.Printf("[DesignCache] gagal marshal designs: %v", err)
		return
	}

	flowerComboJSON, _ := json.Marshal(flowers)

	cache := models.DesignCacheDB{
		CacheKey:      cacheKey,
		BouquetTypeID: bouquetTypeID,
		FlowerCombo:   string(flowerComboJSON),
		DesignsJSON:   string(designsJSON),
		Message:       result.Message,
		HitCount:      0,
	}

	if err := database.DB.Create(&cache).Error; err != nil {
		log.Printf("[DesignCache] gagal simpan cache: %v", err)
	} else {
		log.Printf("[DesignCache] SAVED key=%s", cacheKey)
	}
}

// ────────────────────────────────────────────────────────────
// Agent 2 — Generate desain bouquet
// SINKRON: stem_count di output AI = total tangkai yg dipilih user
// ────────────────────────────────────────────────────────────

func Agent2GenerateBouquet(req models.GenerateBouquetRequest) (*models.GenerateBouquetResponse, error) {
	// ── Cache lookup: jika kombinasi bunga sama dan tidak ada hint, pakai hasil lama ──
	cacheKey := buildCacheKey(req.BouquetTypeID, req.SelectedFlowers)
	// Skip cache jika user memberikan style/description hint (agar AI customisasi)
	hasHints := req.StyleHint != "" || req.DescriptionHint != ""
	if !hasHints {
		if cached, ok := lookupDesignCache(cacheKey); ok {
			return cached, nil
		}
	}

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

	// Harga paket FLAT:
	// Mini  = Rp35.000
	// Premium = Rp75.000
	const priceSmall int64 = 35000
	const priceLarge int64 = 75000

	// Stem count:
	// Mini  = total bunga user
	// Premium = total bunga user × 2 (lebih lebat)
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

	system := `Kamu adalah desainer florist profesional berpengalaman dari Indonesia dengan keahlian menciptakan bouquet yang elegan dan balance.
PRINSIP DESAIN BOUQUETMU:
1. KOMPOSISI: Buat focal point dengan bunga utama (rose, peony) di tengah, dikelilingi bunga sekunder, dan baby's breath/eucalyptus sebagai filler
2. WARNA: Pilih palette yang harmonis - kombinasi warna warm (pink, peach, red) atau cool (white, purple) yang cocok dengan bunga pilihan
3. VOLUME: Mini = bunga padat namun elegant, Premium = lebih voluminous dan mewah dengan layer yang dalam
4. WRAPPING: Selalu gunakan elegant wrapping - kraft paper dengan ribbon atau silk paper yang sophisticated
5. DETAIL: Tambahkan greenery (eucalyptus, ruscus, salal) untuk depth dan texture yang indah

ATURAN WAJIB - SANGAT PENTING:
- SEMUA bunga yang dipilih customer HARUS ada dalam bouquet. Tidak boleh ada yang dihilangkan.
- Bunga pilihan customer adalah KOMPONEN UTAMA bouquet - ornamen/filler hanya pelengkap tambahan
- Sebutkan setiap bunga pilihan customer secara eksplisit dalam description dan image_prompt
- Jawab HANYA dengan JSON murni tanpa markdown
- Stem count harus SAMA dengan nilai di prompt - jangan ubah
- Image prompt harus DETAILED dan SPESIFIK untuk generate image berkualitas tinggi`

	// Build optional hints section
	var hintsSection string
	if req.StyleHint != "" || req.DescriptionHint != "" {
		hintsSection = "\n\nPreferensi Customer:"
		if req.StyleHint != "" {
			hintsSection += fmt.Sprintf("\n- Gaya: %s", req.StyleHint)
		}
		if req.DescriptionHint != "" {
			hintsSection += fmt.Sprintf("\n- Konteks: %s", req.DescriptionHint)
		}
	}

	userPrompt := fmt.Sprintf(`Desain bouquet untuk acara: %s
Bunga pilihan customer (SEMUA HARUS ADA dalam bouquet):
%s%s

TECHNICAL REQUIREMENTS:
- Total bunga: %d jenis
- Mini Bouquet: Rp%d untuk %d stem (balance antara beauty dan size)
- Premium Bouquet: Rp%d untuk %d stem (lebih voluminous)

WAJIB: Buat 1 desain bouquet yang menggunakan SEMUA bunga di atas sebagai komponen utama.
Kamu boleh menambahkan baby's breath, eucalyptus, atau greenery lain sebagai filler pelengkap.
Tapi SEMUA bunga pilihan customer harus disebutkan dan ada di bouquet.
Wrapping elegant dengan color yang cocok.

Jawab ONLY dengan JSON (no markdown):
{
  "message": "Pesan inspirasi 2-3 kalimat yang menyebut semua bunga pilihan, momen, dan design bouquet",
  "designs": [
    {
      "id": "design_1",
      "name": "Nama design (3-4 kata bahasa Indonesia yang elegan)",
      "description": "Deskripsi design 1-2 kalimat - jelaskan semua bunga pilihan, focal point, color palette, dan overall feel",
      "style": "Style description (contoh: Romantic Modern, Minimalist Elegant, Luxe Glamour)",
      "image_prompt": "professional florist bouquet with [LIST ALL SELECTED FLOWERS BY NAME] arranged in [arrangement style], elegant [color] wrapping, premium quality, professional photography, studio lighting, white background, flowers are fresh and vibrant, composition is balanced, perfect for [occasion]",
      "small": {
        "label": "Mini Bouquet",
        "price": %d,
        "description": "Ukuran elegan untuk personal atau hadiah intimate",
        "stem_count": %d
      },
      "large": {
        "label": "Premium Bouquet",
        "price": %d,
        "description": "Ukuran mewah dengan volume lebih besar dan impact visual yang kuat",
        "stem_count": %d
      }
    }
  ]
}`,
		bouquetTypeName, flowerDetails.String(), hintsSection, totalStemCount,
		priceSmall, stemSmall,
		priceLarge, stemLarge,
		// design_1
		priceSmall, stemSmall, priceLarge, stemLarge,
	)

	response, err := callGroq(system, userPrompt)
	if err != nil {
		log.Printf("[Agent2] callGroq error: %v — menggunakan fallback", err)
		return agent2Fallback(stemSmall, stemLarge), nil
	}

	cleaned := extractJSON(response)
	log.Printf("[Agent2] Cleaned JSON: %.500s", cleaned)

	var result models.GenerateBouquetResponse
	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		log.Printf("[Agent2] JSON parse error: %v — menggunakan fallback", err)
		return agent2Fallback(stemSmall, stemLarge), nil
	}

	// Pastikan stem_count & price konsisten dengan pilihan user (override jika AI mengubah)
	for i := range result.Designs {
		result.Designs[i].SmallSize.Price = priceSmall
		result.Designs[i].SmallSize.StemCount = stemSmall
		result.Designs[i].LargeSize.Price = priceLarge
		result.Designs[i].LargeSize.StemCount = stemLarge
	}

	// ── Simpan ke cache untuk penggunaan berikutnya ──
	go saveDesignCache(cacheKey, req.BouquetTypeID, req.SelectedFlowers, &result)

	return &result, nil
}

func agent2Fallback(stemSmall, stemLarge int) *models.GenerateBouquetResponse {
	const small int64 = 35000
	const large int64 = 75000
	return &models.GenerateBouquetResponse{
		Message: "Kombinasi bunga yang indah! Berikut desain bouquet elegan untuk momen spesialmu.",
		Designs: []models.BouquetDesign{
			{
				ID: "design_1", Name: "Elegance Bloom", Style: "Romantic Elegance",
				Description: "Desain balanced dengan focal point bunga utama, dikelilingi filler yang memberi depth dan texture.",
				ImagePrompt: "professional florist bouquet arrangement with elegant wrapping, balanced composition, premium quality, fresh flowers, studio lighting",
				SmallSize:   models.SizeVariant{Label: "Mini Bouquet", Price: small, Description: "Ukuran elegan untuk personal atau hadiah intimate", StemCount: stemSmall},
				LargeSize:   models.SizeVariant{Label: "Premium Bouquet", Price: large, Description: "Ukuran mewah dengan volume lebih besar", StemCount: stemLarge},
			},
		},
	}
}
