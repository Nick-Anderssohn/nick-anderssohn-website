package auth

import (
	"database/sql"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"nick-anderssohn-website/ops-gateway/serverutil"
)

type UserReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var LoginEndpoint = &serverutil.Endpoint {
	Path: "/SubmitLogin",
	HandleFunc: handleLogin,
}

func recoverAndLog() {
	if e := recover(); e != nil {
		log.Println(e)
	}
}

// This endpoint ignores any existing cookies the caller may have.
// Will always go through the login process.
func handleLogin(respWriter http.ResponseWriter, request *http.Request) {
	defer recoverAndLog()

	// Defer sending result
	statusCode := http.StatusInternalServerError
	var msg string
	defer func() {
		respWriter.WriteHeader(statusCode)
		respWriter.Write([]byte(msg))
	}()

	// Get login info from request
	var loginReq UserReq
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&loginReq)
	if err != nil {
		log.Println(err)
		statusCode = http.StatusBadRequest
		msg = "Bad request."
		return // Bad Request
	}

	log.Println("user req", loginReq)

	// Check if credentials pass
	userEntity, err := selectUser(loginReq.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("user not found")
			statusCode = http.StatusNotFound
			msg = "User not found."
			return // User not found
		}
		log.Println("unknown error ", err)
		msg = "Unknown error."
		return // Unknown error
	}

	err = bcrypt.CompareHashAndPassword([]byte(userEntity.PasswordHash), []byte(loginReq.Password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			log.Println("incorrect password provided for user: ", loginReq.Username)
			statusCode = http.StatusUnauthorized
			msg = "Incorrect password."
			return // Unauthorized
		}
		log.Println("unknown error ", err)
		msg = "Unknown error."
		return
	}

	// Password matched. Create session.
	session := NewSession(userEntity.Username)
	err = insertSession(session)
	if err != nil {
		log.Println("unknown error ", err)
		msg = "Unknown error."
		return
	}
	cacheSession(session)

	// Reply with cookie and successful response
	log.Println("cookie:", session.toCookie())
	http.SetCookie(respWriter, session.toCookie())
	statusCode = http.StatusOK
	msg = "Login successful."
}