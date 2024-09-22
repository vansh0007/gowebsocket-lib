# WebSocket Server

A simple WebSocket server implemented in Go. This server allows clients to connect and communicate through WebSocket, sending and receiving messages in real-time. It also sends periodic messages to connected clients.

## Table of Contents

- [Getting Started](#getting-started)
- [Prerequisites](#prerequisites)
- [Running the Server](#running-the-server)
- [Testing with Postman](#testing-with-postman)
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

