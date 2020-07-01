package storage

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"

	protov1 "github.com/bryk-io/did-method/proto/v1"
	"github.com/pkg/errors"
	"go.bryk.io/x/ccg/did"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var useUpsert bool = true

// Collection used to store the DID documents managed.
const didCol = "identifiers"

// MongoStore provides a storage handler utilizing MongoDB as underlying
// database. The connection strings must be of the form "mongodb://...";
// for example: "mongodb://localhost:27017"
type MongoStore struct {
	db *mongo.Database
}

// Open establish the connection and database selection for the instance.
// Must be called before any further operations.
func (ms *MongoStore) Open(info string) error {
	ctx := context.TODO()
	cl, err := mongo.Connect(ctx, options.Client().ApplyURI(info))
	if err != nil {
		return err
	}
	ms.db = cl.Database("didctl")
	return nil
}

// Close the client connection with the backend server.
func (ms *MongoStore) Close() error {
	return ms.db.Client().Disconnect(context.TODO())
}

// Exists returns true if the provided DID instance is already available
// in the store.
func (ms *MongoStore) Exists(id *did.Identifier) bool {
	query := bson.M{
		"method":  id.Method(),
		"subject": id.Subject(),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	res := ms.db.Collection(didCol).FindOne(ctx, query)
	return res.Err() != mongo.ErrNoDocuments
}

// Get a previously stored DID instance.
func (ms *MongoStore) Get(req *protov1.QueryRequest) (*did.Identifier, error) {
	// Run query
	query := bson.M{
		"method":  req.Method,
		"subject": req.Subject,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	res := ms.db.Collection(didCol).FindOne(ctx, query)

	// Check for result
	if res.Err() == mongo.ErrNoDocuments {
		return nil, errors.New("no information available")
	}

	// Decode result
	record := map[string]string{}
	if err := res.Decode(&record); err != nil {
		return nil, errors.New("invalid record contents")
	}
	if _, ok := record["document"]; !ok {
		return nil, errors.New("invalid record contents")
	}
	data, err := base64.RawStdEncoding.DecodeString(record["document"])
	if err != nil {
		return nil, errors.New("invalid record contents")
	}

	// Restore DID document
	doc := &did.Document{}
	if err = json.Unmarshal(data, doc); err != nil {
		return nil, errors.New("invalid record contents")
	}
	return did.FromDocument(doc)
}

// Save will create or update an entry for the provided DID instance.
func (ms *MongoStore) Save(id *did.Identifier) error {
	data, err := json.Marshal(id.Document(true))
	if err != nil {
		return err
	}
	filter := bson.M{
		"method":  id.Method(),
		"subject": id.Subject(),
	}
	record := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "method", Value: id.Method()},
			{Key: "subject", Value: id.Subject()},
			{Key: "document", Value: base64.RawStdEncoding.EncodeToString(data)},
		}},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err = ms.db.Collection(didCol).UpdateOne(ctx, filter, record, &options.UpdateOptions{Upsert: &useUpsert})
	return err
}

// Delete any existing record for the provided DID instance.
func (ms *MongoStore) Delete(id *did.Identifier) error {
	query := bson.M{
		"method":  id.Method(),
		"subject": id.Subject(),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := ms.db.Collection(didCol).DeleteOne(ctx, query)
	return err
}

// Description returns a brief summary for the storage instance.
func (ms *MongoStore) Description() string {
	return "MongoDB data store"
}
