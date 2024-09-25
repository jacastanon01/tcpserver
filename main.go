package main

import (
	"log/slog"
	"net"
	"os"
	"time"
)

// Instead of using log, I am opting for more granular logging by using depedency injection. This pattern also allows us to create new universal methods attached to our application that will recieve the app struct
type application struct {
	logger *slog.Logger
}

func (app *application) connect(conn net.Conn) {
	app.logger.Info("Inside the connect function")
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
	pool := NewPool(5)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	app := &application{
		logger: logger,
	}
	// First we estalbish the connection with our local server and listen to a specified port on the server
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		app.logger.Error(err.Error())
		return
	}
	// This is an infinite loop to keep the conneciton alive
	for {
		logger.Info("Waiting for client to connect")
		// This is awaiting the connection and establishing a client. This is a blocking call, so the server will not proceed until this receives data
		conn, err := listener.Accept()
		if err != nil {
			logger.Error(err.Error())
			continue // continue to next iteration if there is an error
		}
		app.logger.Info("Client connected!")

		// Process the request and configure response
		// by using the go keyword we are attaching that connection to a thread, allowing our server to handle multiple requests
		// go app.connect(conn)

		job := func() {
			// Set a dealine to handle a graceful exit if connection is open with no activity for 5 seconds
			err := conn.SetDeadline(time.Now().Add(5 * time.Second))
			// TODO Handle shutdown of connections and server
			if err != nil {
				os.Exit(1)
			}
			defer conn.Close() // Ensure the connection is closed after processing the request
			app.connect(conn)  // Process the connection (read request and send response)
			app.logger.Info("Job completed!")
		}
		// Add the job to the pool to be processed by an available worker
		pool.AddJob(job)
	}
}
