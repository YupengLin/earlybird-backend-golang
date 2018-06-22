package chat

import (
	"database/sql"
	"log"

	"../common"
	model "../models"
	//"strings"

	"github.com/labstack/echo"
)

func GetMessageListHandler(c echo.Context) (err error) {
	rows, err := common.DB.Query(`select m.message_uuid, m.message_id,m.user_id,m.message_content,m.created_at,
		u.username,message_type from message m
		join user_ u on u.user_id =  m.user_id
		where m.message_content !=""
		order by created_at asc limit 100`)
	if err != nil && err != sql.ErrNoRows {
		return c.JSON(400, err.Error())
	}
	defer rows.Close()
	messages := []model.Message{}
	for rows.Next() {
		m := model.Message{}
		rows.Scan(&m.UUID, &m.MessageID, &m.UserID, &m.MessageContent,
			&m.CreatedAt, &m.Username, &m.MessageType)
		messages = append(messages, m)
	}
	return c.JSON(200, messages)
}

func GetAllMessageList() (messages []*model.Message, err error) {
	rows, err := common.DB.Query(`select m.uuid,m.message_id,m.user_id,m.message_content,m.created_at,
		u.username,message_type from message m
		join user_ u on u.user_id =  m.user_id
		order by created_at asc limit 100`)
	if err != nil && err != sql.ErrNoRows {
		return
	}
	defer rows.Close()
	messages = []*model.Message{}
	for rows.Next() {
		m := model.Message{}
		rows.Scan(&m.UUID, &m.MessageID, &m.UserID, &m.MessageContent,
			&m.CreatedAt, &m.Username, &m.MessageType)
		messages = append(messages, &m)
	}
	return messages, err
}

func CreateNewMessage(message *model.Message) (err error) {
	_, err = common.DB.Exec(`insert into message(message_uuid,user_id,message_type,message_content,created_at)
		values($1, $2, $3, $4,now())`, message.UUID, message.UserID, message.MessageType, message.MessageContent)
	if err != nil {
		log.Print(err.Error())
		return err
	}
	return
}

func GetMessageExampleHandler(c echo.Context) (err error) {
	message := model.Message{}
	username := "yupeng"
	message.Username = &username
	messageContent := "test text message"
	message.MessageContent = &messageContent
	return c.JSON(200, message)
}
