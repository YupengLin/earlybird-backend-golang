package auth

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/crypto/pbkdf2"

	"../common"
	model "../models"
	"github.com/labstack/echo"
)

var (
	saltMax          *big.Int = big.NewInt(9223372036854775807)
	password_version int64    = 1
)

func PostSignupHandler(c echo.Context) (err error) {
	var u model.User
	// if err = c.Bind(u); err != nil {
	// 	return echo.NewHTTPError(http.StatusBadRequest, "err in parsing format"+err.Error())
	// }

	decoder := json.NewDecoder(c.Request().Body)
	err = decoder.Decode(&u)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "err in formatting "+err.Error())
	}

	u.Email = strings.TrimSpace(strings.ToLower(u.Email))
	if u.Password == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "password required")
	}

	passwordHash, saltStr, err := ComputePasswordHashAndSaltByPasswordAndSaltAndVersion(*u.Password, nil, password_version)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	_, err = common.DB.Exec(`INSERT INTO user_ (email, username, password, salt) VALUES ($1, $2, $3, $4)`, u.Email, u.Username, passwordHash, saltStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func ComputePasswordHashAndSaltByPasswordAndSaltAndVersion(password string, saltStr *string, version int64) (hash string, salt string, err error) {
	if saltStr == nil {
		if version == 1 {
			var saltBigInt *big.Int
			saltBigInt, err = rand.Int(rand.Reader, saltMax)
			if err != nil {
				return
			}
			salt = saltBigInt.String()
			saltStr = &salt
		} else {
			err = errors.New("users: ComputePasswordHash: a salt is required for password version 0")
		}
	}

	switch version {
	case 0: // old magento password hashing....md5? pffff
		md5Sum := md5.Sum([]byte(*saltStr + password))
		hash = hex.EncodeToString(md5Sum[:])
	case 1: // hack this fuckers
		var saltInt64 int64
		saltInt64, err = strconv.ParseInt(*saltStr, 10, 64)
		if err != nil {
			return
		}
		dk := pbkdf2.Key([]byte(password), big.NewInt(saltInt64).Bytes(), 20000, 512, sha512.New)
		hash = base64.StdEncoding.EncodeToString(dk)
	default:
		err = errors.New("users: password version not supported")
	}
	return
}
