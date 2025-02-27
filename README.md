AuthServer
AuthServer is a simple authentication server implemented in Go. It provides basic functionality for managing and authenticating users using either API keys or username/password combinations.

Features
Add new users with API keys or username/password
Authenticate users using API keys or username/password
In-memory storage using a map-based datastore
Web interface for interaction with server
Usage
To run the AuthServer:

Ensure you have Go installed on your system
Clone this repository
Navigate to the project directory
Run the command: go run main.go
Interacting with the Server
The server provides a command-line interface with the following options:

Add an API key
Add a user ID and password
Authenticate with an API key
Authenticate with a user ID and password
Follow the prompts to perform the desired action.

Future Improvements
Implement persistent storage
Add password update functionality
Implement API key reissuance
Improve error handling and user feedback
Add unit tests for better code coverage
Contributing
Contributions are welcome! Please feel free to submit a Pull Request.

License
This project is open source and available under the MIT License.
