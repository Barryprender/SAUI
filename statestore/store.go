package statestore

import (
	"database/sql"
	"log/slog"

	_ "modernc.org/sqlite"
)

type Store struct {
	db     *sql.DB
	logger *slog.Logger
}

func New(path string, logger *slog.Logger) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(`PRAGMA journal_mode=WAL`); err != nil {
		db.Close()
		return nil, err
	}
	if _, err := db.Exec(`PRAGMA foreign_keys=ON`); err != nil {
		db.Close()
		return nil, err
	}

	s := &Store{db: db, logger: logger}
	if err := s.migrate(); err != nil {
		db.Close()
		return nil, err
	}
	if err := s.migrateFoodDemo(); err != nil {
		db.Close()
		return nil, err
	}
	if err := s.migrateBankingDemo(); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}
