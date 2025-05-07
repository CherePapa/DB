package repositories

import (
	"DB-project/internal/models"

	"gorm.io/gorm"
)

type MedicineRepository interface {
	Create(medicine *models.Medicine) error
	FindAll() ([]models.Medicine, error)
	FindByID(id uint) (*models.Medicine, error)
	Update(medicine *models.Medicine) error
	Delete(id uint) error
}

type medicineRepo struct {
	db *gorm.DB
}

func NewMedicineRepository(db *gorm.DB) MedicineRepository {
	return &medicineRepo{db: db}
}

func (r *medicineRepo) Create(medicine *models.Medicine) error {
	return r.db.Create(medicine).Error
}

func (r *medicineRepo) FindAll() ([]models.Medicine, error) {
	var medicines []models.Medicine
	err := r.db.Find(&medicines).Error
	return medicines, err
}

func (r *medicineRepo) FindByID(id uint) (*models.Medicine, error) {
	var medicine models.Medicine
	err := r.db.First(&medicine, id).Error
	return &medicine, err
}

func (r *medicineRepo) Update(medicine *models.Medicine) error {
	return r.db.Save(medicine).Error
}

func (r *medicineRepo) Delete(id uint) error {
	result := r.db.Unscoped().Delete(&models.Medicine{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
