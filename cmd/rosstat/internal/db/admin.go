package db

import (
	"path"

	"github.com/recoilme/pudge"
)

type AdminData struct {
	Users   []User `json:"users"`
}

func getUsers() ([]User, error) {
	defer CloseAllDB()
	keys, err := pudge.Keys(path.Join(".", "db", "users"), 0, 0, 0, true)
	if err != nil {
		return nil, err
	}

	users := make([]User, 0, len(keys))
	for _, key := range keys {
		var u User
		err := pudge.Get(path.Join(".", "db", "users"), key, &u)
		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}

func FormData() (AdminData, error) {
	users, err := getUsers()
	if err != nil {
		return AdminData{}, err
	}
	return AdminData{
		Users:   users,
	}, nil
}

func AddUser(login, pass string, role int) (AdminData, error) {
	defer CloseAllDB()
	u := &User{
		Login:    login,
		Password: pass,
		Role:    role,
	}
	err := pudge.Set(path.Join(".", "db", "users"), u.Login, u)
	if err != nil {
		return AdminData{}, err
	}

	return FormData()
}

func DeleteUser(login string) error {
	defer pudge.CloseAll()
	err := pudge.Delete(path.Join(".", "db", "users"), login)
	if err != nil {
		return err
	}
	return nil
}
