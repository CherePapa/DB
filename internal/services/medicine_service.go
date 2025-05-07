package services

import (
	"DB-project/internal/models"
	"DB-project/internal/repositories"
)

type MedicineService interface {
	AddMedicine(medicine models.Medicine) error
	GetAllMedicines() ([]models.Medicine, error)
	GetMedicineByID(id uint) (*models.Medicine, error)
	UpdateMedicine(medicine models.Medicine) error
	DeleteMedicine(id uint) error
}

type medicineService struct {
	repo repositories.MedicineRepository
}

func NewMedicineService(repo repositories.MedicineRepository) MedicineService {
	return &medicineService{repo: repo}
}

func (s *medicineService) AddMedicine(medicine models.Medicine) error {
	return s.repo.Create(&medicine)
}

func (s *medicineService) GetAllMedicines() ([]models.Medicine, error) {
	return s.repo.FindAll()
}

func (s *medicineService) GetMedicineByID(id uint) (*models.Medicine, error) {
	return s.repo.FindByID(id)
}

func (s *medicineService) UpdateMedicine(medicine models.Medicine) error {
	return s.repo.Update(&medicine)
}

func (s *medicineService) DeleteMedicine(id uint) error {
	// Добавьте предварительную проверку существования записи
	if _, err := s.repo.FindByID(id); err != nil {
		return err
	}

	return s.repo.Delete(id)
}
