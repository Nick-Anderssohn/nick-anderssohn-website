package auth

import (
	"sync"
	"time"
)

var cachedSessions = sync.Map{}

func cacheSession(session *Session) {
	cachedSessions.Store(session.SessionId, session)
	go deleteAfterOneDay(session.SessionId)
}

func deleteAfterOneDay(key string) {
	time.Sleep(time.Hour * 24)
	cachedSessions.Delete(key)
}

// Will grab cached if exists, otherwise queries and performs cache before returning.
func grabSessionFromId(sessionId string) (*Session, error) {
	session, found := cachedSessions.Load(sessionId)
	if found {
		return session.(*Session), nil
	}

	session, err := selectSession(sessionId)
	if err != nil {
		return nil, err
	}

	cacheSession(session.(*Session))

	return session.(*Session), nil
}