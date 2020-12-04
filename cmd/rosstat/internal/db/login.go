package db

import (
	"path"

	"github.com/recoilme/pudge"
)

type User struct {
	Login    string `json:"login"`
	Role     int    `json:"role"`
	Password string `json:"password"`
}

// depends on user role. if admin - admin panel. otherwise - scan panel
func AuthorizeUser(login, pass string) (string, string, error) {
	defer CloseAllDB()

	// authorization and authentication logic
	user := User{}
	err := pudge.Get(path.Join(".", "db", "users"), login, &user)
	if err != nil {
		return "", "", err
	}

	var addr string
	var apiKey string

	if user.Password == pass {

		if user.Role == 0 {
			apiKey = "admin"
			addr = "admin"
		}else if user.Role == 1 {
			apiKey = "collector"
			addr = "orders"
		} else {
			apiKey = "storekeeper"
			addr = "shipment"
		}
	} else {
		addr = "login"
	}

	return addr, apiKey, nil
}
