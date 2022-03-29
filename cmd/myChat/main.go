package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) // connected clients

var broadcast = make(chan Message) // broadcast channel

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Define our message object

type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

func main() {
	fmt.Println(34222222)
	// Create a simple file server
	//	fs := http.FileServer(http.Dir("public"))
	//	http.Handle("/", fs)
	// Configure websocket route
	//	http.HandleFunc("/ws/22", handleConnections)
	r := mux.NewRouter()
	r.HandleFunc("/ws/{key}", handleConnections)
	// Start listening for incoming chat messages
	go handleMessages()

	// Start the server on localhost port 8000 and log any errors
	log.Println("http server started on :3000")

	err := http.ListenAndServe(":3000", r)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	chatId, err := strconv.Atoi(params["key"])
	if err != nil {
		log.Fatal("invalid chat id")
	}

	fmt.Println(chatId)
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Fatal(err)
	}

	// Make sure we close the connection when the function returns
	defer ws.Close()

	// Register our new client
	clients[ws] = true

	fmt.Println(clients)
	for {
		var msg Message
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)

		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}

		// Send the newly received message to the broadcast channel
		broadcast <- msg
	}
}

func handleMessages() {

	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast

		// Send it out to every client that is currently connected
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
