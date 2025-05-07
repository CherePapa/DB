package main

import (
	"DB-project/internal/cli/handlers"
	"DB-project/internal/config"
	"DB-project/internal/models"
	"DB-project/internal/repositories"
	"DB-project/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация БД
	db, err := config.ConnectDB()
	if err != nil {
		panic(err)
	}

	// Миграции
	db.AutoMigrate(&models.Medicine{})

	// Инициализация зависимостей
	medicineRepo := repositories.NewMedicineRepository(db)
	medicineService := services.NewMedicineService(medicineRepo)
	medicineHandler := handlers.NewMedicineHandler(medicineService)

	// Создание Gin-приложения
	router := gin.Default()

	// Загрузка шаблонов и статических файлов
	router.LoadHTMLGlob("../internal/web/templates/*")
	router.Static("/assets", "../internal/web/assets")

	// Маршруты
	router.GET("/", medicineHandler.ListMedicines)
	router.GET("/add", medicineHandler.ShowAddForm)
	router.POST("/add", medicineHandler.AddMedicine)
	router.GET("/edit/:id", medicineHandler.ShowEditForm)
	router.POST("/edit/:id", medicineHandler.UpdateMedicine)
	router.POST("/delete/:id", medicineHandler.DeleteMedicine)

	// Запуск сервера
	router.Run(":8080")
}
