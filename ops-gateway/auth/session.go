package auth

import (
	"github.com/satori/go.uuid"
	"net/http"
	"time"
)

type Session struct {
	SessionId string
	Expires *time.Time
	Username string
	user *userEntity // Not a field in db. Select based off of Username.
}

const (
	createSessionTableSql = `
CREATE TABLE IF NOT EXISTS session(
  sessionId TEXT NOT NULL PRIMARY KEY,
  expires TIMESTAMP NOT NULL,
  username TEXT REFERENCES Users(Username) NOT NULL ON DELETE CASCADE,
  CHECK (sessionId != '')
)`

	insertSessionSql = `INSERT INTO session VALUES ($1, $2, $3)`
	selectSessionSql   = `SELECT * FROM session WHERE sessionId = $1`
)

const sessionDuration = time.Hour * 24 * 365
const sessionCookieName = "session"

func NewSession(username string) *Session {
	sessionId := uuid.NewV4().String()
	expires := time.Now().Add(sessionDuration)
	return &Session{
		SessionId: sessionId,
		Expires: &expires,
		Username: username,
	}
}

func (s *Session) getUser() (*userEntity, error) {
	if s.user != nil {
		return s.user, nil
	}

	user, err := selectUser(s.Username)
	if err != nil {
		return nil, err
	}
	s.user = user
	return s.user, err
}

func insertSession(session *Session) error {
	_, err := db.Exec(insertSessionSql, session.SessionId, session.Expires, session.Username)
	return err
}

// Queries based off of sessionId
func selectSession(sessionId string) (*Session, error) {
	var session Session
	row := db.QueryRow(selectSessionSql, sessionId)
	err := row.Scan(&session.SessionId, &session.Expires, &session.Username)

	return &session, err
}

func (s *Session) toCookie() *http.Cookie {
	return &http.Cookie{
		Name: sessionCookieName,
		Value: s.SessionId,
		Expires: *s.Expires,
		Secure: true,
		HttpOnly: true,
		//SameSite: http.SameSiteStrictMode,
	}
}