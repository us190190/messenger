package services

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
	"log"
	"messenger/contracts"
	"messenger/models"
	"net/http"
	"sync"
)

var (
	connectionUpgrade = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	connections       = make(map[int]*websocket.Conn) // Map to store WebSocket connections
	mutex             = sync.Mutex{}                  // Mutex for thread-safe access to the connections map
	activeConnections = sync.WaitGroup{}
)

func HandleWebSocket(responseWriter http.ResponseWriter, request *http.Request) {

	// Upgrade HTTP connection to WebSocket
	conn, err := connectionUpgrade.Upgrade(responseWriter, request, nil)
	if err != nil {
		log.Println("Error upgrading connection to WebSocket:", err)
		return
	}
	defer conn.Close()

	// Read authentication message
	_, authMsg, err := conn.ReadMessage()
	if err != nil {
		log.Println("Error reading authentication message from WebSocket:", err)
		return
	}

	// Decode authentication message
	var authData contracts.CreateUserRequest
	err = json.Unmarshal(authMsg, &authData)
	if err != nil {
		log.Println("Error decoding authentication message:", err)
		return
	}

	// Perform authentication (e.g., check credentials against database)
	// Retrieve user from database
	currentUser, err := models.GetUserByUsername(authData.Username)
	if err != nil {
		log.Println("Authentication failed for user:", authData.Username)
		return
	}
	storedPassword := currentUser.Password

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(authData.Password))
	if err != nil {
		log.Println("Authentication failed for user:", authData.Username)
		return
	}

	// Authentication successful, continue with WebSocket connection
	log.Println("Authentication successful for user:", currentUser.Username)

	// Authentication successful, add connection to map
	mutex.Lock()
	_, ok := connections[currentUser.ID]
	if ok {
		delete(connections, currentUser.ID)
		fmt.Printf("Closing previous websocket connection for user %s \n", currentUser.Username)
	}
	connections[currentUser.ID] = conn
	mutex.Unlock()

	// Deliver pending pvt messages
	undeliveredPvtMsgs, err := models.GetUndeliveredPvtMsgsByUserID(currentUser.ID)
	if err != nil {
		log.Printf("Error fetching undelivered message for user %s: %s\n", currentUser.Username, err)
		return
	}
	mutex.Lock()
	for _, currPvtMsg := range undeliveredPvtMsgs {
		byteMsg, err := json.Marshal(currPvtMsg)
		if err != nil {
			log.Printf("Error writing message to WebSocket from user %s: %s\n", currentUser.Username, err)
			continue
		}
		err = conn.WriteMessage(1, byteMsg)
		if err != nil {
			log.Printf("Error writing message to WebSocket from user %s: %s\n", currentUser.Username, err)
			break
		}
		_, err = models.MarkPvtMsgDelivered(currPvtMsg.ID)
		if err != nil {
			log.Printf("Unable to mark message as delivered msgID: %d\n", currPvtMsg.ID)
			continue
		}
	}
	mutex.Unlock()

	// Deliver pending group messages
	// TODO write logic

	// Start a goroutine to handle messages from this WebSocket connection
	activeConnections.Add(1)
	go handleWebSocketMessages(currentUser, conn)

	activeConnections.Wait()
	// All goroutines have finished
	log.Println("All connections have finished. Exiting main logic.")
}

func handleWebSocketMessages(currentUser *models.User, conn *websocket.Conn) {
	defer activeConnections.Done()
	for {
		var msg models.Message
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from WebSocket from user %s: %s\n", currentUser.Username, err)
			break
		}
		// Decode JSON message
		err = json.Unmarshal(message, &msg)
		if err != nil {
			log.Printf("Error decoding JSON message from user %s: %s\n", currentUser.Username, err)
			continue
		}
		msg.SenderID = currentUser.ID
		// Process message (e.g., log or handle business logic)
		log.Printf("Received JSON message from user %s : %+v\n", currentUser.Username, msg)

		if msg.ReceiverID > 0 {
			// Pvt message
			mutex.Lock()
			receiverConn, ok := connections[msg.ReceiverID]
			if ok {
				// if receiver is connected send the message immediately
				err = receiverConn.WriteMessage(messageType, message)
				if err != nil {
					log.Printf("Error writing message to WebSocket from user %s: %s\n", currentUser.Username, err)
					break
				}
				msg.IsDelivered = true
			}
			mutex.Unlock()
		} else if msg.GroupID > 0 {
			// group messages

			// TODO insert message and get inserted id, then try sending the message to group user, and update last delivered message in db for the user

			// TODO send message to active group members
			// TODO store in db members who were not active
		}

		// Insert message into database along with the status whether it was delivered or not
		_, err = models.InsertNewMessage(msg)
		if err != nil {
			log.Printf("Failed to insert message from user %s: %+v : %s\n", currentUser.Username, msg, err)
			break
		}

	}

	// Connection closed or error occurred, remove connection from map
	mutex.Lock()
	delete(connections, currentUser.ID)
	mutex.Unlock()

	log.Printf("WebSocket connection closed for user %s\n", currentUser.Username)
}
