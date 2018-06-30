package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"io/ioutil"
	"encoding/json"
	"time"
)

type dbConfig struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
	DbName string `json:"DbName"`
}

var config dbConfig

const dbConfigPath = "db_config.json"

type FileInfo struct {
	Code       string
	Name       string
	FileSize   int
	UploadedOn *time.Time
}

const (
	createFilesTableIFNotExistsSql = `
CREATE TABLE IF NOT EXISTS files(
  code VARCHAR(36) PRIMARY KEY,
  name VARCHAR(50) NOT NULL,
  fileSize INTEGER DEFAULT 0,
  uploadedOn DATE 
);
`

	insertIntoFilesTableSQL = `
INSERT INTO files
VALUES ($1, $2, $3);
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
	readDBConfig()
	ConnectToDb()
	CreateTablesIfNotExists()
}

func readDBConfig() {
	configBytes, _ := ioutil.ReadFile(dbConfigPath)
	json.Unmarshal(configBytes, &config)
	if config.DbName == "" {
		panic("no database name")
	}
}

func ConnectToDb() {
	connStr := fmt.Sprintf("postgres://%s:%s@localhost/%s?sslmode=disable", config.Username, config.Password, config.DbName)
	conn, _ = sql.Open("postgres", connStr)
}

func CreateTablesIfNotExists() {
	conn.Exec(createFilesTableIFNotExistsSql)
}

func InsertNewFileInfo(code, name string, fileSize int) (err error) {
	now := time.Now().UTC()
	if _, err = conn.Exec(insertIntoFilesTableSQL, code, name, fileSize, &now); err != nil {
		panic(err.Error())
	}
	return
}

func SelectFileInfo(code string) *FileInfo {
	var fileInfo FileInfo
	row := conn.QueryRow(selectFileInfoSql, code)
	if err := row.Scan(&fileInfo.Code, &fileInfo.Name, &fileInfo.FileSize, &fileInfo.UploadedOn); err != nil {
		panic(err.Error())
	}
	return &fileInfo
}

func FileInfoExists(code string) bool {
	var count int
	conn.QueryRow(countFileWithCode, code).Scan(&count)
	return count >= 1
}
