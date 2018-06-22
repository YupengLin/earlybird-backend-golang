package main

import (
	"./auth"
	"./chat"
	"./chatroom"
	"./checkin"
	"./profile"

	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // not checking origin
}

func main() {

	/*http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		var result string
		if err := common.DB.QueryRow(`SELECT col FROM test`).Scan(&result); err != nil {
			log.Panic(err)
		}

		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{
			"result":  result,
			"backend": "go",
		}); err != nil {
			log.Panic(err)
		}
	})*/
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", chat.GetMessageListHandler)
	e.GET("/test", chat.GetMessageListHandler)
	e.POST("/api/v1/auth/signup", auth.PostSignupHandler)
	e.GET("/ap1/v1/auth/token", auth.GetToken)
	e.GET("/api/v1/auth/user", auth.GetUser)
	e.POST("/api/v1/checkin/image", checkin.PostImageHandler)
	e.POST("/api/v1/checkin/audio", checkin.PostAudioHanlder)
	e.GET("/api/v1/activity", checkin.GetUserOwnActivityHandler)
	// message api
	//e.POST("/api/v1/messages/", chat.PostMessageHandler)
	e.GET("/api/v1/message_example", chat.GetMessageExampleHandler)

	//profile
	e.POST("/api/v1/profile/thumbnail", profile.PostProfileThumbnailHandler)
	e.POST("/api/v1/profile/username", profile.PostProfileUsernameHandler)

	server := chatroom.NewServer()
	go server.Init()
	e.GET("/ws", func(c echo.Context) error {
		chatroom.Listen(server, c)
		return nil
	})
	e.GET("/api/v1/messages/", func(c echo.Context) error {
		//	server.GetUserList()
		err := chat.GetMessageListHandler(c)
		return err
	})

	e.Logger.Fatal(e.Start(":8000"))

	/*	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}*/
}
