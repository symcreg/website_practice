package model

// TODO: add searcher

type User struct {
	ID           int `gorm:"primary_key"`
	Name         string
	Email        string
	Password     string
	ResisterTime string
}
type BlackList struct {
	ID    int `gorm:"primary_key"`
	Token string
}
type UserDB interface {
	InsertUser(user *User) error
	GetUserByEmail(email string) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(user *User) error
	IsExist(email string) bool
}
type BlackListDB interface {
	InsertBlackList(token string) error
	IsTokenBlackListed(token string) bool
}
