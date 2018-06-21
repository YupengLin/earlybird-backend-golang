package profile

import (
	"encoding/json"
	"net/http"

	"../auth"
	"../common"

	model "../models"
	"github.com/labstack/echo"
)

func PostProfileThumbnailHandler(c echo.Context) (err error) {
	userId, _, err := auth.GetUserIdAndPassPartFromRequest(c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	var u model.User
	decoder := json.NewDecoder(c.Request().Body)
	err = decoder.Decode(&u)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "err in formatting "+err.Error())
	}
	_, err = common.DB.Exec(`update user_ set img_url=$1 where id=$2`, u.Thumbnail, userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "err in db "+err.Error())
	}

	return c.NoContent(http.StatusOK)

}

func PostProfileUsernameHandler(c echo.Context) (err error) {
	userId, _, err := auth.GetUserIdAndPassPartFromRequest(c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	var u model.User
	decoder := json.NewDecoder(c.Request().Body)
	err = decoder.Decode(&u)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "err in formatting "+err.Error())
	}
	_, err = common.DB.Exec(`update user_ set username=$1 where id=$2`, u.Username, userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "err in db "+err.Error())
	}

	return c.NoContent(http.StatusOK)
}
