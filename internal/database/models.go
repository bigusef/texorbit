// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package database

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type AccountStatus string

const (
	AccountStatusActive    AccountStatus = "active"
	AccountStatusSuspended AccountStatus = "suspended"
	AccountStatusDeleted   AccountStatus = "deleted"
)

func (e *AccountStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = AccountStatus(s)
	case string:
		*e = AccountStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for AccountStatus: %T", src)
	}
	return nil
}

type NullAccountStatus struct {
	AccountStatus AccountStatus
	Valid         bool // Valid is true if AccountStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullAccountStatus) Scan(value interface{}) error {
	if value == nil {
		ns.AccountStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.AccountStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullAccountStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.AccountStatus), nil
}

type City struct {
	ID       int64
	NameEn   string
	NameAr   string
	IsActive bool
}

type User struct {
	ID          uuid.UUID
	Name        string
	Email       string
	PhoneNumber pgtype.Text
	Avatar      pgtype.Text
	Status      AccountStatus
	IsStaff     bool
	JoinDate    pgtype.Timestamptz
	LastLogin   pgtype.Timestamptz
}
