package main

import (
	"context"
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/lib/pq"

	"rosstat/cmd/rosstat/internal/debug"
	"rosstat/cmd/rosstat/internal/handlers"
)

var dbConnStr = "root:whatever@tcp(127.0.0.1:13306)/chat"
var port = 80
var debugFlag = false

func init() {
	flag.BoolVar(&debugFlag, "debug", debugFlag, "proxy to use react-create-app instead of using static handler")
	flag.IntVar(&port, "port", port, "listen port")
	flag.StringVar(&dbConnStr, "db-conn", dbConnStr, "db connection string")
}

func main() {
	flag.Parse()

	appCtx, appCancel := context.WithCancel(context.Background())

	// run react debug server for convenient development
	if debugFlag {
		debug.RunReactDevServer(appCtx)
	}

	// http
	ec := echo.New()

	// some helpful stuff
	ec.Use(middleware.Recover())
	ec.Use(middleware.RequestID())
	ec.Use(handlers.RosContextMiddleware(dbConnStr))
	ec.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `REQU[${time_rfc3339}] ${remote_ip} ${status} ${method} ${uri} ${latency_human} ${error}${message}` + "\n",
	}))
	ec.Use(middleware.CORS())

	// login settings
	ec.POST("/api/login", handlers.Login)

	// static files
	if debugFlag {
		reactUrl, err := url.Parse("http://localhost:3000")
		if err != nil {
			log.Fatalf("unable to parse react url: %s", err)
		}
		// proxy static to webpack for debug mode
		cfg := middleware.ProxyConfig{
			Skipper: func(c echo.Context) bool {
				return strings.HasPrefix(c.Path(), "/api/")
			},
			Balancer: middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{
				{
					Name: "React",
					URL:  reactUrl,
				},
			}),
		}
		ec.Use(middleware.ProxyWithConfig(cfg))
	} else {
		// otherwise use embedded statics
		ec.Use(middleware.Static(path.Join("ui", "build")))
		ec.GET("/*", func(c echo.Context) error {
			return c.File(path.Join("ui", "build", "index.html"))
		})
	}

	// everything else is restricted
	api := ec.Group("/api")
	api.Use(middleware.JWT([]byte(handlers.LoginSecretKey)))

	// ---- ALL API ENDPOINTS ARE DEFINED HERE !!! ----
	api.GET("/orders/build", handlers.GetToBuildOrders)
	api.GET("/orders/big/build/:id", handlers.GetBigToBuildOrders)
	api.GET("/orders/small/build/:id", handlers.GetSmallToBuildOrders)
	api.POST("/orders/small/build/:id/finish", handlers.FinishSmallToBuildOrders)
	api.GET("/orders/big/pallet/:id", handlers.GetBigPalletOrders)
	api.GET("/orders/big/pallet/:id/num/:num", handlers.GetBigPalletNum)
	api.GET("/orders/big/pallet/:id/barcode/:barcode", handlers.GetBigPalletBarcodeOrders)
	api.POST("/orders/big/pallet/:id/finish", handlers.FinishBigPalletOrders)
	api.GET("/shipment/ready", handlers.GetReadyForShipment)
	api.GET("/shipment/pallet/:id", handlers.GetPalletShipment)
	api.POST("/shipment/pallet/:id/finish", handlers.FinishPalletShipment)
	api.GET("/admin/users", handlers.GetUsers)
	api.POST("/admin/users", handlers.AddUser)
	api.DELETE("/admin/users/:login", handlers.DeleteUser)
	// ------------------------------------------------

	// start server
	go func() {
		if err := ec.Start(":" + strconv.Itoa(port)); err != nil {
			log.Println("shutting down the server")
		}
	}()

	// -- below is just some graceful stop stuff

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// stop all app related jobs like db connections
	appCancel()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// jobsQuit <- struct{}{}
	if err := ec.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	time.Sleep(2 * time.Second)
}

