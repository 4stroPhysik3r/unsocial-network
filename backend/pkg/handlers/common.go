package handlers

import (
	"backend/pkg/db"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Check the origin of your requests here
		// For development purposes, allow all origins
		return true
	},
}

func generateUniqueFileName(fileExtension string) string {
	currentTime := time.Now()
	randomPart := rand.Int63() // Generates a pseudo-random int64 number

	// Combine the current time (as Unix time in nanoseconds) with the random number
	fileName := fmt.Sprintf("%d-%d%s", currentTime.UnixNano(), randomPart, fileExtension)

	return fileName
}

func userIDFromSession(r *http.Request) int {
	sessionToken := extractSessionToken(r)

	if sessionToken == "" {
		return 0 // Return 0 if session token is not found
	}

	if !db.ValidateSessionToken(sessionToken) {
		return 0 // Return 0 if session token is invalid
	}

	userID, err := db.GetUserIDFromSessionToken(sessionToken)
	if err != nil {
		log.Printf("Error retrieving user ID from session token: %v", err)
		return 0 // Return 0 if there is an error retrieving user ID
	}

	return userID
}

func extractSessionToken(r *http.Request) string {
	const sessionToken = "session_token"
	cookie, err := r.Cookie(sessionToken)
	if err == nil {
		return cookie.Value
	}

	log.Printf("session token not found")
	return ""
}

// customFileServer creates a handler to serve static files from a given root.
// It prevents directory listings by not serving directory paths.
func CustomFileServer(root http.FileSystem) http.Handler {
	fs := http.FileServer(root)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Prevent directory listing by checking if the path is a directory
		path := r.URL.Path
		if f, err := root.Open(path); err == nil {
			defer f.Close()
			if stat, err := f.Stat(); err == nil && stat.IsDir() {
				http.NotFound(w, r)
				return
			}
		}

		fs.ServeHTTP(w, r)
	})
}

func saveBase64File(base64Data, directoryPath string) (string, error) {

	// Split the string to separate the metadata from the data itself
	dataParts := strings.Split(base64Data, ",")
	if len(dataParts) != 2 {
		return "", errors.New("invalid base64 data")
	}

	// Decode the file data from base64
	decodedData, err := base64.StdEncoding.DecodeString(dataParts[1])
	if err != nil {
		return "", err
	}

	// Extract MIME type and choose file extension
	mimeType := strings.Split(dataParts[0], ";")[0]
	mimeType = strings.TrimPrefix(mimeType, "data:")

	var fileExtension string
	switch mimeType {
	case "image/jpeg":
		fileExtension = ".jpg"
	case "image/png":
		fileExtension = ".png"
	case "image/gif":
		fileExtension = ".gif"
	default:
		return "", errors.New("unsupported file type")
	}

	// Generate a unique filename for the file to avoid conflicts
	fileName := generateUniqueFileName(fileExtension)

	// Write the data to a file
	filePath := directoryPath + "/" + fileName

	err = os.WriteFile(filePath, decodedData, 0666)
	if err != nil {
		return "", err
	}

	return fileName, nil
}
