package utills

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

var LogFilePath = "logs/app.log"

type RequestLogger struct {
	Method     string
	Url        string
	Message    string
	StatusCode int
	TimeStamp  string
}

type ErrorLogger struct {
	Error     string
	TimeStamp string
}

func NewReqErr(method, url, message string, statusCode int) *RequestLogger {
	return &RequestLogger{
		method,
		url,
		message,
		statusCode,
		time.Now().Format(time.RFC3339),
	}
}

func NewErrLog(err string) *ErrorLogger {
	return &ErrorLogger{
		err,
		time.Now().Format(time.RFC3339),
	}
}

func EnsureFile() {
	logDir := filepath.Dir(LogFilePath)
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
			log.Println(err)
		}
		fmt.Println("Path created")
	}

	if _, err := os.Stat(LogFilePath); os.IsNotExist(err) {
		if err := os.WriteFile(LogFilePath, []byte("[]"), 0644); err != nil {
			log.Println(err)
		}
	}

	file, err := os.OpenFile(LogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	log.SetOutput(file)

}

func CreateLog(logger ErrorLogger) {
	EnsureFile()

	logJson, err := json.Marshal(logger)
	if err != nil {
		log.Printf("Error serializing error log: %v", err)
		return
	}

	log.Println(string(logJson))
}

func CreateReqLog(req RequestLogger) {
	EnsureFile()

	reqJson, err := json.Marshal(req)
	if err != nil {
		log.Printf("Error serializing error log: %v", err)
		return
	}

	log.Println(string(reqJson))
}
