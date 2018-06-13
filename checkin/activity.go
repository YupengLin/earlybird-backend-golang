package checkin


import (

	// "encoding/json"
	"log"
	"net/http"

	"../auth"
	"../common"
	model "../models"
	"github.com/labstack/echo"

)

func GetUserOwnActivityHandler(c echo.Context) (err error) {
    log.Print("get user own activites")
    userId, _, err := auth.GetUserIdAndPassPartFromRequest(c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "err parsing user id "+err.Error())
	}

    var activities []model.Checkin

    rows, err := common.DB.Query(`select uca.activity_created_at, ST_AsText(uca.geog) geo_location, uca.note, uca.resource_url, ct.name from user_checkin_activity uca
    join checkin_type ct on ct.id=uca.type_id
    where uca.user_id = $1 and uca.deleted_at IS NULL order by uca.created_at desc`, userId)
    defer rows.Close()
    if err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "err in db "+err.Error())
    }

    for rows.Next() {
        var ac model.Checkin
        err = rows.Scan(&ac.ActivityCreatedAt, &ac.Location,  &ac.Note, &ac.ResourceUrl, &ac.Type)
        if err != nil {
            return echo.NewHTTPError(http.StatusBadRequest, "err in scan "+err.Error())
        }
        activities = append(activities, ac)
    }
    return c.JSON(http.StatusOK, activities)
}
