package chat

import (
	"database/sql"
	"log"

	"../common"
	//"net/http"
	//	model "app/models"
	//	"github.com/labstack/echo"
)

func CheckUserExist(username string) (userId *int64) {

	err := common.DB.QueryRow(`select id from user_ where username=$1`,
		username).Scan(&userId)
	if err != nil && err != sql.ErrNoRows {
		log.Print(err.Error())
		return
	}

	return
}

func UpdateUser(username string, newName string) {
	_, err := common.DB.Exec(`update user_ set username = ? where username = ?`, newName, username)
	if err != nil && err != sql.ErrNoRows {
		return
	}
	return
}

func CreateNewUser(username string) (userId int64) {
	uid := CheckUserExist(username)
	if uid == nil {
		ins, err := common.DB.Prepare(`insert into user_(username,created_at,last_login)
			values(?,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP())`)
		if err != nil {
			return
		}
		res, err := ins.Exec(username)
		if err != nil {
			return
		}
		userId, _ = res.LastInsertId()
	} else {
		userId = *uid
	}
	return userId
}
