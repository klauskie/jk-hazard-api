package repository

import (
	"klaus.com/jkapi/app/database"
	"klaus.com/jkapi/app/entity"
)

type UserRepository interface {
	FindById(key int64) entity.User
	Create(user *entity.User) error
	Save(user *entity.User) error
	FindAll() ([]entity.User, error)
	FindByUsername(username string) (entity.User, error)
}

type repoUser struct{}

// UserRepository
func NewUserRepository() UserRepository {
	return &repoUser{}
}

const (
	projectId = "jkhazardapp"
	collectionName = "players"
	credentialsPath = "/go/src/app/jkhazardapp-firebase-adminsdk.json"
)

func (*repoUser) FindById(key int64) entity.User {
	var user entity.User
	database.DB.First(&user, key)
	return user
}

func (*repoUser) Create(user *entity.User) error {
	return database.DB.Create(user).Error
}

func (*repoUser) Save(user *entity.User) error {
	return database.DB.Save(user).Error
}

func (*repoUser) FindAll() ([]entity.User, error) {
	var users []entity.User
	err := database.DB.Find(&users).Error
	return users, err
}

func (*repoUser) FindByUsername(username string) (entity.User, error) {
	var user entity.User
	err := database.DB.Where("username = ?", username).First(&user).Error
	return user, err
}