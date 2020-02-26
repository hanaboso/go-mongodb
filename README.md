# Hanaboso GO MongoDB

**Download**
```
go mod download github.com/hanaboso/go-mongodb
```

**Usage**
```
import "github.com/hanaboso/go-mongodb"

mongodb := &mongodb.Connection{}
mongodb.Connect("mongodb://mongodb/database?connectTimeoutMS=2500&heartbeatFrequencyMS=2500")

context, cancel := mongodb.Context()
defer cancel()

mongodb.Database.Collection("collection").CountDocuments(context, primitive.D{{}})
mongodb.Disconnect()
```
