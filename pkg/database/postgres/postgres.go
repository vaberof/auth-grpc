package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
}

type ManagedDatabase struct {
	Db *sqlx.DB
}

func New(config *Config) (*ManagedDatabase, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Host, config.Port, config.User, config.Password, config.Database)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	managedDatabase := &ManagedDatabase{
		Db: db,
	}

	return managedDatabase, nil
}
