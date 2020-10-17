package entity


type User struct {
	ID int64 `json:"-"`
	Username string `gorm:"unique;not null" json:"username"`
	Password string `gorm:"not null" json:"-"`
}

func (user *User) IsEmpty() bool {
	return user.Username == "" || user.Password == ""
}