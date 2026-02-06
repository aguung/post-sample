package database

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var (
	ErrDuplicateKey        = errors.New("duplicate key value violates unique constraint")
	ErrForeignKeyViolation = errors.New("foreign key violation")
	ErrRecordNotFound      = errors.New("record not found")
)

// ParseError checks the error type and returns a standardized error
func ParseError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrRecordNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return ErrDuplicateKey
		case "23503":
			return ErrForeignKeyViolation
		}
	}
	return err
}
