package utils

import (
	"log"
	"os"
	"time"
)

var (
	logFile *os.File
)

func InitLogger() {
	var err error
	logFile, err = os.OpenFile("backend.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	log.SetOutput(logFile)
}

func LogAction(userID, action, status, details string) {
	entry := time.Now().Format(time.RFC3339) + " | userID=" + userID + " | action=" + action + " | status=" + status + " | details=" + details
	log.Println(entry)
}

func CloseLogger() {
	if logFile != nil {
		logFile.Close()
	}
}
