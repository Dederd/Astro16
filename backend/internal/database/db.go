package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"bouquet-app/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	dsn := buildDSN()

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
		PreferSimpleProtocol: true, // ← tambahkan ini
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	
	if err != nil {
		log.Fatalf("❌ Gagal konek ke database: %v", err)
	}


	// Auto-migrate semua tabel
	if err := db.AutoMigrate(
		&models.UserDB{},
		&models.OrderDB{},
		&models.FlowerDB{},
		&models.CatalogBouquetDB{},
		&models.GenerateSessionDB{},
		&models.DesignCacheDB{},
		&models.PasswordResetTokenDB{},
	); err != nil {
		log.Fatalf("❌ Gagal auto-migrate: %v", err)
	}

	DB = db
	log.Println("✅ Database PostgreSQL terhubung")

	// Seed data bunga jika tabel kosong
	seedFlowers()
	seedCatalog()
}

func buildDSN() string {
	if url := os.Getenv("DATABASE_URL"); url != "" {
		return url
	}
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "bouquet_db")
	sslmode := getEnv("DB_SSLMODE", "disable")
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Jakarta",
		host, port, user, password, dbname, sslmode,
	)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func toJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func seedFlowers() {
	var count int64
	DB.Model(&models.FlowerDB{}).Count(&count)
	if count > 0 {
		return
	}
	log.Println("🌸 Seeding bunga ke database...")

	flowers := []models.FlowerDB{
		{ID: "rose_red", Name: "Red Rose", NameID: "Mawar Merah", Description: "Simbol cinta sejati dan gairah.", Emoji: "🌹", Price: 15000, IsAvailable: true, Stock: 100, Colors: []string{"#C62828"}, Occasions: []string{"wedding", "valentines", "anniversary"}, Meaning: "Cinta dan gairah"},
		{ID: "rose_pink", Name: "Pink Rose", NameID: "Mawar Pink", Description: "Melambangkan kelembutan dan kasih sayang.", Emoji: "🌸", Price: 15000, IsAvailable: true, Stock: 100, Colors: []string{"#F48FB1"}, Occasions: []string{"graduation", "mothers_day", "birthday"}, Meaning: "Kelembutan dan apresiasi"},
		{ID: "rose_white", Name: "White Rose", NameID: "Mawar Putih", Description: "Melambangkan kesucian dan kemurnian.", Emoji: "🤍", Price: 15000, IsAvailable: true, Stock: 80, Colors: []string{"#FAFAFA"}, Occasions: []string{"wedding", "sympathy"}, Meaning: "Kesucian dan keanggunan"},
		{ID: "tulip_yellow", Name: "Yellow Tulip", NameID: "Tulip Kuning", Description: "Simbol kebahagiaan dan keceriaan.", Emoji: "🌷", Price: 18000, IsAvailable: true, Stock: 60, Colors: []string{"#FDD835"}, Occasions: []string{"birthday", "congratulations"}, Meaning: "Kebahagiaan yang cerah"},
		{ID: "tulip_purple", Name: "Purple Tulip", NameID: "Tulip Ungu", Description: "Melambangkan royalti dan keistimewaan.", Emoji: "💜", Price: 18000, IsAvailable: true, Stock: 60, Colors: []string{"#7B1FA2"}, Occasions: []string{"graduation", "congratulations"}, Meaning: "Royalti dan pencapaian"},
		{ID: "sunflower", Name: "Sunflower", NameID: "Bunga Matahari", Description: "Simbol kesetiaan dan keceriaan.", Emoji: "🌻", Price: 12000, IsAvailable: true, Stock: 120, Colors: []string{"#F9A825"}, Occasions: []string{"birthday", "congratulations", "graduation"}, Meaning: "Kesetiaan dan semangat"},
		{ID: "lily_white", Name: "White Lily", NameID: "Lili Putih", Description: "Elegan dan anggun.", Emoji: "🌼", Price: 20000, IsAvailable: true, Stock: 50, Colors: []string{"#FFFFFF"}, Occasions: []string{"wedding", "sympathy", "anniversary"}, Meaning: "Kesucian dan keanggunan"},
		{ID: "chrysanthemum", Name: "Chrysanthemum", NameID: "Krisan", Description: "Simbol panjang umur dan keberuntungan.", Emoji: "💮", Price: 8000, IsAvailable: true, Stock: 150, Colors: []string{"#FFF176", "#FFFFFF", "#FF8F00"}, Occasions: []string{"birthday", "anniversary", "congratulations"}, Meaning: "Panjang umur dan keberuntungan"},
		{ID: "baby_breath", Name: "Baby's Breath", NameID: "Baby Breath", Description: "Bunga pengisi yang cantik dan lembut.", Emoji: "🌬️", Price: 5000, IsAvailable: true, Stock: 200, Colors: []string{"#FFFFFF"}, Occasions: []string{"wedding", "graduation", "birthday", "anniversary"}, Meaning: "Ketulusan dan kemurnian"},
		{ID: "carnation_red", Name: "Red Carnation", NameID: "Anyelir Merah", Description: "Melambangkan cinta yang dalam.", Emoji: "🌺", Price: 10000, IsAvailable: true, Stock: 80, Colors: []string{"#B71C1C"}, Occasions: []string{"mothers_day", "valentines", "anniversary"}, Meaning: "Cinta dan kekaguman"},
		{ID: "carnation_pink", Name: "Pink Carnation", NameID: "Anyelir Pink", Description: "Melambangkan kasih sayang seorang ibu.", Emoji: "🌸", Price: 10000, IsAvailable: true, Stock: 80, Colors: []string{"#EC407A"}, Occasions: []string{"mothers_day", "birthday"}, Meaning: "Kasih sayang dan kelembutan"},
		{ID: "orchid_purple", Name: "Purple Orchid", NameID: "Anggrek Ungu", Description: "Mewah dan eksotis.", Emoji: "🌸", Price: 25000, IsAvailable: false, Stock: 0, Colors: []string{"#9C27B0"}, Occasions: []string{"wedding", "anniversary", "graduation"}, Meaning: "Kemewahan dan keistimewaan"},
		{ID: "lavender", Name: "Lavender", NameID: "Lavender", Description: "Harum dan menenangkan.", Emoji: "💜", Price: 12000, IsAvailable: true, Stock: 90, Colors: []string{"#B39DDB"}, Occasions: []string{"wedding", "birthday", "anniversary"}, Meaning: "Ketenangan dan cinta murni"},
		{ID: "daisy", Name: "Daisy", NameID: "Aster", Description: "Ceria dan penuh semangat.", Emoji: "🌼", Price: 7000, IsAvailable: true, Stock: 130, Colors: []string{"#FFFFFF", "#FFF59D"}, Occasions: []string{"birthday", "congratulations"}, Meaning: "Keceriaan dan semangat"},
		{ID: "peony", Name: "Peony", NameID: "Peony", Description: "Mewah dan romantis.", Emoji: "🌹", Price: 35000, IsAvailable: true, Stock: 30, Colors: []string{"#F48FB1", "#FFCCBC"}, Occasions: []string{"wedding", "anniversary"}, Meaning: "Kemakmuran dan kebahagiaan"},
	}

	for i := range flowers {
		DB.Create(&flowers[i])
	}
	log.Printf("✅ %d bunga berhasil di-seed", len(flowers))
}

