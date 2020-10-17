package controller

import (
	"encoding/json"
	"errors"
	"klaus.com/jkapi/app/repository"
	"klaus.com/jkapi/app/service"
	"log"
	"net/http"
)


var (
	userRepo repository.UserRepository = repository.NewUserRepository()
)

type responseCallback func(w http.ResponseWriter)


// GET: users
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users, err := userRepo.FindAll()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Error can't retreive users"}`))
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

/*
func GetUserByToken(token string) (entity.User, error) {
	var user entity.User
	key := service.GetSessionHandlerInstance().GetPlayerIDByKey(token)
	if len(key) == 0 {
		return user, errors.New("Error token not identified")
	}
	user = userRepo.FindById(key)
	return user, nil
}
 */

func GetPlayerKeyByToken(token string) (string, error) {
	key := service.GetSessionHandlerInstance().GetPlayerIDByKey(token)
	if len(key) == 0 {
		return key, errors.New("Error token not identified")
	}
	return key, nil
}