package cmd

import (
	"anybakup/util"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type FileOperation struct {
	ID         int64
	SrcFile    string
	DestFile   string
	IsFile     bool
	AddTime    time.Time
	UpdateTime time.Time
}

type sqldb struct {
	db     *sql.DB
	dbfile string
}

func NewSqldb() (*sqldb, error) {
	conf := util.Config{}
	configroot, err := conf.Configdir()
	if err != nil {
		return nil, fmt.Errorf("git repo %v", err)
	}
	dbPath := filepath.Join(configroot, "file_operations.db") // Default database file path

	// Create directory if it doesn't exist
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory for database: %v", err)
	}

	// Open database connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Create table if it doesn't exist
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS file_operations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		srcfile TEXT NOT NULL,
		destfile TEXT NOT NULL,
		isfile BOOLEAN NOT NULL,
		add_time DATETIME DEFAULT CURRENT_TIMESTAMP,
		update_time DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	if _, err := db.Exec(createTableSQL); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create table: %v", err)
	}

	s := &sqldb{
		db:     db,
		dbfile: dbPath,
	}

	return s, nil
}

func (s *sqldb) Close() error {
	return s.db.Close()
}

func BackupOptAdd(srcFile, destFile string, isFile bool) error {
	db, err := NewSqldb()
	if err != nil {
		return err
	}
	defer db.Close()

	query := `
	INSERT INTO file_operations (srcfile, destfile, isfile, add_time, update_time)
	VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

	_, err = db.db.Exec(query, srcFile, destFile, isFile)
	if err != nil {
		return fmt.Errorf("failed to insert file operation: %v", err)
	}

	return nil
}

func BackupOptRm(file string) error {
	db, err := NewSqldb()
	if err != nil {
		return err
	}
	defer db.Close()

	// Remove entries where either srcfile or destfile matches the given file
	query := `DELETE FROM file_operations WHERE srcfile = ? OR destfile = ?`

	_, err = db.db.Exec(query, file, file)
	if err != nil {
		return fmt.Errorf("failed to delete file operation: %v", err)
	}

	return nil
}
