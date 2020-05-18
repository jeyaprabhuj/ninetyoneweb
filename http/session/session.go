package session

type session struct {
	externalStorage SessionExternalStore
}

type UserSession struct {
	ID         string
	refSession *session
}

type SessionExternalStore interface {
	Store(string, string, interface{}) error
	Find(string, string) (interface{}, error)
	IsSessionAvailable(string) bool
}
