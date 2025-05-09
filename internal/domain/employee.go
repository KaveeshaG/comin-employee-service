package domain

import (
	"time"

	"github.com/google/uuid"
)

type Employee struct {
	Base
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;not null"`
	OrganizationID  uuid.UUID  `json:"organization_id" gorm:"type:uuid;not null"`
	PositionID      *uuid.UUID `json:"position_id,omitempty" gorm:"type:uuid"`
	DepartmentID    *uuid.UUID `json:"department_id,omitempty" gorm:"type:uuid"`
	FirstName       string     `json:"first_name" gorm:"not null"`
	LastName        string     `json:"last_name" gorm:"not null"`
	Email           string     `json:"email" gorm:"not null"`
	Phone           string     `json:"phone"`
	DateOfBirth     *time.Time `json:"date_of_birth"`
	HireDate        time.Time  `json:"hire_date" gorm:"not null"`
	TerminationDate *time.Time `json:"termination_date"`
	EmployeeID      string     `json:"employee_id"`
	Status          string     `json:"status" gorm:"default:'active'"`
	ManagerID       *uuid.UUID `json:"manager_id,omitempty" gorm:"type:uuid"`
	WorkType        string     `json:"work_type" gorm:"default:'full-time'"`
	Position        *Position  `json:"position,omitempty" gorm:"foreignKey:PositionID"`
}

type Position struct {
	Base
	OrganizationID uuid.UUID  `json:"organization_id" gorm:"type:uuid;not null"`
	Title          string     `json:"title" gorm:"not null"`
	Description    string     `json:"description"`
	Level          int        `json:"level"`
	DepartmentID   *uuid.UUID `json:"department_id,omitempty" gorm:"type:uuid"`
	Status         string     `json:"status" gorm:"default:'active'"`
}

type EmployeeContract struct {
	Base
	EmployeeID   uuid.UUID  `json:"employee_id" gorm:"type:uuid"`
	StartDate    time.Time  `json:"start_date"`
	EndDate      *time.Time `json:"end_date"`
	ContractType string     `json:"contract_type"`
	SalaryAmount float64    `json:"salary_amount"`
	Currency     string     `json:"currency"`
	Status       string     `json:"status" gorm:"default:'active'"`
}

type EmployeeDocument struct {
	Base
	EmployeeID   uuid.UUID  `json:"employee_id" gorm:"type:uuid"`
	DocumentType string     `json:"document_type"`
	DocumentName string     `json:"document_name"`
	DocumentURL  string     `json:"document_url"`
	ExpiryDate   *time.Time `json:"expiry_date"`
	Status       string     `json:"status" gorm:"default:'active'"`
}

type CreateEmployeeRequest struct {
	ID             uuid.UUID  `json:"id"`
	OrganizationID uuid.UUID  `json:"organization_id" binding:"required"`
	PositionID     *uuid.UUID `json:"position_id"`
	DepartmentID   *uuid.UUID `json:"department_id"`
	FirstName      string     `json:"first_name" binding:"required"`
	LastName       string     `json:"last_name" binding:"required"`
	Email          string     `json:"email" binding:"required,email"`
	Phone          string     `json:"phone"`
	DateOfBirth    *time.Time `json:"date_of_birth"`
	HireDate       time.Time  `json:"hire_date" binding:"required"`
	EmployeeID     string     `json:"employee_id"`
	WorkType       string     `json:"work_type"`
	ManagerID      *uuid.UUID `json:"manager_id"`
	LeaveTypeID    *uuid.UUID `json:"leave_type_id"`
	LeaveTotalDays *float64   `json:"leave_total_days"`
}

type UpdateEmployeeRequest struct {
	PositionID   *uuid.UUID `json:"position_id"`
	DepartmentID *uuid.UUID `json:"department_id"`
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	Email        string     `json:"email" binding:"required,email"`
	Phone        string     `json:"phone"`
	ManagerID    *uuid.UUID `json:"manager_id"`
	WorkType     string     `json:"work_type"`
	Status       string     `json:"status,omitempty"`
	DateOfBirth  *time.Time `json:"date_of_birth"`
	HireDate     time.Time  `json:"hire_date" binding:"required"`
	EmployeeID   string     `json:"employee_id"`
}
