package types

import(
	"github.com/gorilla/websocket"
)

type Client_Object struct {
	Id string
	Conn *websocket.Conn 
	Name string
}

type Room_Object struct {
	RoomId string 
	Token string
	Den string
	Ongoing bool
	Answer string
}


