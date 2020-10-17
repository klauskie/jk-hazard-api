package entity

import (
	"strconv"
	"time"
)

type Player struct {
	ID string
	Username string
	Score int
	IsHost bool
	Hand []Card
}

func NewPlayer(username string, isHost bool) Player {
	id := username + "-" + strconv.FormatInt(time.Now().UnixNano() / 1000000, 10) // This is changed when the room is created or a player joins in [username-roomTAG]
	player := Player{id, username, 0, isHost, []Card{}}
	return player
}

func (player *Player) GetCardById(id int64) Card {
	var result Card
	for _, card := range player.Hand {
		if card.ID == id {
			result = card
			break
		}
	}
	return result
}

func (player *Player) AddCardToHand(card Card) {
	player.Hand = append(player.Hand, card)
}

func (player *Player) RemoveCardById(id int64) {
	index := -1
	for i, card := range player.Hand {
		if card.ID == id {
			index = i
			break
		}
	}
	player.Hand[index] = player.Hand[len(player.Hand)-1] // Copy last element to index i.
	player.Hand[len(player.Hand)-1] = Card{}   // Erase last element (write zero value).
	player.Hand = player.Hand[:len(player.Hand)-1]
}

func (player *Player) AddToScore(num int) {
	player.Score += num
}

func (player *Player) IsEmpty() bool {
	return len(player.Username) == 0
}

func (player *Player) Equals(p2 Player) bool {
	return player.Username == p2.Username
}