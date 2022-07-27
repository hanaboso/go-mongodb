package mongodb

import (
	"context"
	"sort"
	"time"

	"github.com/hanaboso/go-log/pkg/zap"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"

	log "github.com/hanaboso/go-log/pkg"
)

// Connection represents MongoDB connection
type Connection struct {
	Database *mongo.Database
	timeout  time.Duration
	Log      log.Logger
}

// Connect creates MongoDB connection
func (connection *Connection) Connect(dsn string) {
	if connection.Log == nil {
		connection.Log = zap.NewLogger()
	}

	connectionString, err := connstring.Parse(dsn)

	if err != nil {
		connection.logContext().Error(err)
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(dsn))
	if err != nil {
		connection.Log.Error(err)
		connection.reconnect(dsn, err)

		return
	}

	if err := client.Connect(context.Background()); err != nil {
		_ = client.Disconnect(context.Background())
		connection.Log.Error(err)
		connection.reconnect(dsn, err)

		return
	}

	if err := client.Ping(context.Background(), nil); err != nil {
		_ = client.Disconnect(context.Background())
		connection.Log.Error(err)
		connection.reconnect(dsn, err)

		return
	}

	connection.Database = client.Database(connectionString.Database, nil)
	connection.timeout = getTimeout(connectionString)
}

// Disconnect from MongoDB
func (connection *Connection) Disconnect() {
	err := connection.Database.Client().Disconnect(context.Background())

	if err != nil {
		connection.logContext().Error(err)
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

// StartSession creates session
func (connection *Connection) StartSession(options ...*options.SessionOptions) (mongo.Session, error) {
	return connection.Database.Client().StartSession(options...)
}

func getTimeout(connectionString connstring.ConnString) time.Duration {
	timeouts := []int{
		int(connectionString.ConnectTimeout.Milliseconds()),
		int(connectionString.SocketTimeout.Milliseconds()),
		int(connectionString.ServerSelectionTimeout.Milliseconds()),
	}

	sort.Ints(timeouts)

	if timeout := time.Duration(timeouts[len(timeouts)-1]) * time.Millisecond; timeout > 0 {
		return timeout
	}

	return 25 * time.Second
}

func (connection *Connection) logContext() log.Logger {
	return connection.Log.WithFields(map[string]interface{}{
		"package": "MongoDB",
	})
}

func (connection *Connection) reconnect(dsn string, err error) {
	connection.logContext().Error(err)
	time.Sleep(createBackOff(5*time.Second, 60*time.Second, 5, false).Duration())
	connection.Connect(dsn)
}
