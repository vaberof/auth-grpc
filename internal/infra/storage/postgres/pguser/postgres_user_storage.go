package pguser

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/vaberof/auth-grpc/internal/domain/user"
	"github.com/vaberof/auth-grpc/internal/infra/storage"
	"github.com/vaberof/auth-grpc/pkg/domain"
)

type PgUserStorage struct {
	db *sqlx.DB
}

func NewPgUserStorage(db *sqlx.DB) *PgUserStorage {
	return &PgUserStorage{
		db: db,
	}
}

func (us *PgUserStorage) Create(email domain.Email, password domain.Password) (domain.UserId, error) {
	query := `
			INSERT INTO users(
			                  email,
			                  password
			) VALUES ($1, $2)
			RETURNING id
	`

	row := us.db.QueryRow(query, email.String(), password.String())

	var uid int64

	err := row.Scan(&uid)
	if err != nil {
		return 0, nil
	}

	return domain.UserId(uid), nil
}

func (us *PgUserStorage) GetByEmail(email domain.Email) (*user.User, error) {
	query := `
			SELECT * FROM users
			WHERE email=$1
	`

	row := us.db.QueryRow(query, email)

	var pgUser User

	err := row.Scan(
		&pgUser.Id,
		&pgUser.Email,
		&pgUser.Password,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrPostgresUserNotFound
		}
		return nil, err
	}

	return toDomainUser(&pgUser), nil
}

func (us *PgUserStorage) ExistsByEmail(email domain.Email) (bool, error) {
	query := `
			SELECT id FROM users
			WHERE email=$1
	`

	var uid int64

	err := us.db.QueryRow(query, email).Scan(&uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
