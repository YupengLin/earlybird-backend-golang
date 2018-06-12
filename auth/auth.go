package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"../common"
	model "../models"
	jwt "github.com/dgrijalva/jwt-go"
	jwtReq "github.com/dgrijalva/jwt-go/request"
	"github.com/labstack/echo"
)

type jwtClaims struct {
	UserId   int64
	PassPart string
	jwt.StandardClaims
}

var (
	ErrJwtClaimsAssertFailed = errors.New("couldn't assert claim type")
	jwtKeyFunc               = func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return common.Config.JwtSecret, nil
	}
)

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

func GetUser(c echo.Context) (err error) {
	userId, passPart, err := GetUserIdAndPassPartFromRequest(c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "err parsing user id "+err.Error())
	}

	user, err := GetUserByIdAndPassPart(userId, passPart)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "err parsing user obj "+err.Error())
	}
	return c.JSON(http.StatusOK, user)

}

func GetUserIdAndPassPartFromRequest(r *http.Request) (userID int64, passPart string, err error) {
	token, err := jwtReq.ParseFromRequestWithClaims(r, jwtReq.AuthorizationHeaderExtractor, &jwtClaims{}, jwtKeyFunc)
	if err != nil {
		return
	}

	claims, ok := token.Claims.(*jwtClaims)
	if !ok {
		err = ErrJwtClaimsAssertFailed
		return
	}

	userID = claims.UserId
	passPart = claims.PassPart

	return
}

func GetUserByIdAndPassPart(id int64, passPart string) (u model.User, err error) {
	var password string
	err = common.DB.QueryRow(`select username, password from user_ where id=$1`, id).Scan(&u.Username, &password)
	if err != nil {
		return
	}

	if len(password) >= 16 && passPart != password[0:16] {
		err = ErrPasspartMismatch
	}

	return
}