//
// func handleRoot(writer http.ResponseWriter, request *http.Request){
// 	http.Redirect(writer, request, "/login", http.StatusFound)
// }
//
// func handleLogin(writer http.ResponseWriter, request *http.Request) {
// 	if request.Method != "GET" {
// 		err := errors.New("don't do it")
// 		db.ResponseForbidden(writer, err)
// 		return
// 	}
// 	err := tmpl.ExecuteTemplate(writer, "login.html", nil)
// 	if err != nil {
// 		db.ResponseInternalError(writer, err)
// 		return
// 	}
// }
//
// func handlerLoginCheck(writer http.ResponseWriter, request *http.Request) {
// 	if request.Method == "GET" {
// 		http.Redirect(writer, request, "/login", http.StatusFound)
// 	} else {
// 		// trying to load a session from user's request
// 		var userSession UserSession
// 		err := session.Load(&userSession, writer, request)
// 		if err != nil {
// 			db.ResponseInternalError(writer, err)
// 			return
// 		}
//
// 		currentUser := request.FormValue("login")
// 		passw := request.FormValue("password")
// 		address, apiKey, err := db.AuthorizeUser(currentUser, passw)
// 		if err != nil {
// 			// incorrect login or password
// 			db.ResponseForbidden(writer, err)
// 			return
// 		}
//
// 		// seems good. we can save the session now
// 		userSession.CurrentUser = currentUser
// 		userSession.ApiKey = apiKey
// 		err = session.Save(userSession, writer, request)
// 		if err != nil {
// 			db.ResponseInternalError(writer, err)
// 			return
// 		}
//
// 		http.Redirect(writer, request, "/"+address, http.StatusFound)
// 	}
// }
//
// func handleOrders(writer http.ResponseWriter, request *http.Request) {
// 	err := tmpl.ExecuteTemplate(writer, "orders.html", nil)
// 	if err != nil {
// 		db.ResponseInternalError(writer, err)
// 		return
// 	}
// }
// func handleBigOrder(writer http.ResponseWriter, request *http.Request) {
// 	err := tmpl.ExecuteTemplate(writer, "big.html", nil)
// 	if err != nil {
// 		db.ResponseInternalError(writer, err)
// 		return
// 	}
// }
// func handleSmallOrder(writer http.ResponseWriter, request *http.Request) {
// 	err := tmpl.ExecuteTemplate(writer, "small.html", nil)
// 	if err != nil {
// 		db.ResponseInternalError(writer, err)
// 		return
// 	}
// }
// func handlePallet(writer http.ResponseWriter, request *http.Request) {
// 	err := tmpl.ExecuteTemplate(writer, "o-pallet.html", nil)
// 	if err != nil {
// 		db.ResponseInternalError(writer, err)
// 		return
// 	}
// }
// func handleShipment(writer http.ResponseWriter, request *http.Request) {
// 	err := tmpl.ExecuteTemplate(writer, "shipment.html", nil)
// 	if err != nil {
// 		db.ResponseInternalError(writer, err)
// 		return
// 	}
// }
// func handleShipmentPallet(writer http.ResponseWriter, request *http.Request) {
// 	err := tmpl.ExecuteTemplate(writer, "s-pallet.html", nil)
// 	if err != nil {
// 		db.ResponseInternalError(writer, err)
// 		return
// 	}
// }
// func handleAdmin(writer http.ResponseWriter, request *http.Request) {
// 	// trying to load a session from user's request
// 	var userSession UserSession
// 	err := session.Load(&userSession, writer, request)
// 	if err != nil {
// 		db.ResponseInternalError(writer, err)
// 		return
// 	}
//
// 	if userSession.ApiKey == "admin" {
// 		data, err := db.FormData()
// 		if err != nil {
// 			db.ResponseInternalError(writer, err)
// 			return
// 		}
//
// 		err = tmpl.ExecuteTemplate(writer, "admin.html", data)
// 		if err != nil {
// 			db.ResponseInternalError(writer, err)
// 			return
// 		}
// 	} else {
// 		http.Redirect(writer, request, "/login", http.StatusFound)
// 	}
// }
//
// func handleAdminNewUser(writer http.ResponseWriter, request *http.Request) {
// 	// trying to load a session from user's request
// 	var userSession UserSession
// 	err := session.Load(&userSession, writer, request)
// 	if err != nil {
// 		db.ResponseInternalError(writer, err)
// 		return
// 	}
//
// 	if request.Method == "POST" {
// 		if userSession.ApiKey == "admin" {
// 			role := 0
// 			switch request.FormValue("role") {
// 			case "0":
// 				role = 0
// 				break
// 			case "1":
// 				role = 1
// 				break
// 			case "2":
// 				role = 2
// 				break
//
// 			}
// 			log.Println(role)
//
// 			_, err := db.AddUser(request.FormValue("login"), request.FormValue("password"), role)
// 			if err != nil {
// 				db.ResponseInternalError(writer, err)
// 				return
// 			}
// 			http.Redirect(writer, request, "/admin", http.StatusFound)
// 		} else {
// 			http.Redirect(writer, request, "/admin", http.StatusFound)
// 		}
// 	}
// }
//
//
// func createFirstAdmin() error {
// 	defer db.CloseAllDB()
// 	u := db.User{
// 		Login:    "admin",
// 		Role: 0,
// 		Password: "admin",
// 	}
// 	err := pudge.Set(path.Join(".", "db", "users"), u.Login, u)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
