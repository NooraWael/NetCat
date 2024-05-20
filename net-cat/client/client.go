package client1

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type Client struct {
	Name string
	Conn net.Conn
}

// ConnectToServer creates a connection to the server at the specified address and port.
func ConnectToServer(address string, port string) (net.Conn, error) {
	serverAddr := fmt.Sprintf("%s:%s", address, port)
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Error connecting to the server:", err)
		return nil, err
	}
	return conn, nil
}

// InitializeClient asks the user for a username and sends it to the server.
func InitializeClient(conn net.Conn) string {
	reader := bufio.NewReader(os.Stdin)

	// Print the welcome message received from the server
	welcomeMessage, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading welcome message:", err)
		return ""
	}
	fmt.Print(welcomeMessage)

	// Read until the first newline to get the server's response
	serverResponse, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading server response:", err)
		return ""
	}
	fmt.Print(serverResponse)

	// Ask the user to enter their name
	fmt.Print("[ENTER YOUR NAME]: ")
	name, err := reader.ReadString('\n')
	if strings.TrimSpace(name) == "" {
		fmt.Println("Please enter a valid Name")
		return ""
	}
	if err != nil {
		fmt.Println("Error reading username:", err)
		return ""
	}

	// Send the name to the server
	_, err = conn.Write([]byte(name))
	if err != nil {
		fmt.Println("Error sending username to server:", err)
		return ""
	}

	return strings.TrimSpace(name)
}

// SendMessage sends a message to the server.
func SendMessage(conn net.Conn, message string) {
	_, err := conn.Write([]byte(message + "\n"))
	if err != nil {
		fmt.Println("Error sending message to server:", err)
	}
}

// ReceiveMessage receives a message from the server.
func ReceiveMessage(conn net.Conn) string {
	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error receiving message from server:", err)
		return ""
	}
	return message
}

// HandleInput handles user input for sending messages to the server.
func HandleInput(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	for {
		// Read user input
		fmt.Print("> ")
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading user input:", err)
			return
		}

		// Send the message to the server
		SendMessage(conn, message)
	}
}

// ListenToServer listens for incoming messages from the server and prints them.
func ListenToServer(conn net.Conn) {
	for {
		// Receive and print messages from the server
		message := ReceiveMessage(conn)
		if message != "" {
			fmt.Print(message)
		}
	}
}

// Add a method named Printinvalid to the Client type
func (c *Client) Printinvalid() {
	c.Conn.Write([]byte("Please Write a valid Username"))
}

// Add a method named Printinvalid to the Client type
func (c *Client) Printfull() {
	c.Conn.Write([]byte("Server is full :("))
}
