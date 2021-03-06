package chatroom

import (
	"log"
	"strings"

	model "../models"
	//	"math"
	"net/http"
	"time"

	chat "../chat"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/twinj/uuid"
)

type ChatServer struct {
	OnlineUsers  map[string]Client
	NewMessage   chan *model.Message
	OfflineUsers map[string]Client
	NewUser      chan *Client
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true }, // not checking origin
	}
)

func NewServer() (server *ChatServer) {
	return &ChatServer{
		//	AllMessages:  []*model.Message{},
		NewMessage:   make(chan *model.Message, 5),
		OnlineUsers:  make(map[string]Client),
		OfflineUsers: make(map[string]Client),
		NewUser:      make(chan *Client, 5),
	}
}

// Initializing the chatroom
func (server *ChatServer) Init() {
	go func() {
		for {
			server.BroadCast()
			time.Sleep(100 * time.Millisecond)
		}
	}()
}

func (server *ChatServer) Join(msg model.Message, conn *websocket.Conn) *Client {

	client := &Client{
		User:   model.User{Username: msg.Username, UserID: msg.UserID},
		Socket: conn,
		Server: server,
	}

	server.OnlineUsers[*msg.Username] = *client
	server.updateOnlineUserList(client)
	server.AddMessage(model.Message{
		UUID:           uuid.NewV4().String(),
		MessageType:    "system-message",
		CreatedAt:      time.Now(),
		MessageContent: getPointer(*msg.Username + " has joined the chat."),
		User:           model.User{UserID: 0, Username: getPointer("system")},
	})

	server.AddMessage(msg)

	return client
}

// Leaving the chatroom
func (server *ChatServer) Leave(name string) {
	server.OfflineUsers[name] = server.OnlineUsers[name]
	delete(server.OnlineUsers, name)

	server.AddMessage(
		model.Message{
			UUID:           uuid.NewV4().String(),
			MessageType:    "system-message",
			CreatedAt:      time.Now(),
			MessageContent: getPointer(name + " has left the chat."),
			User:           model.User{UserID: 0, Username: getPointer("system")},
		})
}

func getPointer(s string) *string {
	return &s
}

// Adding message to queue
func (server *ChatServer) AddMessage(message model.Message) {
	server.NewMessage <- &message
	if message.Username != nil {
		if message.UserID != 0 {
			uid := chat.CheckUserExist(*message.Username)
			if uid != nil {
				message.UserID = *uid

			}
		}
		chat.CreateNewMessage(&message)
	}
}

func Listen(server *ChatServer, c echo.Context) error {
	//c.GET("/ws", server.Listen)
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	if err != nil {
		ws.Close()
		return err
	}
	defer ws.Close()
	log.Print("websocket start")
	msg := model.Message{}
	err = ws.ReadJSON(&msg)
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			log.Printf("WEB SOCKET:error: %v", err)
		} else {
			log.Print("WEB SOCKET: other err" + err.Error())
		}
		ws.Close()
		return err
	}
	msg.UUID = uuid.NewV4().String()

	log.Print("WEB SOCKET: FIRST MESSAGE: ", *msg.MessageContent)
	user := server.Join(msg, ws)

	if user == nil {
		//	log.Print(err)
		return err
	}

	for {
		msg := model.Message{}
		err = ws.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			} else {
				log.Print(err)
			}
			user.Exit()
			server.updateOnlineUserList(user)
			return err
		}
		msg.UUID = uuid.NewV4().String()
		log.Print(*msg.MessageContent)
		if user.Username != nil && msg.Username != nil {
			if strings.TrimSpace(*user.Username) != strings.TrimSpace(*msg.Username) {
				chat.UpdateUser(*user.Username, *msg.Username)
				delete(server.OnlineUsers, *msg.Guestname)
				user.Username = msg.Username
				user.Register = msg.Register
				user.Guestname = msg.Guestname
				server.OnlineUsers[*user.Username] = *user
				server.updateOnlineUserList(user)
			}

		}
		// Write
		user.NewMessage(msg)
	}
	//return err
}

func (server *ChatServer) updateOnlineUserList(client *Client) {
	server.NewUser <- client

}

// Broadcasting all the messages in the queue in one block
func (server *ChatServer) BroadCast() {

	messages := make([]*model.Message, 0)
	userList := make(map[string]interface{})
InfiLoop:
	for {
		select {
		case message := <-server.NewMessage:
			messages = append(messages, message)
			for _, client := range server.OnlineUsers {
				client.SendSingleMessage(message)
			}
		case <-server.NewUser:
			//	user["username"] = *newUser.Username
			//	user["created_at"] = time.Now().String()
			userList["message_type"] = "user_list"
			us := []model.User{}
			for _, c := range server.OnlineUsers {
				us = append(us, c.User)
			}
			userList["list"] = us
		default:
			break InfiLoop
		}
	}
	// if len(userList) > 0 {
	// 	for _, client := range server.OnlineUsers {
	// 		client.Socket.WriteJSON([]map[string]interface{}{
	// 			userList,
	// 		})
	// 	}
	// }

	// if len(messages) > 0 {
	// 	for _, client := range server.OnlineUsers {
	// 		client.Send(messages)
	// 	}
	// }
}
