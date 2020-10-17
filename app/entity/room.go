package entity

import (
	"bytes"
	"encoding/csv"
	"errors"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Room struct {
	TAG				string
	Players			PlayerList
	Deck			DeckStack
	TableDeck		map[int]Card
	Host			Player
	Judge			Player
	CardsOnTrial	map[string]Card
	RollPlayers 	map[string]int64
}

func NewRoom(player Player) *Room {
	room := Room{}
	tag := createTag()
	player.ID = player.Username + "-" + tag

	room.TAG = tag
	room.Players = PlayerList{player}
	room.Deck = NewDeckStack()
	room.TableDeck = make(map[int]Card)
	room.Host = player
	room.Judge = Player{}
	room.CardsOnTrial = make(map[string]Card)
	room.RollPlayers = make(map[string]int64)

	return &room
}

// Join Player to Room
func (room *Room) JoinPlayer(player Player) error {
	if player.Username == room.Host.Username {
		return errors.New("error: host cannot join the guests")
	}
	if room.Players == nil {
		return errors.New("error: list is null")
	}
	if room.Players.Contains(player) {
		return errors.New("error: username provided already exists")
	}
	player.ID = player.Username + "-" + room.TAG
	room.Players = append(room.Players, player)
	return nil
}

func (room *Room) GetPlayerByUsername(username string) Player {
	player, _ := room.Players.GetByUsername(username)
	return player
}

func (room *Room) GetPlayerIndex(player Player) int {
	index := -1
	for i, p := range room.Players {
		if p.Equals(player) {
			index = i
			break
		}
	}
	return index
}

func (room *Room) AddToTrial(card Card, player Player) {
	if !room.Judge.Equals(player) {
		room.CardsOnTrial[player.Username] = card
	}
}

func (room *Room) PopDeckStack() Card {
	return room.Deck.pop()
}

func (room *Room) BatchPopDeckStack(batchSize int) []Card {
	cards := []Card{}
	for i := 0; i < batchSize; i++ {
		cards = append(cards, room.Deck.pop())
	}
	return cards
}

func (room *Room) SavePlayer(player Player) {
	index := room.GetPlayerIndex(player)
	room.Players[index] = player
}

func (room *Room) IsPlayerJudge(player Player) bool {
	return room.Judge.Equals(player)
}

func (room *Room) NextRound() {
	//Choose new judge
	judgeIndex := room.GetPlayerIndex(room.Judge)
	judgeIndex++
	room.Judge = room.Players[judgeIndex % room.PlayerListSize()]

	// Clear table deck
	room.TableDeck = make(map[int]Card)

	// Clear cards on trial
	room.CardsOnTrial = make(map[string]Card)
}

func (room *Room) PlayerListSize () int {
	return len(room.Players)
}

func (room *Room) UpdateLastConnection (user Player) {
	timeNow := time.Now().UnixNano() / 1000000
	lastBeat := room.RollPlayers[user.Username]
	room.RollPlayers[user.Username] = timeNow - lastBeat
}

func (room *Room) InitPlayersHand () {
	for _, p := range room.Players {
		p.Hand = room.BatchPopDeckStack(5)
		room.SavePlayer(p)
	}
}

func createTag() string {
	buff := bytes.NewBufferString("")
	for i := 0; i < 4; i++ {
		letter := rand.Intn(90 - 65) + 65
		buff.WriteString(string(letter))
	}
	return buff.String()
}

func shuffle(list []Card) []Card {
	rand.Seed(time.Now().UnixNano())

	for i := range list {
		j := rand.Intn(i + 1)
		list[i], list[j] = list[j], list[i]
	}
	return list
}

/* ------------------ Deck Stack ------------------ */
type DeckStack struct {
	collection []Card
}

func NewDeckStack() DeckStack {
	cards, err := readCardCSV()
	if err != nil {
		log.Fatal(err)
	}
	shuffledCards := shuffle(cards)
	return DeckStack{shuffledCards}
}

func (s *DeckStack) push(element Card) {
	s.collection = append(s.collection, element)
}

func (s *DeckStack) pop() Card {
	card := s.collection[len(s.collection)-1]
	s.collection = s.collection[:len(s.collection)-1]
	return card
}

func (s *DeckStack) top() Card {
	return s.collection[len(s.collection)-1]
}

func (s *DeckStack) isEmpty() bool {
	return len(s.collection) == 0
}
/* ------------------ END DeckStack ------------------ */

/* ------------------ START CSV READER ------------------ */

func readCardCSV() ([]Card, error) {
	f, err := os.Open("cards.csv")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	csvr := csv.NewReader(f)

	cards := []Card{}
	for {
		row, err := csvr.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return cards, err
		}

		cardName := row[0]
		cardID, _ := strconv.Atoi(cardName[strings.Index(row[0], "-")+1:])

		c := Card{int64(cardID), cardName, row[1] == "True"}

		cards = append(cards, c)
	}
}

/* ------------------ END CSV READER ------------------ */

/* ------------------ START TYPE OVERLOADING ------------------ */

type PlayerList []Player

type IPlayerList interface {
	IsEmpty() bool
	Contains(needle Player) bool
	GetByUsername(username string) (Player, error)
	GetByKey(key string) (Player, error)
}

func (l PlayerList) IsEmpty() bool {
	return len(l) == 0
}

func (l PlayerList) Contains(needle Player) bool {
	for _, p := range l {
		if p.Equals(needle) {
			return true
		}
	}
	return false
}

func (l PlayerList) GetByUsername(username string) (Player, error) {
	for _, p := range l {
		if p.Username == username {
			return p, nil
		}
	}
	return Player{}, errors.New("error: username provided not found in list")
}

// KEY stored in session-handler
func (l PlayerList) GetByKey(key string) (Player, error) {
	username := key[0:strings.Index(key, "-")]

	for _, p := range l {
		if p.Username == username {
			return p, nil
		}
	}
	return Player{}, errors.New("error: username provided not found in list")
}

/* ------------------ END TYPE OVERLOADING ------------------ */