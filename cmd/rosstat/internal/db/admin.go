package db

import (
	"path"

	"github.com/recoilme/pudge"
)

type PublicUser struct {
	Login string `json:"login"`
	Role  string `json:"role"`
}

type User struct {
	Login    string `json:"login"`
	Role     string `json:"role"`
	Password string `json:"password"`
}

type AdminData struct {
	Users []User `json:"users"`
}

var dbPath = path.Join(".", "db", "users")

func GetUsers() ([]User, error) {
	defer pudge.CloseAll()
	keys, err := pudge.Keys(dbPath, 0, 0, 0, true)
	if err != nil {
		return nil, err
	}

	users := make([]User, 0, len(keys))
	for _, key := range keys {
		var u User
		err := pudge.Get(dbPath, key, &u)
		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}

func GetUser(login string) (*User, error) {
	defer pudge.CloseAll()
	var user = &User{}
	err := pudge.Get(dbPath, login, user)
	if err != nil {
		return nil, err
	}
	return user, nil

}

func GetPublicUsers() ([]PublicUser, error) {
	defer pudge.CloseAll()
	keys, err := pudge.Keys(dbPath, 0, 0, 0, true)
	if err != nil {
		return nil, err
	}

	users := make([]PublicUser, 0, len(keys))
	for _, key := range keys {
		var u PublicUser
		err := pudge.Get(dbPath, key, &u)
		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}

func AddUser(login, pass string, role string) error {
	defer pudge.CloseAll()
	u := &User{
		Login:    login,
		Password: pass,
		Role:     role,
	}
	err := pudge.Set(dbPath, u.Login, u)
	if err != nil {
		return err
	}
	return nil
}

func DeleteUser(login string) error {
	defer pudge.CloseAll()
	err := pudge.Delete(dbPath, login)
	if err != nil {
		return err
	}
	return nil
}
