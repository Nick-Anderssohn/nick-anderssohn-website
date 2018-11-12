package auth

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
)
import "golang.org/x/crypto/bcrypt"

type userEntity struct {
	Username      string
	PasswordHash  string
	CanViewKibana bool
}

const (
	// User table holds both Authn and Authz info for a user.
	createUserTableSql = `
CREATE TABLE IF NOT EXISTS Users(
username TEXT PRIMARY KEY,
passwordHash TEXT NOT NULL,
canViewKibana BOOLEAN NOT NULL DEFAULT FALSE
)
`
	insertUserSql = `INSERT INTO Users VALUES($1, $2, $3)`
	selectUserSql = `SELECT * FROM Users WHERE username = $1`
)

const adminPath = "admin.json"

func insertAdmin() {
	adminBytes, _ := ioutil.ReadFile(adminPath)
	var admin UserReq
	json.Unmarshal(adminBytes, &admin)
	if admin.Username == "" {
		panic("no admin user provided. create an admin.json file.")
	}
	adminEntity, _ := userReqToEntity(&admin)
	adminEntity.CanViewKibana = true
	_, err := selectUser(adminEntity.Username)
	if err == sql.ErrNoRows {
		err = insertUser(adminEntity)
		if err != nil {
			panic("could not insert admin user.")
		}
	}
}

func insertUser(user *userEntity) error {
	_, err := db.Exec(insertUserSql, user.Username, user.PasswordHash, user.CanViewKibana)
	return err
}

func selectUser(username string) (*userEntity, error) {
	var user userEntity
	row := db.QueryRow(selectUserSql, username)
	err := row.Scan(&user.Username, &user.PasswordHash, &user.CanViewKibana)
	return &user, err
}

func hashPassword(password string) (string, error) {
	log.Println("WHAT")
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	log.Println("THE FUCK")
	return string(hashBytes), err
}

func userReqToEntity(userReq *UserReq) (*userEntity, error) {
	passHash, err := hashPassword(userReq.Password)
	if err != nil {
		return nil, err
	}

	return &userEntity{
		Username:     userReq.Username,
		PasswordHash: passHash,
	}, nil
}