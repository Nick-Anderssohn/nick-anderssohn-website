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
}

const (
	createFilesTableIFNotExistsSql = `
CREATE TABLE IF NOT EXISTS files(
  code VARCHAR(6) PRIMARY KEY,
  isCompressed BOOLEAN DEFAULT FALSE,
  compressionType VARCHAR(15) DEFAULT '',
  path VARCHAR(50) NOT NULL
);
`

	insertIntoFilesTableSQL = `
INSERT INTO files
VALUES ($1, $2, $3, $4);
`

	selectFileInfoSql = `
SELECT * FROM files WHERE code = $1;
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

func InsertNewFileInfo(code, compressionType, path string, isCompressed bool) (err error) {
	if _, err = conn.Exec(insertIntoFilesTableSQL, code, isCompressed, compressionType, path); err != nil {
		log.Println(err.Error())
	}
	return
}

func SelectFileInfo(code string) *FileInfo {
	var fileInfo FileInfo
	row := conn.QueryRow(selectFileInfoSql, code)
	if err := row.Scan(&fileInfo.Code); err != nil {
		log.Println(err.Error())
		return nil
	}
	if err := row.Scan(&fileInfo.IsCompressed); err != nil {
		log.Println(err.Error())
		return nil
	}
	if err := row.Scan(&fileInfo.CompressionType); err != nil {
		log.Println(err.Error())
		return nil
	}
	if err := row.Scan(&fileInfo.Path); err != nil {
		log.Println(err.Error())
		return nil
	}
	return &fileInfo
}

func FileInfoExists(code string) bool {
	var dontCare string
	err := conn.QueryRow(selectFileInfoSql, code).Scan(dontCare)
	return err != sql.ErrNoRows
}
