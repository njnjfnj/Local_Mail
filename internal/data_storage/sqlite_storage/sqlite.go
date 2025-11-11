package sqlite_storage

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Repository struct {
	db *sql.DB
}

func New(dbPath string) (*Repository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("can not open db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can not connect to db: %w", err)
	}

	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("can not create tables: %w", err)
	}

	return &Repository{db: db}, nil
}

func createTables(db *sql.DB) error {
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS messages (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        from_to TEXT NOT NULL,
        message TEXT,
        file_path TEXT,
        image_path TEXT,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`

	_, err := db.Exec(createTableSQL)
	return err
}

// TODO: если надо будет сделать сохранение чата. Возможно добавить
// специальный тумблер в чате чтоб можно было сохранять сообщения
