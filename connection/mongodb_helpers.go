package connection

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Settings struct {
	MongoDBUser     string
	MongoDBPassword string
}

// ConnectDB: Return a data base model.
func ConnectDB(dbName string, settings Settings) *mongo.Database {

	// MongoDBClusterUri
	uri := fmt.Sprintf("mongodb+srv://%s:%s@roommatepaymentcluster.baqgzpz.mongodb.net/?retryWrites=true&w=majority",
		settings.MongoDBUser, settings.MongoDBPassword)

	// Client options
	clientOptions := options.Client().ApplyURI(uri)

	// Context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!!!!")

	// Connect with DB

	db := client.Database(dbName)

	return db
}

type ErrorResponse struct {
	StatusCode int    `json:"status"`
	Message    string `json:"message"`
}

// GetError : Helper function to get the error as an ErrorResponse object.
func GetError(err error, w http.ResponseWriter) {
	log.Fatal(err.Error())

	var response = ErrorResponse{
		Message:    err.Error(),
		StatusCode: http.StatusInternalServerError,
	}

	message, _ := json.Marshal(response)

	w.WriteHeader(response.StatusCode)
	w.Write(message)
}
