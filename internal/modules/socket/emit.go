package socket

import( 
	"encoding/json"
	"github.com/gorilla/websocket"
	"sgcodes7471/damsharaz.io-server/internal/pkg"
	"sgcodes7471/damsharaz.io-server/internal/db"
	"sgcodes7471/damsharaz.io-server/internal/types"
	"sgcodes7471/damsharaz.io-server/internal/server"
)


func Emit(conn *websocket.Conn , msgObj string) (string , error) {
	event , author , msg , err := pkg.Parse_Payload(msgObj);

	if err != nil {
		return "" , err
	}

	switch event {
	case "START" :
		roomId := msg[:7]
		roomObjectData , err := db.Redis_Get(roomId + "_data");
		if err != nil {
			return "" , err;
		}

		var roomObject types.Room_Object;
		err := json.Unmarshal(roomObjectData , &roomObject);

		if err != nil {
			return "" , err;
		}

		var den_client types.Client_Object = server.Get_Clients_From_Id(roomId , roomObject.Den);

		if den_client.Conn === conn {
			payload = author + "/r/n" + event + "/r/n" + roomObject.Answer + "/r/n";
		} else {
			payload = author + "/r/n" + event + "/r/n" + den_client.name + "/r/n";
		}

		return payload , nil;

	default :
		return Errorf("Not a Valid Event")
	}

	return nil
}