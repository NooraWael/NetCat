package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"netcat/utilities"
	"strings"
	"sync"
)

type Client struct {
	Name string
	Conn net.Conn
	// Other relevant client information
}

type Server struct {
	listener net.Listener
}

var clientsMutex sync.Mutex
var connectedClients = make(map[*Client]bool)

var messageHistory []string

func AddClient(client *Client) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	connectedClients[client] = true
}

func RemoveClient(client *Client) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	delete(connectedClients, client)
}

// start the sever on the specified port
func (s *Server) StartServer(port string) {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		// Handle the error if there's an issue with creating the listener.
		fmt.Println("Error starting the server:", err)
		return
	}
	s.listener = listener
	fmt.Println("server is listening on address:" + port)
	messageHistory = []string{}
	s.AcceptConnections()
}

// new clients
func (s *Server) AcceptConnections() {

	for {

		newConn, err := s.listener.Accept()
		if err != nil {
			//handle error
			fmt.Println("Error with accepting Connection: ", err)
			continue
		}

		//to be able to handle more than one connection we have to do this
		// the go basically means new route
		go s.HandleClientConnections(newConn)

	}
}

// handles the communication with a connected client
func (s *Server) HandleClientConnections(conn net.Conn) {
	// Send a welcome message to the client
	SendWelcomeMessage(conn)
	defer conn.Close() // Close the connection when the function exits.

	// Read the client's name (assuming it's sent first)
	scanner := bufio.NewScanner(conn)
	if !scanner.Scan() {
		fmt.Println("Error reading client name:", scanner.Err())
		return
	}

	// Trim the newline character from the client's name
	clientName := strings.TrimSpace(scanner.Text())
	fmt.Println("Received client name:", clientName)

	// Create a new Client instance with the name and connection
	client := &Client{
		Name: clientName,
		Conn: conn,
	}
	if clientName == "" {
		fmt.Println("Empty client name received. Closing connection.")
		client.Printinvalid()

		return
	}
	print(len(connectedClients))
	if len(connectedClients) >= 10 {
		fmt.Println("Server is Full")
		return
	}
	AddClient(client)

	// Notify other clients that a new client has joined
	SendJoinMessage(client)
	fmt.Println(utilities.GetTime(), "Join message sent for:", clientName)

	// Send chat history to the new client
	SendChatHistory(client)
	fmt.Println(utilities.GetTime(), "Chat history sent to:", clientName)

	// Main loop to handle incoming messages from the client
	for scanner.Scan() {
		message := scanner.Text()
		message1 := strings.TrimSpace(message)

		// Check for an empty message
		if message1 != "" {
			// Get the current timestamp as a string

			// Format the message with the timestamp
			formattedMessage := fmt.Sprintf(" [%s]: %s", client.Name, message)

			// Broadcast the formatted message to all clients
			BroadcastMessage(formattedMessage, client)

			// Add the formatted message to the history
			messageHistory = append(messageHistory, formattedMessage)
		}
	}

	// Handle the case where the client disconnects
	// Notify other clients that this client has left
	SendLeaveMessage(client)
}

func BroadcastMessage(message string, sender *Client) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	// Create a copy of connected clients to avoid concurrent map iteration and modification
	clientsCopy := make(map[*Client]bool)
	for client := range connectedClients {
		clientsCopy[client] = true
	}

	// Iterate over connected clients
	for client := range clientsCopy {
		// Skip the sender
		//		if client == sender {
		//			continue
		//		}

		// Send the message to each client
		_, err := client.Conn.Write([]byte(utilities.GetTime() + message + "\n"))
		if err != nil {
			// Handle error (e.g., log it)
			fmt.Println("Error broadcasting message:", err)
		}
	}
}

func (c *Client) Printinvalid() {
	c.Conn.Write([]byte("Please Write a valid Username"))
}

func SendWelcomeMessage(conn net.Conn) {
	linuxLogo := "Welcome to TCP-Chat!\n" +
		"         _nnnn_\n" +
		"        dGGGGMMb\n" +
		"       @p~qp~~qMb\n" +
		"       M|@||@) M|\n" +
		"       @,----.JM|\n" +
		"      JS^\\__/  qKL\n" +
		"     dZP        qKRb\n" +
		"    dZP          qKKb\n" +
		"   fZP            SMMb\n" +
		"   HZM            MMMM\n" +
		"   FqM            MMMM\n" +
		" __| \".        |\\dS\"qML\n" +
		" |    `.       | `' \\Zq\n" +
		"_)      \\.___.,|     .'\n" +
		"\\____   )MMMMMP|   .'\n" +
		"     `-'       `--'\n" +
		"[ENTER YOUR NAME]: "

	_, err := conn.Write([]byte(linuxLogo))
	if err != nil {
		// Handle error (e.g., log it)
		fmt.Println("Error sending welcome message:", err)
	}
}

// SendJoinMessage sends a message to all clients that a new client has joined.
func SendJoinMessage(client *Client) {
	fmt.Println("Sending join message")
	message := fmt.Sprintf("[SYSTEM]: %s has joined the chat...\n", client.Name)
	BroadcastMessage(message, nil) // Broadcast the join message to all clients
}

// sends a message that the user left
func SendLeaveMessage(client *Client) {
	message := fmt.Sprintf("[SYSTEM]: %s has left the chat...\n", client.Name)
	BroadcastMessage(message, nil) // Broadcast the leave message to all clients

	// Remove the client from the connected clients
	RemoveClient(client)
}

// SendChatHistory sends chat history to a newly joined client.
func SendChatHistory(client *Client) {
	// Lock the mutex to safely access the history
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	// Send each message from the history to the client
	for _, message := range messageHistory {
		_, err := client.Conn.Write([]byte(message + "\n"))
		if err != nil {
			// Handle error (e.g., log it)
			fmt.Println("Error sending chat history:", err)
		}
	}
}

// Close the server listener
func (s *Server) CloseServer() {
	if s.listener != nil {
		err := s.listener.Close()
		if err != nil {
			fmt.Println("Error closing the server:", err)
		} else {
			fmt.Println("Server closed successfully.")
		}
	}
}

func GetIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	localaddr := conn.LocalAddr().(*net.UDPAddr)
	return localaddr.IP.String()
}
