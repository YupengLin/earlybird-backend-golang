package auth

import (
	"database/sql"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"../common"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type jwtClaims struct {
	UserId   int64
	PassPart string
	jwt.StandardClaims
}

func GetToken(c echo.Context) (err error) {
	email := strings.ToLower(c.QueryParam("email"))
	password := c.QueryParam("password")

	if len(email) == 0 || len(password) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "email and password required")
	}

	var userid int64
	var dbPassword string
	var salt string

	err = common.DB.QueryRow(`select id, salt, password from user_ where email=$1`, email).Scan(&userid, &salt, &dbPassword)
	if err == sql.ErrNoRows {
		time.Sleep(time.Millisecond * time.Duration(500+rand.Intn(200))) // wait to prevent against dictionary attacks
		return echo.NewHTTPError(http.StatusUnauthorized, "Please provide valid credentials")

	} else if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())

	}

	var computedPasswordHash string
	computedPasswordHash, _, err = ComputePasswordHashAndSaltByPasswordAndSaltAndVersion(password, &salt, password_version)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if computedPasswordHash != dbPassword {
		return echo.NewHTTPError(http.StatusBadRequest, "Cannot compute password hash and salt.")
	}

	claims := jwtClaims{
		userid,
		computedPasswordHash[0:16],
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 216000).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString(common.Config.JwtSecret)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusOK, tokenString)
}
