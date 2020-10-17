package entity


type Player struct {
	UserData User
	Score int
	IsHost bool
	Hand []Card
}

func NewPlayer(user User, isHost bool) Player {
	player := Player{user, 0, isHost, []Card{}}
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
	return player.UserData.IsEmpty()
}