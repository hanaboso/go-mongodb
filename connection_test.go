package mongodb

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"os/exec"
	"testing"
	"time"
)

const dsn = "MONGO_DSN"
const stop = "stop"
const start = "start"
const docker = "docker"
const container = "go-mongodb_mongodb_1"
const collection = "Collection"

var connection = Connection{}

func TestConnect(t *testing.T) {
	connection.Connect(os.Getenv(dsn))
	_ = connection.Database.Drop(context.Background())

	innerContext, cancel := connection.Context()
	defer cancel()
	document, err := connection.Database.Collection(collection).InsertOne(innerContext, map[string]interface{}{"key": "value"})

	assert.Nil(t, err)
	assert.NotNil(t, document)

	var result interface{}
	err = connection.Database.Collection(collection).FindOne(innerContext, primitive.D{{"key", "value"}}).Decode(&result)

	assert.Nil(t, err)
	assert.Equal(t, getInsertedID(document), getSelectedID(result))
}

func TestDisconnect(t *testing.T) {
	connection.Connect(os.Getenv(dsn))
	_ = connection.Database.Drop(context.Background())

	connection.Disconnect()

	assert.Nil(t, nil)
}

func TestIsConnected(t *testing.T) {
	if os.Getenv("GITLAB_CI") != "" {
		t.Skip()
	}

	connection.Connect(os.Getenv(dsn))
	_ = connection.Database.Drop(context.Background())
	assert.True(t, connection.IsConnected())

	_ = exec.Command(docker, stop, container).Run()
	assert.False(t, connection.IsConnected())

	_ = exec.Command(docker, start, container).Run()
	assert.True(t, connection.IsConnected())
}

func TestConnectError(t *testing.T) {
	_ = exec.Command(docker, stop, container).Run()

	go func() {
		for range time.After(time.Second) {
			_ = exec.Command(docker, start, container).Run()
		}
	}()

	connection.Connect(os.Getenv(dsn))
	_ = connection.Database.Drop(context.Background())

	innerContext, cancel := connection.Context()
	defer cancel()
	document, err := connection.Database.Collection(collection).InsertOne(innerContext, map[string]interface{}{"key": "value"})

	assert.Nil(t, err)
	assert.NotNil(t, document)

	var result interface{}
	err = connection.Database.Collection(collection).FindOne(innerContext, primitive.D{{"key", "value"}}).Decode(&result)

	assert.Nil(t, err)
	assert.Equal(t, getInsertedID(document), getSelectedID(result))
}

func getInsertedID(result *mongo.InsertOneResult) string {
	return result.InsertedID.(primitive.ObjectID).Hex()
}

func getSelectedID(result interface{}) string {
	return result.(primitive.D).Map()["_id"].(primitive.ObjectID).Hex()
}
