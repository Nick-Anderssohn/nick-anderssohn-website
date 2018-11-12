package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

type dbConfig struct {
	Username         string `json:"Username"`
	Password         string `json:"Password"`
	PostgresPassword string `json:"PostgresPassword"`
	DbName           string `json:"DbName"`
	Host             string `json:"Host"`
}

const dbConfigPath = "db_config.json"

var config dbConfig
var db *sql.DB

func init() {
	log.Println("read db config")
	readDbConfig()
	log.Println("create db if not exist")
	createDbIfNotExist()
	log.Println("create conn pool")
	createConnPool()
	log.Println("create upgrade table")
	createUpgradeTable()
	log.Println("run upgrades")
	runUpgrades()
	log.Println("insert admin")
	insertAdmin()
	log.Println("done with db init")
}

func readDbConfig() {
	configBytes, _ := ioutil.ReadFile(dbConfigPath)
	json.Unmarshal(configBytes, &config)
	if config.DbName == "" {
		panic("no database name")
	}
	config.DbName = strings.ToLower(config.DbName)
}

func getConnStr(host, user, password, dbname string) string {
	return fmt.Sprintf("host=%s user='%s' password='%s' dbname='%s' sslmode=disable", host, user, password, dbname)
}


func createDbIfNotExist() {
	var count int
	connStr := getConnStr(config.Host, "postgres", config.PostgresPassword, "postgres")
	postgresConn, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	rows := postgresConn.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM pg_database WHERE datname = '%s'", config.DbName))
	rows.Scan(&count)
	if count == 0 {
		fmt.Println("Creating", config.DbName, "database")
		_, err = postgresConn.Exec(fmt.Sprintf("CREATE DATABASE %s OWNER postgres", config.DbName))
		if err != nil {
			panic(err)
		}
	}
}

func createConnPool() {
	var err error
	connStr := getConnStr(config.Host, config.Username, config.Password, config.DbName)
	if db, err = sql.Open("postgres", connStr); err != nil {
		panic(err)
	}
}