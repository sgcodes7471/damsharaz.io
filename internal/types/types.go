package types

import(
	"github.com/gorilla/websocket"
)

type Client_Object struct {
	Conn *websocket.Conn 
	Name string
}

type Room_Object struct {
	RoomId string 
	Token string
	Den Client_Object
	Ongoing bool
	Answer string
	// add a new field of Answer which will not be compusolory
}


