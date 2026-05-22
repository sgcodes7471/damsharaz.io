package socket

import(
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"sgcodes7471/damsharaz.io-server/internal/pkg"
	"sgcodes7471/damsharaz.io-server/internal/db"
	"sgcodes7471/damsharaz.io-server/internal/types"
)


func Publish(payload string , r *http.Request) error {
	event , author, msg , err := pkg.Parse_Payload(payload);

	if err != nil {
		return err
	}

	switch event {
	case "START" :
		roomId = strings.TrimSpace(msg[:7]);
		token = strings.TrimSpace(msg[7:]);

		if token == "" {
			return fmt.Errorf("Unauthorized access");
		}

		token_claims , err := pkg.VerifyToken(token);
		if err != nil {
			return err;
		}

		if token_claims["roomId"] != roomId {
			return fmt.Errorf("Unauthorized access");
		}

		var den_client_Id string;
		den_client_Id , err = db.Redis_Random(roomId + "_members") ;
		if err != nil {
			return err;
		} 

		// choose a random word here .
		
		var ans string = "";

		roomObject := types.Room_Object{
			RoomId : r.Header.Get("roomId") ,
			Token : token ,
			Den : den_client_Id ,
			Ongoing : true ,
			Answer : ans
		}

		var data string;
		data , err = json.Marshal(roomObject);

		if err != nil {
			return err;
		}

		if err := db.Redis_Set(roomId + "_data", data) ; err != nil {
			return err;
		}

		if err := db.Redis_Publish(roomId , payload) ; err != nil {
			return err;
		}

	default :
		return Errorf("Not a Valid Event")
	}

	return nil;
}


