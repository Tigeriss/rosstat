package main

import (
	"errors"
	"github.com/recoilme/pudge"
	"html/template"
	"log"
	"net/http"
	"path"
	"rosstat/internal"
	"rosstat/internal/session"
)

var (

	tmpl = template.Must(template.ParseFiles(path.Join(".", "ui", "templates", "login.html"),
		path.Join(".", "ui", "templates", "admin.html"),
		path.Join(".", "ui", "templates", "orders.html"),
		path.Join(".", "ui", "templates", "small.html"),
		path.Join(".", "ui", "templates", "o-pallet.html"),
		path.Join(".", "ui", "templates", "shipment.html"),
		path.Join(".", "ui", "templates", "s-pallet.html"),
		path.Join(".", "ui", "templates", "big.html")))
)

type UserSession struct {
	ApiKey      string `json:"apiKey"`
	CurrentUser string `json:"currentUser"`
}


func main() {
	createFirstAdmin()
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/login/enter", handlerLoginCheck)
	http.HandleFunc("/orders", handleOrders)
	http.HandleFunc("/orders/big", handleBigOrder)
	http.HandleFunc("/orders/small", handleSmallOrder)
	http.HandleFunc("/orders/pallet", handlePallet)
	http.HandleFunc("/shipment", handleShipment)
	http.HandleFunc("/shipment/pallet", handleShipmentPallet)
	http.HandleFunc("/admin", handleAdmin)
	http.HandleFunc("/admin/new_user", handleAdminNewUser)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("ui/"))))
	log.Fatal(http.ListenAndServe(":9090", nil))
}

func handleRoot(writer http.ResponseWriter, request *http.Request){
	http.Redirect(writer, request, "/login", http.StatusFound)
}

func handleLogin(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		err := errors.New("don't do it")
		internal.ResponseForbidden(writer, err)
		return
	}
	err := tmpl.ExecuteTemplate(writer, "login.html", nil)
	if err != nil {
		internal.ResponseInternalError(writer, err)
		return
	}
}

func handlerLoginCheck(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		http.Redirect(writer, request, "/login", http.StatusFound)
	} else {
		// trying to load a session from user's request
		var userSession UserSession
		err := session.Load(&userSession, writer, request)
		if err != nil {
			internal.ResponseInternalError(writer, err)
			return
		}

		currentUser := request.FormValue("login")
		passw := request.FormValue("password")
		address, apiKey, err := internal.AuthorizeUser(currentUser, passw)
		if err != nil {
			// incorrect login or password
			internal.ResponseForbidden(writer, err)
			return
		}

		// seems good. we can save the session now
		userSession.CurrentUser = currentUser
		userSession.ApiKey = apiKey
		err = session.Save(userSession, writer, request)
		if err != nil {
			internal.ResponseInternalError(writer, err)
			return
		}

		http.Redirect(writer, request, "/"+address, http.StatusFound)
	}
}

func handleOrders(writer http.ResponseWriter, request *http.Request) {
	err := tmpl.ExecuteTemplate(writer, "orders.html", nil)
	if err != nil {
		internal.ResponseInternalError(writer, err)
		return
	}
}
func handleBigOrder(writer http.ResponseWriter, request *http.Request) {
	err := tmpl.ExecuteTemplate(writer, "big.html", nil)
	if err != nil {
		internal.ResponseInternalError(writer, err)
		return
	}
}
func handleSmallOrder(writer http.ResponseWriter, request *http.Request) {
	err := tmpl.ExecuteTemplate(writer, "small.html", nil)
	if err != nil {
		internal.ResponseInternalError(writer, err)
		return
	}
}
func handlePallet(writer http.ResponseWriter, request *http.Request) {
	err := tmpl.ExecuteTemplate(writer, "o-pallet.html", nil)
	if err != nil {
		internal.ResponseInternalError(writer, err)
		return
	}
}
func handleShipment(writer http.ResponseWriter, request *http.Request) {
	err := tmpl.ExecuteTemplate(writer, "shipment.html", nil)
	if err != nil {
		internal.ResponseInternalError(writer, err)
		return
	}
}
func handleShipmentPallet(writer http.ResponseWriter, request *http.Request) {
	err := tmpl.ExecuteTemplate(writer, "s-pallet.html", nil)
	if err != nil {
		internal.ResponseInternalError(writer, err)
		return
	}
}
func handleAdmin(writer http.ResponseWriter, request *http.Request) {
	// trying to load a session from user's request
	var userSession UserSession
	err := session.Load(&userSession, writer, request)
	if err != nil {
		internal.ResponseInternalError(writer, err)
		return
	}

	if userSession.ApiKey == "admin" {
		data, err := internal.FormData()
		if err != nil {
			internal.ResponseInternalError(writer, err)
			return
		}

		err = tmpl.ExecuteTemplate(writer, "admin.html", data)
		if err != nil {
			internal.ResponseInternalError(writer, err)
			return
		}
	} else {
		http.Redirect(writer, request, "/login", http.StatusFound)
	}
}

func handleAdminNewUser(writer http.ResponseWriter, request *http.Request) {
	// trying to load a session from user's request
	var userSession UserSession
	err := session.Load(&userSession, writer, request)
	if err != nil {
		internal.ResponseInternalError(writer, err)
		return
	}

	if request.Method == "POST" {
		if userSession.ApiKey == "admin" {
			role := 0
			switch request.FormValue("role") {
			case "0":
				role = 0
				break
			case "1":
				role = 1
				break
			case "2":
				role = 2
				break

			}
			log.Println(role)

			_, err := internal.AddUser(request.FormValue("login"), request.FormValue("password"), role)
			if err != nil {
				internal.ResponseInternalError(writer, err)
				return
			}
			http.Redirect(writer, request, "/admin", http.StatusFound)
		} else {
			http.Redirect(writer, request, "/admin", http.StatusFound)
		}
	}
}


func createFirstAdmin() error {
	defer internal.CloseAllDB()
	u := internal.User{
		Login:    "admin",
		Role: 0,
		Password: "admin",
	}
	err := pudge.Set(path.Join(".", "db", "users"), u.Login, u)
	if err != nil {
		return err
	}
	return nil
}
