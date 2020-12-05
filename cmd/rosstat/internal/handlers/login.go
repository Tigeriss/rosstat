package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"

	"rosstat/cmd/rosstat/internal/db"
)

const LoginSecretKey = "pleasechangeme"

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Login string `json:"login"`
	Role  string `json:"role"`
	Token string `json:"token"`
}

func Login(c echo.Context) error {
	req := new(LoginRequest)

	if err := c.Bind(req); err != nil {
		return echo.ErrBadRequest
	}

	// replace it with real db
	user, err := db.GetUser(req.Login)
	if err != nil {
		log.Println(err)
		return echo.ErrUnauthorized
	}

	if user == nil {
		return echo.ErrUnauthorized
	}

	if req.Password != user.Password {
		return echo.ErrUnauthorized
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["login"] = user.Login
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(LoginSecretKey))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, LoginResponse{
		Login: req.Login,
		Role:  user.Role,
		Token: t,
	})
}
