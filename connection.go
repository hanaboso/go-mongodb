package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Connection struct {
	Database *mongo.Database
	timeout  *time.Duration
}

func (connection *Connection) Connect(clientOptions *options.ClientOptions, database string) {
	client, err := mongo.NewClient(clientOptions)

	if err != nil {

	}

	err = client.Connect(context.Background())

	if err != nil {

	}

	connection.Database = client.Database(database, nil)

	if clientOptions.ConnectTimeout != nil {
		connection.timeout = clientOptions.ConnectTimeout
	}

	if clientOptions.SocketTimeout != nil {
		connection.timeout = clientOptions.SocketTimeout
	}

}

func (connection *Connection) Context() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), *connection.timeout)
}

func (connection *Connection) Disconnect() {
	_ = connection.Database.Client().Disconnect(context.Background())
}
