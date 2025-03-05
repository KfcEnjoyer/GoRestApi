# About 
A simple desktop application with a GUI that allows users to send HTTP requests to APIs, similar to Postman. This app supports various request methods and now it supports only JSON format but I will add other formats in the future, also allows you to manage and save requests for easy reuse.
This tool will act as Postman API without a need to log in 

# Features
- Send HTTP requests: GET, POST, PUT, DELETE, PATCH.
- Save Requests: Save your commonly used requests for easy access.
- Response Handling: View formatted responses, with JSON pretty-printing and error handling.
- Cross-Platform: Built as a desktop application accessible on all major platforms (Windows, Mac, Linux).
- Customizable Environment: Support for multiple API environments (development, staging, production).

# Instalation 

## Prerequisites
- Make sure you have Go installed on your machine.
- GCC (for Fyne dependencies on Windows/Linux)
- If you don't have Go installed, follow the instructions [here](https://golang.org/doc/install) to set up Go.
- Depending on OS follow [these](https://docs.fyne.io/started/) steps to ensure fyne is working properly

## Clone the Repository
```
git clone https://github.com/KfcEnjoyer/GoRestApi.git
cd http-request-sender
```

## Install Dependencies
```
go mod tidy
```

## Run the Application
After installing dependencies, you can build and run the app:
```
go run main.go

or

go build -o myapp
```

# Usage

1. Select HTTP Method: Choose from available HTTP methods (GET, POST, PUT, DELETE).
2. Enter API URL: Input the API endpoint URL.
4. Enter Request Body: Input data in JSON format.
5. Send the Request: Click on "Send" to make the API call and view the response.
6. Save Requests: Option to save requests for later use.

# Contributing
I welcome contributions! If you find a bug or want to add a new feature, feel free to create an issue or submit a pull request.

## How to Contribute
### Fork the repository.
- Create a new branch: git checkout -b feature/your-feature-name.
- Commit your changes: git commit -m 'Add new feature'.
- Push to the branch: git push origin feature/your-feature-name.
- Open a pull request.

  