func seedCatalog() {
	var count int64
	DB.Model(&models.CatalogBouquetDB{}).Count(&count)
	if count > 0 {
		return
	}
	log.Println("🛍️ Seeding katalog bouquet...")

	catalogs := []models.CatalogBouquetDB{
		{ID: "cat_wisuda_classic", Name: "Wisuda Classic", Description: "Bouquet classic untuk wisuda dengan mawar pink dan baby breath yang elegan.", ImageURL: "", Style: "Classic", Occasion: "graduation", Price: 185000, StemCount: 15, IsAvailable: true, Stock: 10, SortOrder: 1},
		{ID: "cat_wisuda_premium", Name: "Wisuda Premium", Description: "Bouquet premium wisuda dengan kombinasi tulip ungu dan mawar putih.", ImageURL: "", Style: "Premium", Occasion: "graduation", Price: 350000, StemCount: 25, IsAvailable: true, Stock: 5, SortOrder: 2},
		{ID: "cat_valentine_romantic", Name: "Valentine Romantic", Description: "Bouquet penuh cinta dengan mawar merah dan anyelir merah yang menawan.", ImageURL: "", Style: "Romantic", Occasion: "valentines", Price: 220000, StemCount: 20, IsAvailable: true, Stock: 8, SortOrder: 3},
		{ID: "cat_birthday_sunshine", Name: "Birthday Sunshine", Description: "Rangkaian bunga matahari dan aster yang ceria untuk ulang tahun.", ImageURL: "", Style: "Playful", Occasion: "birthday", Price: 160000, StemCount: 12, IsAvailable: true, Stock: 15, SortOrder: 4},
		{ID: "cat_wedding_elegant", Name: "Wedding Elegance", Description: "Bouquet pernikahan mewah dengan peony, lili putih dan mawar putih.", ImageURL: "", Style: "Elegant", Occasion: "wedding", Price: 550000, StemCount: 30, IsAvailable: true, Stock: 5, SortOrder: 5},
		{ID: "cat_ibu_kasih", Name: "Kasih Ibu", Description: "Hadiah penuh kasih untuk ibu tercinta — anyelir pink dan mawar pink.", ImageURL: "", Style: "Warm", Occasion: "mothers_day", Price: 175000, StemCount: 15, IsAvailable: true, Stock: 12, SortOrder: 6},
	}

	for i := range catalogs {
		DB.Create(&catalogs[i])
	}
	log.Printf("✅ %d katalog berhasil di-seed", len(catalogs))
}
