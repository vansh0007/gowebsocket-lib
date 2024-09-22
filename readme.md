# WebSocket Server

A simple WebSocket server implemented in Go. This server allows clients to connect and communicate through WebSocket, sending and receiving messages in real-time. It also sends periodic messages to connected clients.

## Table of Contents

- [Getting Started](#getting-started)
- [Prerequisites](#prerequisites)
- [Running the Server](#running-the-server)
- [Testing with Postman](#testing-with-postman)
- [How It Works](#how-it-works)
- [License](#license)

## Getting Started

These instructions will help you set up and run the WebSocket server on your local machine.

### Prerequisites

- [Go](https://golang.org/dl/) (version 1.14 or later)
- A terminal or command prompt
- Postman (optional, for testing)

### Running the Server

1. **Clone the Repository**:

   1. Open your terminal and run the following command to clone the repository:

   ```bash
   git clone https://github.com/username/your-repo.git
   ```
Replace username and your-repo with your GitHub username and repository name.

2. Navigate to the Project Directory:
 ```
cd your-repo
```

3. Install Dependencies:

If you haven't already initialized your Go module, do so:

```
go mod init gowebscoket-lib
```
Then, install any dependencies (if applicable):

```
go mod tidy
```

4. Run the Server:

Execute the following command to start the server:
````
go run main.go
````

You should see a message indicating that the WebSocket server has started.

5. Testing with Postman
   1. Open Postman.
   2. Create a new WebSocket request.
   3. Enter the WebSocket URL: ws://localhost:8080/ws.
   4. Click "Connect".
   5. You can send messages through Postman, and the server will echo them back to you. Additionally, you will receive periodic messages from the server every 5 seconds.

How It Works
Overview of WebSocket Server
WebSocket Protocol:
WebSockets provide a full-duplex communication channel over a single, long-lived connection. This allows for real-time data exchange between a client (like a web browser or Postman) and a server.
Key Components
HTTP Server Setup:

The server listens for incoming HTTP requests on a specified port (in this case, 8080) using the http.ListenAndServe function.
WebSocket Upgrade:

When a client sends a request to the /ws endpoint, the server upgrades this HTTP connection to a WebSocket connection using the websocket.Upgrade function. This step is crucial as it switches from HTTP to the WebSocket protocol.
Hijacking the Connection:

The server uses http.Hijacker to hijack the connection, allowing it to take control over the raw network connection (a net.Conn). This enables the server to read from and write to the WebSocket connection directly.
Creating a WebSocket Connection:

A new WebSocket connection (wsConn) is created using the hijacked net.Conn. This object provides methods for reading and writing messages in a WebSocket-compatible format.
Message Handling
Listening for Incoming Messages:

In a loop, the server listens for messages sent by the client using the wsConn.ReadMessage() method. When a message is received, it logs the message to the console.
Echoing Messages:

After receiving a message, the server immediately sends it back to the client using wsConn.WriteMessage(). This echo functionality demonstrates bidirectional communication.
Sending Periodic Messages
Ticker for Periodic Messages:

A time.Ticker is set up to send messages at specified intervals (e.g., every 5 seconds). This runs in a separate goroutine, which means it operates concurrently with the message listening loop.
Sending Messages:

Inside the ticker’s loop, the server sends a predefined message to the connected client. If an error occurs while sending the message, it logs the error.
Client Interaction
Connecting from Clients:

Clients (like browsers or Postman) can connect to the WebSocket server by targeting the ws://localhost:8080/ws endpoint. Once connected, they can send messages and receive responses, including the periodic updates from the server.
Real-time Communication:

This setup allows for real-time communication where both the client and server can send messages at any time without needing to establish a new connection.
