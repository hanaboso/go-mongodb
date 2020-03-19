package mongodb

import (
	"context"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const stop = "stop"
const start = "start"
const docker = "docker"
const sudo = "su-exec"
const user = "root"
const container = "go-mongodb_mongodb_1"
const collection = "Collection"

var connection = Connection{}

func TestConnect(t *testing.T) {
	connection.Connect(getDsn())
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
	connection.Connect(getDsn())
	_ = connection.Database.Drop(context.Background())

	connection.Disconnect()

	assert.Nil(t, nil)
}

func TestIsConnected(t *testing.T) {
	connection.Connect(getDsn())
	_ = connection.Database.Drop(context.Background())
	assert.True(t, connection.IsConnected())

	_ = getCmdContext(stop).Run()
	assert.False(t, connection.IsConnected())

	_ = getCmdContext(start).Run()
	assert.True(t, connection.IsConnected())
}

func TestConnectError(t *testing.T) {
	_ = getCmdContext(stop).Run()

	go func() {
		for range time.After(time.Second) {
			_ = getCmdContext(start).Run()
		}
	}()

	connection.Connect(getDsn())
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

func getDsn() string {
	if dsn := os.Getenv("MONGO_DSN"); dsn != "" {
		return dsn
	}

	return "mongodb://127.0.0.25/database?connectTimeoutMS=2500&serverSelectionTimeoutMS=2500&socketTimeoutMS=2500&heartbeatFrequencyMS=2500"
}

func getCmdContext(action string) *exec.Cmd {
	return exec.Command(sudo, user, docker, action, container)
}
