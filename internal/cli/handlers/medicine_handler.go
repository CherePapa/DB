package handlers

import (
	"DB-project/internal/models"
	"DB-project/internal/services"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MedicineWebHandler struct {
	service services.MedicineService
}

func NewMedicineHandler(service services.MedicineService) *MedicineWebHandler {
	return &MedicineWebHandler{service: service}
}

func (h *MedicineWebHandler) ListMedicines(c *gin.Context) {
	medicines, err := h.service.GetAllMedicines()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Failed to fetch medicines",
		})
		return
	}

	c.HTML(http.StatusOK, "list.html", gin.H{
		"medicines": medicines,
	})
}

func (h *MedicineWebHandler) ShowAddForm(c *gin.Context) {
	// Передаем текущую дату в формате YYYY-MM-DD
	currentDate := time.Now().Format("2006-01-02")

	c.HTML(http.StatusOK, "add.html", gin.H{
		"currentDate": currentDate,
	})
}

func (h *MedicineWebHandler) AddMedicine(c *gin.Context) {
	var medicine models.Medicine

	// Парсинг данных формы
	medicine.Name = c.PostForm("name")
	medicine.Manufacturer = c.PostForm("manufacturer")
	medicine.Price, _ = strconv.ParseFloat(c.PostForm("price"), 64)
	medicine.Quantity, _ = strconv.Atoi(c.PostForm("quantity"))
	expiryDate, _ := time.Parse("2006-01-02", c.PostForm("expiryDate"))
	medicine.ExpiryDate = expiryDate

	// Валидация
	if medicine.Name == "" || medicine.Manufacturer == "" {
		c.HTML(http.StatusBadRequest, "add.html", gin.H{
			"error": "Все обязательные поля должны быть заполнены",
		})
		return
	}

	if err := h.service.AddMedicine(medicine); err != nil {
		c.HTML(http.StatusInternalServerError, "add.html", gin.H{
			"error": "Ошибка сохранения: " + err.Error(),
		})
		return
	}

	c.Redirect(http.StatusFound, "/")
}

func (h *MedicineWebHandler) ShowEditForm(c *gin.Context) {
	idStr := c.Param("id") // Получаем ID из URL
	id, _ := strconv.Atoi(idStr)

	medicine, err := h.service.GetMedicineByID(uint(id))
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Medicine not found"})
		return
	}

	c.HTML(http.StatusOK, "edit.html", gin.H{
		"medicine": medicine,
	})
}

func (h *MedicineWebHandler) DeleteMedicine(c *gin.Context) {
	// Получаем и валидируем ID
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing medicine ID"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid medicine ID format"})
		return
	}

	// Логирование попытки удаления
	log.Printf("Attempting to delete medicine with ID: %d", id)

	// Выполняем удаление
	if err := h.service.DeleteMedicine(uint(id)); err != nil {
		log.Printf("Delete failed: %v", err)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Medicine not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete medicine"})
		}
		return
	}

	log.Printf("Successfully deleted medicine with ID: %d", id)
	c.JSON(http.StatusOK, gin.H{"message": "Medicine deleted successfully"})
}

// handlers/medicine_handler.go
func (h *MedicineWebHandler) UpdateMedicine(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	medicine, err := h.service.GetMedicineByID(uint(id))
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Лекарство не найдено"})
		return
	}

	medicine.Name = c.PostForm("name")
	medicine.Manufacturer = c.PostForm("manufacturer")
	medicine.Price, _ = strconv.ParseFloat(c.PostForm("price"), 64)
	medicine.Quantity, _ = strconv.Atoi(c.PostForm("quantity"))
	expiryDate, _ := time.Parse("2006-01-02", c.PostForm("expiryDate"))
	medicine.ExpiryDate = expiryDate

	if err := h.service.UpdateMedicine(*medicine); err != nil {
		c.HTML(http.StatusInternalServerError, "edit.html", gin.H{
			"error":    "Ошибка обновления: " + err.Error(),
			"medicine": medicine,
		})
		return
	}

	c.Redirect(http.StatusFound, "/")
}
