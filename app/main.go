package main

import (
	"fmt"
	"klaus.com/jkapi/app/controller"
	"klaus.com/jkapi/app/database"
	"klaus.com/jkapi/app/service"
	"log"
	"net/http"



	"github.com/gorilla/mux"
)

//docker build -t jk-hazard-img .
//docker run -it --rm -p 127.0.0.1:8080:8080/tcp --name jk-api jk-hazard-img


func init() {
	service.GetSessionHandlerInstance()
}

func main() {
	defer database.DB.Close()
	initRoutes()
}

func initRoutes() {
	router := mux.NewRouter()
/*
	router.HandleFunc("/api/login", controller.Login).Methods("POST")
	router.HandleFunc("/api/logout", controller.Logout).Methods("GET")
	router.HandleFunc("/api/signup", controller.SignUp).Methods("POST")
 */
	router.HandleFunc("/api/users", controller.GetAllUsers).Methods("GET")

	router.HandleFunc("/api/rooms", controller.GetAllRooms).Methods("GET")
	router.HandleFunc("/api/room", controller.NewRoom).Methods("POST")
	router.HandleFunc("/api/room/join", controller.JoinRoom).Methods("PUT")
	router.HandleFunc("/api/room/{roomTAG}", controller.GetRoom).Methods("GET")
	router.HandleFunc("/api/room/{roomTAG}/start", controller.StartGame).Methods("GET")
	router.HandleFunc("/api/room/{roomTAG}/send-card", controller.UpdateCardsOnTrial).Methods("PUT")
	router.HandleFunc("/api/room/{roomTAG}/judge-cards", controller.GetJudgeCards).Methods("GET")
	router.HandleFunc("/api/room/{roomTAG}/round-winner", controller.SetRoundWinner).Methods("PUT")
	router.HandleFunc("/api/room/{roomTAG}/heart-beat", controller.HeartBeat).Methods("GET")
	router.HandleFunc("/api/room/{roomTAG}/host", controller.GetHost).Methods("GET")
	router.HandleFunc("/api/room/{roomTAG}/judge", controller.GetCurrentJudge).Methods("GET")
	router.HandleFunc("/api/room/{roomTAG}/player", controller.GetPlayer).Methods("GET")

	http.Handle("/", router)
	fmt.Println("Starting up on 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}