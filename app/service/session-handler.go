package service

import(
	"klaus.com/jkapi/app/entity"
	"strings"

	"github.com/google/uuid"
)

var instance *sessionHandler = nil

type sessionHandler struct {
	sessions map[string]string
	loggedUsers map[string]bool
}

// Singleton Constructor
func GetSessionHandlerInstance() *sessionHandler {
	if instance == nil {
		instance = new(sessionHandler)
		instance.sessions = make(map[string]string)
		instance.loggedUsers = make(map[string]bool)
	}
	return instance
}


func (handler* sessionHandler) AddNewSession(player entity.Player) string {
	if _, ok := handler.loggedUsers[player.ID]; ok {
		return ""
	}
	key := newKey()
	handler.sessions[key] = player.ID
	handler.loggedUsers[player.ID] = true
	return key
}

func (handler* sessionHandler) GetPlayerIDByKey(key string) string {
	if _, ok := handler.sessions[key]; ok {
		return handler.sessions[key]
	}
	return ""
}

func (handler* sessionHandler) RemoveSession(key string) {
	if _, ok := handler.loggedUsers[handler.sessions[key]]; ok {
		delete(handler.loggedUsers, handler.sessions[key])
		delete(handler.sessions, key)
	}
}

func newKey() string {
	key := uuid.New().String()
	key = strings.ReplaceAll(key, "-", "")
	return key
}

