package session

import (
	"github.com/google/uuid"
)

var efsSession *session

//NewSessionManager creates a new session map containing userSession
//return existing session if any
func NewSessionManager(storage SessionExternalStore) *session {
	if efsSession == nil {
		efsSession = &session{externalStorage: storage}
	}
	return efsSession
}

func (s *session) NewUserSession() *UserSession {
	return &UserSession{
		ID:         uuid.New().String(),
		refSession: s,
	}
}

func (s *session) GetUserSession(id string) (*UserSession, bool) {
	if s.externalStorage != nil {
		isAvailable := s.externalStorage.IsSessionAvailable(id)
		return &UserSession{ID: id, refSession: s}, isAvailable
	}
	return nil, false
}

func (s *session) DeleteSession(id string) error {
	return nil
}

func (us *UserSession) Store(key string, value interface{}) error {
	if us.refSession.externalStorage != nil {
		return us.refSession.externalStorage.Store(us.ID, key, value)
	}
	return nil
}

func (us *UserSession) Find(key string) interface{} {
	if us.refSession.externalStorage != nil {
		tempUs, _ := us.refSession.externalStorage.Find(us.ID, key)
		return tempUs
	}
	return nil

}
