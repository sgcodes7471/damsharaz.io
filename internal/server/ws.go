package server

import( 
	"net/http"
	"encoding/json"
	"sgcodes7471/damsharaz.io-server/internal/pkg"
	"sgcodes7471/damsharaz.io-server/internal/db"
	"sgcodes7471/damsharaz.io-server/internal/types"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

var Rooms_Connections = make(
	map[string]map[string]types.Client_Object,
)

func Get_Clients_From_Id(roomId string, clientId string) types.Client_Object {
	return Rooms_Connections[roomId][clientId];
}

var upgrader = websocket.Upgrader{};

func WSServer(w http.ResponseWriter , r *http.Request) {
	roomId := r.Header.Get("roomId");
	clientId := r.Header.Get("clientId");
	name := r.URL.Query().Get("name");

	if roomId == "" || name == "" {
		pkg.Api_Error("Required Query Params missing in /ws" , "GET /ws" , 400 , w);
		return
	}

	for _, _ch := range name {
		if _ch == '/' {
			pkg.Api_Error("Name has '/' : FORBIDDEN" , "GET /ws" , 403 , w);
			return
		}
	}

	var value string;
	value , err := db.Redis_Get(roomId);

	if  err == redis.Nil {
		pkg.Api_Error("RoomId invalid! Room NOT FOUND in Redis" , "GET /ws" , 404 , w);
		return
	} else if err != nil {
		pkg.Api_Error("Error in connecting to Redis : " + err.Error() , "GET /ws" , 500 , w);
		return
	}

	
	var Room = types.Room_Object{
		RoomId : roomId ,
		Token : value ,
	};

	data , err := json.Marshal(Room);
	
	if err != nil {
		pkg.Api_Error("Error in Serializing Room Object : " + err.Error() , "GET /ws" , 500 , w);
		return
	}

	key := roomId + "_data";
	if _ , err := db.Redis_Get(key); err == redis.Nil {
		err = db.Redis_Set(key , string(data) , 0);
		
		if err != nil {
			pkg.Api_Error("Error in writing Room Object to Redis : " + err.Error() , "GET /ws" , 500 , w);
			return
		}
	} else if err != nil {
		pkg.Api_Error("Error in connecting to Redis : " + err.Error() , "GET /ws" , 500 , w);
		return
	}

	
	conn, err := upgrader.Upgrade(w, r, nil);

	if err != nil {
		pkg.Api_Error("Websocket did not start : " + err.Error() , "GET /ws" , 500 , w);
		return
	}
	defer conn.Close();


	sub := db.Redis_Client.Subscribe(db.CTX , roomId);
	defer sub.Close();

	ch := sub.Channel();

	client := types.Client_Object{
		Id : clientId ,
		Conn : conn ,
		Name : name ,
	}

	if Rooms_Connections[roomId] == nil {
		Room_Connections[roomId] = make(map[string]types.Client_Object);
	}

	Rooms_Connections[roomId][clientId] = client;

	key = roomId + "_members";
	if err := db.Redis_Client.SAdd(db.CTX , key , clientId).Err(); err != nil {
		pkg.Api_Error("Error in adding member to Redis : " + err.Error() , "GET /ws" , 500 , w);
		return
	}


	go func() {
		_, msg, err := conn.ReadMessage();

		if(err != nil) {
			pkg.Api_Error("Error reading in /ws from connection " + name + " : " + err.Error() , "GET /ws" , 500 , w);
			return
		}

		finalMsgPayload := name + "/r/n" + string(msg);

		if err := db.Redis_Publish(roomId , finalMsgPayload) ; err != nil {
			pkg.Api_Error("Error publishing to Redis in /ws from connection " + name + " : " + err.Error() , "GET /ws" , 500 , w);
			return
		}
	} ();


	for msg := range ch {
		if err := conn.WriteMessage(websocket.TextMessage , []byte(msg.Payload)); err != nil {
			pkg.Api_Error("Error writing from Redis in /ws to connection " + name + " : " + err.Error() , "GET /ws" , 500 , w);
			return
		}
	}
}