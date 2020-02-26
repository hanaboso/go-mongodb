package mongodb

import (
	"context"
	"fmt"
	"log"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

const errorFormat = "[MongoDB] %+v"

// Connection represents MongoDB connection
type Connection struct {
	Database *mongo.Database
	timeout  time.Duration
}

// Connect creates MongoDB connection
func (connection *Connection) Connect(dsn string) {
	connectionString, err := connstring.Parse(dsn)

	if err != nil {
		panic(fmt.Sprintf(errorFormat, err))
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(dsn))

	if err != nil {
		log.Println(fmt.Sprintf(errorFormat, err))
		connection.Connect(dsn)

		return
	}

	if err := client.Connect(context.Background()); err != nil {
		log.Println(fmt.Sprintf(errorFormat, err))
		connection.Connect(dsn)

		return
	}

	if err := client.Ping(context.Background(), nil); err != nil {
		log.Println(fmt.Sprintf(errorFormat, err))
		connection.Connect(dsn)

		return
	}

	connection.Database = client.Database(connectionString.Database, nil)
	connection.timeout = getTimeout(connectionString)
}

// Disconnect from MongoDB
func (connection *Connection) Disconnect() {
	err := connection.Database.Client().Disconnect(context.Background())

	if err != nil {
		log.Println(fmt.Sprintf(errorFormat, err))
		connection.Disconnect()

		return
	}
}

// IsConnected checks connection status
func (connection *Connection) IsConnected() bool {
	return connection.Database.Client().Ping(context.Background(), nil) == nil
}

// Context creates context with timeout from MongoDB connection string
func (connection *Connection) Context() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), connection.timeout)
}

func getTimeout(connectionString connstring.ConnString) time.Duration {
	timeouts := []int{
		int(connectionString.ConnectTimeout.Milliseconds()),
		int(connectionString.SocketTimeout.Milliseconds()),
		int(connectionString.ServerSelectionTimeout.Milliseconds()),
	}

	sort.Ints(timeouts)

	return time.Duration(timeouts[len(timeouts)-1]) * time.Millisecond
}
