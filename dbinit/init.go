package dbinit

import (
	"database/sql"
	"log"
	"os"
)

// InitDB выполняет SQL-скрипт для инициализации базы данных
func InitDB(db *sql.DB, scriptPath string) error {
	file, err := os.ReadFile(scriptPath)
	if err != nil {
		return err
	}

	query := string(file)
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	log.Println("Database initialized successfully")
	return nil
}
