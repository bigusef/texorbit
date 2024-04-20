package user

import (
	"github.com/bigusef/texorbit/internal/database"
	"github.com/google/uuid"
	"time"
)

type listStaff struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	JoinDate time.Time `json:"join_date"`
}

type newStaff struct {
	Name        string `json:"name" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phone_number" validate:"phone_number"`
}

type staffInfo struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	Avatar      string    `json:"avatar"`
	Status      string    `json:"status"`
	JoinDate    time.Time `json:"joinDate"`
	LastLogin   time.Time `json:"lastLogin"`
}

type updateStaff struct {
	Name        string                 `json:"name"`
	Email       string                 `json:"email" validate:"required,email"`
	PhoneNumber string                 `json:"phone_number"`
	Status      database.AccountStatus `json:"status" validate:"required"`
}
