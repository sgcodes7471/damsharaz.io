package socket

import( 
	"encoding/json"
	"github.com/gorilla/websocket"
	"sgcodes7471/damsharaz.io-server/internal/pkg"
	"sgcodes7471/damsharaz.io-server/internal/db"
	"sgcodes7471/damsharaz.io-server/internal/types"
)


func Emit(conn *websocket.Conn , msgObj string) (string , error) {
	event , author , msg , err := pkg.Parse_Payload(msgObj);

	if err != nil {
		return "" , err
	}

	switch event {
	case "START" :
		roomId := msg
		roomObjectData , err := db.Redis_Get(roomId + "_data");
		if err != nil {
			return "" , err;
		}

		var roomObject type.Room_Object;
		err := json.Unmarshal(roomObjectData , &roomObject);

		if err != nil {
			return "" , err;
		}

		var den_client types.Client_Object;
		err = json.Unmarshal(roomObject.Den , &den_client);

		if err != nil {
			return "" , err;
		}

		if den_client.Conn === conn {
			payload = author + "/r/n" + event + "/r/n" + roomObject.Answer + "/r/n";
		} else {
			payload = author + "/r/n" + event + "/r/n" + Den.name + "/r/n";
		}

		return payload , nil;

	default :
		return Errorf("Not a Valid Event")
	}

	return nil
}