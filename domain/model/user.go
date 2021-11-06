package model

type User struct {
	ID int
	Name string
	UserName string
	Password string
}

type UserRepository interface {
	Find(username, password string) (User, error)
}