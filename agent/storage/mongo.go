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

var useUpsert = true

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
func (ms *MongoStore) Get(req *protov1.QueryRequest) (*did.Identifier, *did.ProofLD, error) {
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
		return nil, nil, errors.New("no information available")
	}

	// Decode result
	record := map[string]string{}
	if err := res.Decode(&record); err != nil {
		return nil, nil, errors.New("invalid record contents")
	}
	return decodeRecord(record)
}

// Save will create or update an entry for the provided DID instance.
func (ms *MongoStore) Save(id *did.Identifier, proof *did.ProofLD) error {
	data, err := json.Marshal(id.Document(true))
	if err != nil {
		return err
	}
	pp, err := json.Marshal(proof)
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
			{Key: "proof", Value: base64.RawStdEncoding.EncodeToString(pp)},
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

func decodeRecord(r map[string]string) (*did.Identifier, *did.ProofLD, error) {
	if _, ok := r["document"]; !ok {
		return nil, nil, errors.New("invalid record contents")
	}
	if _, ok := r["proof"]; !ok {
		return nil, nil, errors.New("invalid record contents")
	}
	d1, err := base64.RawStdEncoding.DecodeString(r["document"])
	if err != nil {
		return nil, nil, errors.New("invalid record contents")
	}
	d2, err := base64.RawStdEncoding.DecodeString(r["proof"])
	if err != nil {
		return nil, nil, errors.New("invalid record contents")
	}

	// Restore DID document
	doc := &did.Document{}
	if err = json.Unmarshal(d1, doc); err != nil {
		return nil, nil, errors.New("invalid record contents")
	}
	id, err := did.FromDocument(doc)
	if err != nil {
		return nil, nil, err
	}

	// Restore proof
	proof := &did.ProofLD{}
	if err = json.Unmarshal(d2, proof); err != nil {
		return nil, nil, errors.New("invalid record contents")
	}
	return id, proof, nil
}
