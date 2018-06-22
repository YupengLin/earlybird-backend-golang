package checkin

import (
	"encoding/json"
	"log"
	"net/http"

	"../auth"
	"../common"
	model "../models"
	"github.com/labstack/echo"
)

var (
	activityType map[string]int64
)

func init() {
	activityType = make(map[string]int64)
	rows, err := common.DB.Query(`select id, name from checkin_type`)
	if err != nil {
		log.Panic(err.Error())
	}
	for rows.Next() {
		var typeid int64
		var typename string
		err = rows.Scan(&typeid, &typename)
		if err != nil {
			log.Panic(err.Error())
		}
		activityType[typename] = typeid
	}
}

func PostImageHandler(c echo.Context) (err error) {
	userId, _, err := auth.GetUserIdAndPassPartFromRequest(c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "err parsing user id "+err.Error())
	}

	var imgCheckin model.Checkin
	decoder := json.NewDecoder(c.Request().Body)
	err = decoder.Decode(&imgCheckin)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "err in formatting "+err.Error())
	}

	_, err = common.DB.Exec(`INSERT INTO user_checkin_activity (user_id, type_id, resource_url, note, activity_created_at, geog)
    VALUES ($1, $2, $3, $4, $5, ST_GeomFromText($6, 4326))`,
		userId, activityType["image"], imgCheckin.ResourceUrl, imgCheckin.Note, imgCheckin.ActivityCreatedAt, "POINT(-122.431297 37.773972)")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "err in formatting "+err.Error())
	}
	return c.NoContent(http.StatusOK)
}

func isValidCheckin(checkin model.Checkin) bool {
	return true
}

func PostAudioHanlder(c echo.Context) (err error) {
	userId, _, err := auth.GetUserIdAndPassPartFromRequest(c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "err parsing user id "+err.Error())
	}

	var imgCheckin model.Checkin
	decoder := json.NewDecoder(c.Request().Body)
	err = decoder.Decode(&imgCheckin)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "err in formatting "+err.Error())
	}

	_, err = common.DB.Exec(`INSERT INTO user_checkin_activity (user_id, type_id, resource_url, note, activity_created_at, geog)
    VALUES ($1, $2, $3, $4, $5, ST_GeomFromText($6, 4326))`,
		userId, activityType["audio"], imgCheckin.ResourceUrl, imgCheckin.Note, imgCheckin.ActivityCreatedAt, "POINT(-122.431297 37.773972)")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "err in formatting "+err.Error())
	}
	return c.NoContent(http.StatusOK)
}

// func DeleteCheckInHandler(c echo.Context) (err error) {
// 	userId, _, err := auth.GetUserIdAndPassPartFromRequest(c.Request())
// 	if err != nil {
// 		return echo.NewHTTPError(http.StatusUnauthorized, "err parsing user id "+err.Error())
// 	}
//
//
// }
