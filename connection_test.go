package mongodb

import (
	"context"
	"log"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const stop = "stop"
const start = "start"
const docker = "docker"
const sudo = "su-exec"
const user = "root"
const container = "go-mongodb-mongodb"
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
	err = connection.Database.Collection(collection).FindOne(innerContext, bson.D{{"key", "value"}}).Decode(&result)

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

	err := getCmdContext(stop).Run()
	log.Printf("Command finished with error: %v", err)
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
	err = connection.Database.Collection(collection).FindOne(innerContext, bson.D{{"key", "value"}}).Decode(&result)

	assert.Nil(t, err)
	assert.Equal(t, getInsertedID(document), getSelectedID(result))
}

func getInsertedID(result *mongo.InsertOneResult) string {
	return result.InsertedID.(bson.ObjectID).Hex()
}

func getSelectedID(result interface{}) string {

	type Result struct {
		ID bson.ObjectID `bson:"_id"`
	}

	var a Result

	val, _ := bson.Marshal(result)
	_ = bson.Unmarshal(val, &a)

	return a.ID.Hex()
}

func getDsn() string {
	if dsn := os.Getenv("MONGO_DSN"); dsn != "" {
		return dsn
	}

	return "mongodb://127.0.0.1/database?connectTimeoutMS=2500&serverSelectionTimeoutMS=2500&socketTimeoutMS=2500&heartbeatFrequencyMS=2500"
}

func getCmdContext(action string) *exec.Cmd {
	return exec.Command(sudo, user, docker, action, container)
}
