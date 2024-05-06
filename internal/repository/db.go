package repository

import (
	"database/sql"
	"fmt"
)

type Repository struct {
	Dbh *sql.DB
}

func (repo *Repository) Insert(name string, size int) error {
	var q = "INSERT  INTO fileinfo(file_name, file_size) VALUES(?,?) ON DUPLICATE KEY UPDATE file_size=?"
	_, err := repo.Dbh.Exec(q, name, size, size)

	if err != nil {
		return fmt.Errorf("insert error %w", err)
	}
	return nil
}
