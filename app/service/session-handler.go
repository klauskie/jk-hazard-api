package service

import(
	"klaus.com/jkapi/app/entity"
	"strings"

	"github.com/google/uuid"
)

var instance *sessionHandler = nil

type sessionHandler struct {
	sessions map[string]int64
	loggedUsers map[int64]bool
}

// Singleton Constructor
func GetSessionHandlerInstance() *sessionHandler {
	if instance == nil {
		instance = new(sessionHandler)
		instance.sessions = make(map[string]int64)
		instance.loggedUsers = make(map[int64]bool)
	}
	return instance
}


func (handler* sessionHandler) AddNewSession(user entity.User) string {
	if _, ok := handler.loggedUsers[user.ID]; ok {
		return ""
	}
	key := newKey()
	handler.sessions[key] = user.ID
	handler.loggedUsers[user.ID] = true
	return key
}

func (handler* sessionHandler) GetUserByKey(key string) int64 {
	if _, ok := handler.sessions[key]; ok {
		return handler.sessions[key]
	}
	return -1;
}

func (handler* sessionHandler) RemoveSession(key string) {
	if _, ok := handler.loggedUsers[handler.sessions[key]]; ok {
		delete(handler.loggedUsers, handler.sessions[key])
		delete(handler.sessions, key);
	}
}

func newKey() string {
	key := uuid.New().String()
	key = strings.ReplaceAll(key, "-", "")
	return key
}

