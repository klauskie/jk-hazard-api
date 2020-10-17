package controller

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"klaus.com/jkapi/app/entity"
	"klaus.com/jkapi/app/repository"
	"klaus.com/jkapi/app/service"
	"net/http"
	"strconv"
)


var (
	//roomRepo repository.RoomRepository = repository.NewRoomRepository()
)

// GET: rooms
// temp
func GetAllRooms(w http.ResponseWriter, r *http.Request) {
	rooms := repository.GetRoomRepository().GetAll()
	FillHeaders(w, http.StatusOK)
	json.NewEncoder(w).Encode(rooms)
}

// GET: room/{id}
func GetRoom(w http.ResponseWriter, r *http.Request) {
	// Auth
	token := r.Header.Get("token")
	_, err := GetUserByToken(token)
	if err != nil {
		FillErrorHeaders(w, http.StatusUnauthorized, "Error token not identified")
		return
	}

	params := mux.Vars(r)
	tag := params["roomTAG"]

	room := repository.GetRoomRepository().Get(tag)
	if room == nil {
		FillErrorHeaders(w, http.StatusNotFound, "Error: cannot find room with given tag.")
		return
	}
	FillHeaders(w, http.StatusOK)
	json.NewEncoder(w).Encode(room)
}

// GET: room/{id}/judge-cards
func GetJudgeCards(w http.ResponseWriter, r *http.Request) {
	// Auth
	token := r.Header.Get("token")
	_, err := GetUserByToken(token)
	if err != nil {
		FillErrorHeaders(w, http.StatusUnauthorized, "Error token not identified")
		return
	}

	params := mux.Vars(r)
	tag := params["roomTAG"]

	room := repository.GetRoomRepository().Get(tag)
	if room == nil {
		FillErrorHeaders(w, http.StatusNotFound, "Error: cannot find room with given tag.")
		return
	}

	result := make(map[string]interface{})
	result["cardsOnTrial"] = room.CardsOnTrial

	FillHeaders(w, http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// POST: room
func NewRoom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("token")
	key := service.GetSessionHandlerInstance().GetUserByKey(token)
	if key == -1 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "Error token not identified"}`))
	}
	user := userRepo.FindById(key)
	player := entity.NewPlayer(user, true)
	room := entity.NewRoom(player)
	room.UpdateLastConnection(user)
	repository.GetRoomRepository().Add(room)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(room)
}


// PUT: room/join
func JoinRoom(w http.ResponseWriter, r *http.Request) {
	// Auth
	token := r.Header.Get("token")
	user, err := GetUserByToken(token)
	if err != nil {
		FillErrorHeaders(w, http.StatusUnauthorized, "Error token not identified")
		return
	}

	// Get roomTAG from body
	result := make(map[string]string)
	err1 := json.NewDecoder(r.Body).Decode(&result)
	if err1 != nil {
		FillErrorHeaders(w, http.StatusInternalServerError, "Error unmarshalling error")
		return
	}
	tag := result["roomTag"]
	player := entity.NewPlayer(user, false)
	err2 := repository.GetRoomRepository().Get(tag).JoinPlayer(player)
	if err2 != nil {
		FillErrorHeaders(w, http.StatusInternalServerError, err2.Error())
	}

	// Set connection time
	repository.GetRoomRepository().Get(tag).UpdateLastConnection(user)

	FillHeaders(w, http.StatusOK)
	json.NewEncoder(w).Encode(repository.GetRoomRepository().Get(tag))
}

// PUT: room/{roomTAG}/send-card
func UpdateCardsOnTrial(w http.ResponseWriter, r *http.Request) {
	// Auth
	token := r.Header.Get("token")
	user, err := GetUserByToken(token)
	if err != nil {
		FillErrorHeaders(w, http.StatusUnauthorized, "Error token not identified")
		return
	}
	// Get Url params
	params := mux.Vars(r)
	tag := params["roomTAG"]

	// Get body params
	bodyParams := make(map[string]string)
	err1 := json.NewDecoder(r.Body).Decode(&bodyParams)
	if err1 != nil {
		FillErrorHeaders(w, http.StatusInternalServerError, "Error unmarshalling error")
		return
	}

	// Get room by tag
	room := repository.GetRoomRepository().Get(tag)
	if room == nil {
		FillErrorHeaders(w, http.StatusNotFound, "Error: cannot find room with given tag.")
		return
	}

	// Get Player
	player := room.GetPlayerByUsername(user.Username)
	if player.IsEmpty() {
		FillErrorHeaders(w, http.StatusNotFound, "Error: cannot find player with given username.")
		return
	}

	// get card from player's hand
	cardID, _ := strconv.ParseInt(bodyParams["cardID"], 10, 64)
	card := player.GetCardById(cardID)
	if !card.IsEmpty() {
		// remove card from player's hand
		player.RemoveCardById(cardID)
	}

	// Add card to Trial
	room.AddToTrial(card, player)

	// New card and add it to player's hand
	newCard := room.PopDeckStack()
	player.AddCardToHand(newCard)

	// Save Room
	room.SavePlayer(player)
	//repository.GetRoomRepository().Save(room)


	// Result Map
	results := make(map[string]interface{})
	results["newCard"] = newCard


	FillHeaders(w, http.StatusOK)
	json.NewEncoder(w).Encode(results)
}

// PUT: /api/room/{roomTAG}/round-winner
func SetRoundWinner(w http.ResponseWriter, r *http.Request) {
	// TODO: Update room with provided player from the judge
	// Auth
	token := r.Header.Get("token")
	user, err := GetUserByToken(token)
	if err != nil {
		FillErrorHeaders(w, http.StatusUnauthorized, "Error token not identified")
		return
	}

	params := mux.Vars(r)
	tag := params["roomTAG"]
	// Get Room
	room := repository.GetRoomRepository().Get(tag)
	if room == nil {
		FillErrorHeaders(w, http.StatusNotFound, "Error: cannot find room with given tag.")
		return
	}

	// Get Player
	judge := room.GetPlayerByUsername(user.Username)
	if judge.IsEmpty() {
		FillErrorHeaders(w, http.StatusNotFound, "Error: cannot find player with given username.")
		return
	}

	// Is Caller a Judge ?
	if !room.IsPlayerJudge(judge) {
		FillErrorHeaders(w, http.StatusForbidden, "Error: player is not a judge.")
		return
	}

	// Get body params
	bodyParams := make(map[string]string)
	err1 := json.NewDecoder(r.Body).Decode(&bodyParams)
	if err1 != nil {
		FillErrorHeaders(w, http.StatusInternalServerError, "Error unmarshalling error")
		return
	}

	// Get round winner
	player := room.GetPlayerByUsername(bodyParams["username"])
	if player.IsEmpty() {
		FillErrorHeaders(w, http.StatusNotFound, "Error: cannot find player with given username.")
		return
	}

	// Update player's score
	player.AddToScore(1)

	// Save Player
	room.SavePlayer(player)

	// Run next round preparations
	room.NextRound()

	// Result Map
	results := make(map[string]interface{})
	results["message"] = "Round winner received and updated"


	FillHeaders(w, http.StatusOK)
	json.NewEncoder(w).Encode(results)
}

// GET: /api/room/{roomTAG}/heart-beat
func HeartBeat(w http.ResponseWriter, r *http.Request) {
	// Auth
	token := r.Header.Get("token")
	user, err := GetUserByToken(token)
	if err != nil {
		FillErrorHeaders(w, http.StatusUnauthorized, "Error token not identified")
		return
	}

	params := mux.Vars(r)
	tag := params["roomTAG"]

	room := repository.GetRoomRepository().Get(tag)
	if room == nil {
		FillErrorHeaders(w, http.StatusNotFound, "Error: cannot find room with given tag.")
		return
	}

	// Update room with last connection difference
	room.UpdateLastConnection(user)


	FillHeaders(w, http.StatusOK)
	json.NewEncoder(w).Encode(room)
}

// GET: /api/room/start/{roomTAG}
func StartGame(w http.ResponseWriter, r *http.Request) {
	// Auth
	token := r.Header.Get("token")
	user, err := GetUserByToken(token)
	if err != nil {
		FillErrorHeaders(w, http.StatusUnauthorized, "Error token not identified")
		return
	}

	params := mux.Vars(r)
	tag := params["roomTAG"]

	room := repository.GetRoomRepository().Get(tag)
	if room == nil {
		FillErrorHeaders(w, http.StatusNotFound, "Error: cannot find room with given tag.")
		return
	}

	// Update room with last connection difference
	room.UpdateLastConnection(user)
	room.InitPlayersHand()
	room.Judge = room.Host


	FillHeaders(w, http.StatusOK)
	json.NewEncoder(w).Encode(room)
}

func FillHeaders(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
}

func FillErrorHeaders(w http.ResponseWriter, status int, errorMessage string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	message := make(map[string]string)
	message["error"] = errorMessage
	json.NewEncoder(w).Encode(message)
}