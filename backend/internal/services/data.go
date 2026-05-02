package services

import (
	"bouquet-app/internal/database"
	"bouquet-app/internal/models"
)

// GetAllBouquetTypes — tipe bouquet tetap hardcoded (jarang berubah)
func GetAllBouquetTypes() []models.BouquetType {
	return []models.BouquetType{
		{ID: "graduation", Name: "Wisuda", Description: "Rayakan momen kelulusan dengan bouquet spesial", Icon: "🎓", Theme: "#7C4DFF"},
		{ID: "wedding", Name: "Pernikahan", Description: "Bouquet romantis untuk hari pernikahan impian", Icon: "💍", Theme: "#E91E8C"},
		{ID: "birthday", Name: "Ulang Tahun", Description: "Ucapkan selamat ulang tahun dengan warna-warni bunga", Icon: "🎂", Theme: "#FF6B35"},
		{ID: "anniversary", Name: "Anniversary", Description: "Peringati momen spesial bersama orang tersayang", Icon: "💑", Theme: "#C62828"},
		{ID: "sympathy", Name: "Dukacita", Description: "Sampaikan rasa duka dengan penuh ketulusan", Icon: "🕊️", Theme: "#546E7A"},
		{ID: "congratulations", Name: "Selamat", Description: "Rayakan pencapaian luar biasa seseorang", Icon: "🏆", Theme: "#F9A825"},
		{ID: "valentines", Name: "Valentine", Description: "Ungkapkan cinta dengan rangkaian bunga indah", Icon: "❤️", Theme: "#D32F2F"},
		{ID: "mothers_day", Name: "Hari Ibu", Description: "Hargai kasih sayang ibu dengan bunga istimewa", Icon: "🌸", Theme: "#AD1457"},
	}
}

// GetAllFlowers — ambil dari PostgreSQL
func GetAllFlowers() []models.Flower {
	var records []models.FlowerDB
	database.DB.Order("name_id ASC").Find(&records)

	flowers := make([]models.Flower, len(records))
	for i, r := range records {
		flowers[i] = r.ToFlower()
	}
	return flowers
}

// GetFlowerByID — ambil satu bunga dari DB
func GetFlowerByID(id string) (*models.Flower, bool) {
	var record models.FlowerDB
	if err := database.DB.First(&record, "id = ?", id).Error; err != nil {
		return nil, false
	}
	f := record.ToFlower()
	return &f, true
}

// GetAllCatalog — ambil semua katalog dari PostgreSQL
func GetAllCatalog() []models.CatalogBouquetDB {
	var records []models.CatalogBouquetDB
	database.DB.Where("is_available = true").Order("sort_order ASC").Find(&records)
	return records
}

// GetAllCatalogAdmin — semua katalog tanpa filter (untuk admin)
func GetAllCatalogAdmin() []models.CatalogBouquetDB {
	var records []models.CatalogBouquetDB
	database.DB.Order("sort_order ASC").Find(&records)
	return records
}
