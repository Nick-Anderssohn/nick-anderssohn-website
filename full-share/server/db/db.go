package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	cfg "nick-anderssohn-website/full-share/server/config"
	"nick-anderssohn-website/full-share/server/slog"
	"strings"
	"time" // Blank import required by PostgreSQL driver.

	"github.com/Nick-Anderssohn/sherlog"
	_ "github.com/lib/pq"
)

type dbConfig struct {
	Username         string `json:"Username"`
	Password         string `json:"Password"`
	PostgresPassword string `json:"PostgresPassword"`
	DbName           string `json:"DbName"`
	Host             string `json:"Host"`
}

var config dbConfig

const dbConfigPath = "db_config.json"

// FileInfo contains the information associated with a file entry in the files database.
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
  uploadedOn TIMESTAMP 
);
`

	createUpgradesTableIfNotExists = `
CREATE TABLE IF NOT EXISTS upgrades(
  upgradeNumber INTEGER PRIMARY KEY
);
`

	insertIntoFilesTableSQL = `
INSERT INTO files
VALUES ($1, $2, $3, $4);
`

	selectFileInfoSql = `SELECT * FROM files WHERE code = $1;`

	deleteOlderThan = `DELETE FROM files WHERE date_part('day', age($1, uploadedOn)) >= $2 RETURNING code;`

	deleteFromCode = `DELETE FROM files WHERE code = $1;`
)

var conn *sql.DB

func init() {
	if !cfg.ApplicationConfig.TestMode {
		readDbConfig()
		createDbIfNotExist()
		connectToDb()
		createTablesIfNotExists()
		runUpgrades()
	}
}

func readDbConfig() {
	configBytes, _ := ioutil.ReadFile(dbConfigPath)
	json.Unmarshal(configBytes, &config)
	if config.DbName == "" {
		panic("no database name")
	}
	config.DbName = strings.ToLower(config.DbName)
}

func connectToDb() {
	var err error
	connStr := getConnStr(config.Host, config.Username, config.Password, config.DbName)
	if conn, err = sql.Open("postgres", connStr); err != nil {
		slog.Logger.OpsError("could not connect to database: ", err)
		panic(err)
	}
}

func createDbIfNotExist() {
	var count int
	connStr := getConnStr(config.Host, "postgres", config.PostgresPassword, "postgres")
	postgresConn, err := sql.Open("postgres", connStr)
	if err != nil {
		slog.Logger.OpsError("could not connect to database: ", err)
		panic(err)
	}
	rows := postgresConn.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM pg_database WHERE datname = '%s'", config.DbName))
	rows.Scan(&count)
	if count == 0 {
		slog.Logger.Info("Creating ", config.DbName, " database")
		_, err = postgresConn.Exec(fmt.Sprintf("CREATE DATABASE %s OWNER postgres", config.DbName))
		if err != nil {
			panic(err)
		}
	}
}

func getConnStr(host, user, password, dbname string) string {
	return fmt.Sprintf("host=%s user='%s' password='%s' dbname='%s' sslmode=disable", host, user, password, dbname)
}

func createTablesIfNotExists() {
	if _, err := conn.Exec(createFilesTableIFNotExistsSql); err != nil {
		slog.Logger.OpsError(err)
		panic(err)
	}
	if _, err := conn.Exec(createUpgradesTableIfNotExists); err != nil {
		slog.Logger.OpsError(err)
		panic(err)
	}
}

// InsertNewFileInfo inserts a new file info into the files table
func InsertNewFileInfo(code, name string, fileSize int) (err error) {
	now := time.Now().UTC()
	_, err = conn.Exec(insertIntoFilesTableSQL, code, name, fileSize, &now)
	return
}

// SelectFileInfo returns a file info based off of code
func SelectFileInfo(code string) (fileInfo *FileInfo, err error, errNoRows bool) {
	fileInfo = &FileInfo{}
	row := conn.QueryRow(selectFileInfoSql, code)
	err = row.Scan(&fileInfo.Code, &fileInfo.Name, &fileInfo.FileSize, &fileInfo.UploadedOn)
	return fileInfo, sherlog.AsOpsError(err), err == sql.ErrNoRows
}

// FileInfoExists checks if a file exists. Returns an error if there is a problem querying the database.
func FileInfoExists(code string) (bool, error) {
	_, err, doesNotExist := SelectFileInfo(code)
	if doesNotExist {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// DeleteFilesOlderThan deletes all file entries from the files table that are older than days.
func DeleteFilesOlderThan(days int) (codes []string, err error) {
	now := time.Now().UTC()
	rows, err := conn.Query(deleteOlderThan, &now, days)
	if err != nil {
		return codes, sherlog.AsOpsError(err)
	}
	for rows.Next() {
		var code string
		err = rows.Scan(&code)
		if err != nil {
			return codes, sherlog.AsError(err)
		}
		codes = append(codes, code)
	}
	return
}

// DeleteDbEntryFromCode deletes the entry in the database with the matching code.
func DeleteDbEntryFromCode(code string) error {
	_, err := conn.Exec(deleteFromCode, code)
	return sherlog.AsOpsError(err)
}
