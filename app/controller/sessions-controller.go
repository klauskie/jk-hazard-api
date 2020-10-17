package controller

import (
	"klaus.com/jkapi/app/entity"
	"klaus.com/jkapi/app/service"
	"log"
	"encoding/json"
	"net/http"
)

type SimpleUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// POST: login
func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var tempUser SimpleUser
	err := json.NewDecoder(r.Body).Decode(&tempUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error unmarshalling error"}`))
		return
	}

	// Get User by Username
	user, err1 := userRepo.FindByUsername(tempUser.Username)
	if err1 != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Error username doesn't exist"}`))
		log.Println(err1)
		return
	}

	// Vaidate
	if user.Password != tempUser.Password {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Error incorrect password"}`))
	}

	token := service.GetSessionHandlerInstance().AddNewSession(user)
	
	result := make(map[string]string)
	result["token"] = token

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// POST: signup
func SignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var tempUser SimpleUser
	err := json.NewDecoder(r.Body).Decode(&tempUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error unmarshalling error"}`))
		return
	}
	user := entity.User{Username:tempUser.Username, Password:tempUser.Password}

	// Create User
	err1 := userRepo.Create(&user)
	if err1 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error creating New User"}`))
		log.Println(err1)
		return
	}

	token := service.GetSessionHandlerInstance().AddNewSession(user)
	
	result := make(map[string]interface{})
	result["token"] = token
	result["user"] = user

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// GET: logout
func Logout(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")
	service.GetSessionHandlerInstance().RemoveSession(token)
	w.WriteHeader(http.StatusOK)
}