package postgres

import (
	"database/sql"
	"fmt"

	"github.com/internships-backend/test-backend-beldurad/internal/config"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/tanimutomo/sqlfile"
)

func formatToDSN(cfg *pq.Config) string {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Database,
		cfg.SSLMode,
	)
	return dsn
}

func New(dbConfig config.DatabaseConfig) (*sql.DB, error) {
	const op = "database.postgres.New"

	cfg := pq.Config{
		Host:     dbConfig.Host,
		Port:     dbConfig.Port,
		User:     dbConfig.User,
		Password: dbConfig.Password,
		Database: dbConfig.DatabaseName,
		SSLMode:  "disable",
	}

	db, err := sqlx.Open("postgres", formatToDSN(&cfg))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	sf := sqlfile.New()

	if err := sf.File(dbConfig.InitSqlFilepath); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if _, err := sf.Exec(db.DB); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return db.DB, nil
}
