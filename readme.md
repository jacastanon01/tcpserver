# TCP Server in Go

This application will demonstrate how to set up a TCP connection using Go. 

## What is TCP?

Transmission Control Protocol (TCP) is built on top of the Internet Protocol Suite (TCP/IP) and establishes a connection between two computers. A request is sent from the client to the server and the response is sent back. A TCP connection guarantees delivery because data packets are transmitted in segments, and the protocol ensures that packets are received in the correct order. If any packet is lost or out of order, TCP handles retransmission, maintaining reliable data transfer. This connection is handled through a three-step "handshake" process: SYN, SYN-ACK, and ACK, which synchronizes both ends. Once the connection is established, data is transmitted in segments, and TCP ensures that all segments arrive intact and in order.

```go
// Instead of using log, I am opting for more granular logging by using depedency injection. This pattern also allows us to create new universal methods attached to our application that will recieve the app struct
type application struct {
	logger *slog.Logger
}

func (app *application) connect(conn net.Conn) {
	var buffer []byte = make([]byte, 1024)
	// Now that the connection is established, we can read the request
	_, err := conn.Read(buffer)
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}

	app.logger.Info("Processing the request")
	time.Sleep(5 * time.Second)
	// Since we will test this with cURL, the response should be formatted as an HTTP response
	res := "HTTP/1.1 200 OK\r\n\r\nHello, World!\r\n"

	conn.Write([]byte(res))
	// After sending the response, we can close the TCP connection
	conn.Close()
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	app := &application{
		logger: logger,
	}
	// First we estalbish the connection with our local server and listen to a specified port on the server
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		app.logger.Error(err.Error())
	}
	// This is an infinite loop to keep the conneciton alive
	for {
		logger.Info("Waiting for client to connect")
		// This is awaiting the connection and establishing a client. This is a blocking call, so the server will not proceed until this receives data
		conn, err := listener.Accept()
		if err != nil {
			logger.Error(err.Error())
		}
		app.logger.Info("Client connected!")
		// Process the request and configure response
		// by using the go keyword we are attaching that connection to a thread, allowing our server to handle multiple requests
        // This also decouples each TCP connection to handle its own read, write processes
		go app.connect(conn)
	}
}
```

Currently, there are optimization issues if multiple clients try to open a thread at once, but this is a general overview of how computers communicate over the internet using TCP.