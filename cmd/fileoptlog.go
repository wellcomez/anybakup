package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"anybakup/util"

	_ "github.com/mattn/go-sqlite3"
)

type FileOperation struct {
	ID         int64
	SrcFile    string
	DestFile   string
	IsFile     bool
	RevCount   int
	Sub        bool
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
	if err := os.MkdirAll(dir, 0o755); err != nil {
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
		revcount INTEGER DEFAULT 0,
		sub BOOLEAN DEFAULT FALSE,
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

func BakupOptAdd(srcFile string, destFile util.RepoPath, isFile bool, sub bool,g GitCmd) error {
	revcount := 0
	if isFile {
		if count, err := g.GetFileLog(srcFile); err != nil {
			return err
		} else {
			revcount = len(count)
		}
	} else {
		revcount = 1
	}
	db, err := NewSqldb()
	if err != nil {
		return err
	}
	defer db.Close()

	// Check if the entry already exists
	checkQuery := `
	SELECT COUNT(*) FROM file_operations
	WHERE srcfile = ?`

	var count int
	err = db.db.QueryRow(checkQuery, srcFile).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check existing file operation: %v", err)
	}

	if count > 0 {
		// Update the existing entry with a new update_time and revcount
		updateQuery := `
		UPDATE file_operations
		SET destfile = ?, isfile = ?, revcount = ?, sub = ?, update_time = CURRENT_TIMESTAMP
		WHERE srcfile = ?`

		_, err = db.db.Exec(updateQuery, destFile, isFile, revcount, sub, srcFile)
		if err != nil {
			return fmt.Errorf("failed to update file operation: %v", err)
		}
	} else {
		// Insert a new entry
		insertQuery := `
		INSERT INTO file_operations (srcfile, destfile, isfile, revcount, sub, add_time, update_time)
		VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

		_, err = db.db.Exec(insertQuery, srcFile, destFile, isFile, revcount, sub)
		if err != nil {
			return fmt.Errorf("failed to insert file operation: %v", err)
		}
	}

	return nil
}

func BakupOptRm(file util.RepoPath) error {
	db, err := NewSqldb()
	if err != nil {
		return err
	}
	defer db.Close()

	// Check if file exists in database
	var count int
	checkQuery := `SELECT COUNT(*) FROM file_operations WHERE destfile = ?`
	if err := db.db.QueryRow(checkQuery, file).Scan(&count); err == nil {
		if count == 0 {
			fmt.Printf("File %s not found in database, no deletion needed", file)
			return nil
		}
	}

	// Remove entries where either srcfile or destfile matches the given file
	query := `DELETE FROM file_operations WHERE  destfile = ?`

	_, err = db.db.Exec(query, file)
	// r.RowsAffected()
	if err != nil {
		return fmt.Errorf("sql failed to delete file operation:  %v", err)
	} else {
		fmt.Printf("SQL Deleted file operation for %s\n", file)
	}
	return nil
}

func GetAllOpt() ([]FileOperation, error) {
	db, err := NewSqldb()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `SELECT id, srcfile, destfile, isfile, revcount, sub, add_time, update_time FROM file_operations ORDER BY add_time DESC`

	rows, err := db.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query file operations: %v", err)
	}
	defer rows.Close()

	var operations []FileOperation
	for rows.Next() {
		var op FileOperation
		err := rows.Scan(&op.ID, &op.SrcFile, &op.DestFile, &op.IsFile, &op.RevCount, &op.Sub, &op.AddTime, &op.UpdateTime)
		if err != nil {
			return nil, fmt.Errorf("failed to scan file operation: %v", err)
		}
		operations = append(operations, op)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	return operations, nil
}
