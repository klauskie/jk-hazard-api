package repository


import (
	"klaus.com/jkapi/app/entity"
)

var instance *repoRoom = nil

type RoomRepository interface {
	Add(room *entity.Room)
	Get(tag string) *entity.Room
	Remove(tag string)
	Save(room *entity.Room)
	GetAll()	map[string]*entity.Room
}

type repoRoom struct {
	collection		map[string]*entity.Room
}

// RoomRepository
func GetRoomRepository() RoomRepository {
	if instance == nil {
		instance = new(repoRoom)
		instance.collection = make(map[string]*entity.Room)
	}
	return instance
}

// Add Room to collection
func (repo *repoRoom) Add(room *entity.Room) {
	repo.collection[room.TAG] = room
}

// Get Room in collection by room Tag
func (repo *repoRoom) Get(key string) *entity.Room {
	if _, ok := repo.collection[key]; ok {
		return repo.collection[key]
	}
	return nil
}

// Save room in collection
func (repo *repoRoom) Save(room *entity.Room) {
	if _, ok := repo.collection[room.TAG]; ok {
		repo.collection[room.TAG] = room
	}
}

// Remove room in collection
func (repo *repoRoom) Remove(key string) {
	if _, ok := repo.collection[key]; ok {
		delete(repo.collection, key);
	}
}

func (repo *repoRoom) GetAll() map[string]*entity.Room {
	return repo.collection
}