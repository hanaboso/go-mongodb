package mongodb

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

func TestConnection_Connect(t *testing.T) {
	options := options.Client().ApplyURI("mongodb://127.0.0.25:27017/example?connectTimeoutMS=5000&heartbeatFrequencyMS=5000")
	connection := Connection{}
	connection.Connect(options, "example")

	context, cancel := connection.Context()
	defer cancel()

	doEvery(1*time.Second, func() {
		a, b := connection.Database.Collection("test").InsertOne(context, map[string]interface{}{"a": "b"})
		fmt.Println(a)
		fmt.Println(b)
	})

}

func doEvery(d time.Duration, f func()) {
	for _ = range time.Tick(d) {
		f()
	}
}
