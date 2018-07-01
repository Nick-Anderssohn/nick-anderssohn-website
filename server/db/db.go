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
	PostgresPassword string `json:"PostgresPassword"`
	DbName string `json:"DbName"`
	Host string `json:"Host"`
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
VALUES ($1, $2, $3, $4);
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
	time.Sleep(5 * time.Second) // give postgres a little bit of time to start.
	createDbIfNotExist()
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
	var err error
	connStr := getConnStr(config.Host, config.Username, config.Password, config.DbName)
	if conn, err = sql.Open("postgres", connStr); err != nil {
		panic(err.Error())
	}
}

func createDbIfNotExist() {
	for tries := 0; tries < 10; tries++ {
		var count int
		connStr := getConnStr(config.Host, "postgres", config.PostgresPassword, "postgres")
		postgresConn, err := sql.Open("postgres", connStr)
		if err == nil {
			tries = 10
			rows := postgresConn.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM pg_database WHERE datname = '%s'", config.DbName))
			rows.Scan(&count)
			if count == 0 {
				fmt.Println("Creating", config.DbName, "database")
				postgresConn.Exec(fmt.Sprintf("CREATE DATABASE %s OWNER postgres", config.DbName))
			}
		} else {
			time.Sleep(1 * time.Second)
		}
	}

}

func getConnStr(host, user, password, dbname string) string {
	return fmt.Sprintf("host=%s user='%s' password='%s' dbname='%s' sslmode=disable", host, user, password, dbname)
}

func CreateTablesIfNotExists() {
	if _, err := conn.Exec(createFilesTableIFNotExistsSql); err != nil {
		panic(err.Error())
	}
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
