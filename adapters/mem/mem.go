package mem

import (
	"errors"
	"nononsensecode.com/session-management/domain/model"
)

var (
	users []model.User = []model.User{
		{
			ID:       1,
			Name:     "Kaushik",
			UserName: "kaushikam",
			Password: "redhat",
		},
		{
			ID:       2,
			Name:     "Ashif",
			UserName: "ashif",
			Password: "redhat",
		},
	}
)

type UserRepo struct {}

func (r UserRepo) Find(username, password string) (model.User, error) {
	for i, u := range users {
		if u.UserName == username && u.Password == password {
			return users[i], nil
		}
	}

	return model.User{}, errors.New("there is no such user")
}

func (r UserRepo) FindByUsername(username string) (model.User, error) {
	for i, u := range users {
		if u.UserName == username {
			return users[i], nil
		}
	}

	return model.User{}, errors.New("there is no such user")
}