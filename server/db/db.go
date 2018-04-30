package db

import (
	"database/sql"
	"fmt"

	"log"

	_ "github.com/lib/pq"
)

// TODO: read username and password from config file
const (
	postgresUserName = "postgres"
	postgresPassword = "123"
	dbName           = "full_share"
)

type FileInfo struct {
	Code            string
	IsCompressed    bool
	CompressionType string
	Path            string
	FileSize        int
}

const (
	createFilesTableIFNotExistsSql = `
CREATE TABLE IF NOT EXISTS files(
  code VARCHAR(6) PRIMARY KEY,
  isCompressed BOOLEAN DEFAULT FALSE,
  compressionType VARCHAR(15) DEFAULT '',
  path VARCHAR(50) NOT NULL,
  fileSize INTEGER DEFAULT 0
);
`

	insertIntoFilesTableSQL = `
INSERT INTO files
VALUES ($1, $2, $3, $4, $5);
`

	selectFileInfoSql = `
SELECT * FROM files WHERE code = $1;
`

	countFileWithCode = `
SELECT COUNT(*) FROM files WHERE code = $1;
`
)

var conn *sql.DB

func init() {
	ConnectToDb()
	CreateTablesIfNotExists()
}

func ConnectToDb() {
	var err error
	connStr := fmt.Sprintf("postgres://%s:%s@localhost/%s?sslmode=disable", postgresUserName, postgresPassword, dbName)
	conn, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Println(err.Error())
	}
}

func CreateTablesIfNotExists() {
	if _, err := conn.Exec(createFilesTableIFNotExistsSql); err != nil {
		log.Println(err.Error())
	}
}

func InsertNewFileInfo(code, compressionType, path string, isCompressed bool, fileSize int) (err error) {
	if _, err = conn.Exec(insertIntoFilesTableSQL, code, isCompressed, compressionType, path, fileSize); err != nil {
		log.Println(err.Error())
	}
	return
}

func SelectFileInfo(code string) *FileInfo {
	var fileInfo FileInfo
	row := conn.QueryRow(selectFileInfoSql, code)
	if err := row.Scan(&fileInfo.Code, &fileInfo.IsCompressed, &fileInfo.CompressionType, &fileInfo.Path, &fileInfo.FileSize); err != nil {
		log.Println(err.Error())
		return nil
	}
	return &fileInfo
}

func FileInfoExists(code string) bool {
	var count int
	conn.QueryRow(countFileWithCode, code).Scan(&count)
	return count >= 1
}
