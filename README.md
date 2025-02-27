# AuthServer

AuthServer is a simple authentication server implemented in Go. It provides basic functionality for managing and authenticating users using either API keys or username/password combinations.

## Features

- Add new users with API keys or username/password
- Authenticate users using API keys or username/password
- In-memory storage using a map-based datastore
- Web interface implemented for interacting with AuthServer
- Remove user using userId or API key

## Usage

To run the AuthServer:

1. Ensure you have Go installed on your system
2. Clone this repository
3. Navigate to the project directory
4. Run the command: `go run main.go`

## Interacting with the Server

The server provides a web interface with the following options:

1. Add an API key
2. Add a user ID and password
3. Authenticate with an API key
4. Authenticate with a user ID and password
5. Remove API key
6. Remove a user ID

Follow the prompts to perform the desired action.

## Future Improvements

- Implement persistent storage
- Add password update functionality
- Implement API key reissuance
- Improve error handling and user feedback
- Add unit tests for better code coverage

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is open source and available under the [MIT License](LICENSE).
