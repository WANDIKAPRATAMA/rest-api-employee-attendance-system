// department_repository.go
package repository

import (
	"employee-attendance-system/internal/entity/domain"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type DepartmentRepository interface {
	CreateDepartment(dept *domain.Department) error
	FindDepartmentByID(id uuid.UUID) (*domain.Department, error)
	UpdateDepartment(dept *domain.Department) error
	DeleteDepartment(id uuid.UUID) error
	FindAllDepartments(offset, limit int) ([]*domain.Department, int64, error)
}

type departmentRepository struct {
	db  *gorm.DB
	log *logrus.Logger
}

func NewDepartmentRepository(db *gorm.DB, log *logrus.Logger) DepartmentRepository {
	return &departmentRepository{db: db, log: log}
}

func (r *departmentRepository) CreateDepartment(dept *domain.Department) error {
	return r.db.Create(dept).Error
}

func (r *departmentRepository) FindDepartmentByID(id uuid.UUID) (*domain.Department, error) {
	var dept domain.Department
	err := r.db.First(&dept, id).Error
	if err != nil {
		return nil, err
	}
	return &dept, nil
}

func (r *departmentRepository) UpdateDepartment(dept *domain.Department) error {
	return r.db.Save(dept).Error
}

func (r *departmentRepository) DeleteDepartment(id uuid.UUID) error {
	return r.db.Delete(&domain.Department{}, id).Error
}

func (r *departmentRepository) FindAllDepartments(offset, limit int) ([]*domain.Department, int64, error) {
	var depts []*domain.Department
	var total int64
	if err := r.db.Model(&domain.Department{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := r.db.Offset(offset).Limit(limit).Find(&depts).Error
	return depts, total, err
}
