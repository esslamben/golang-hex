package database

import (
	"context"
	"fmt"
	"time"

	"github.com/esslamb/golang-hex/pkg/user"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	dbName = "golangHex"
	usersCollection = "users"
)

// Mongo is responsible for representing a mongo instance and all methods
// that it implements
type Mongo struct {
	Client *mongo.Client
}

// OpenClientConnection opens a new mongo connection using the given uri,
// in addition it also pings the new connection.
func OpenClientConnection(uri string) (*Mongo, error) {
	log.Infof("Connecting to mongo instance on uri: %s", uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		cancel()
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		cancel()
		return nil, err
	}

	cancel()

	log.Info("Mongo connection opened")

	return &Mongo{client}, nil
}

// InsertUser inserts new user to users collection
func (m *Mongo) InsertUser(u user.User) error {
	collection := m.Client.Database(dbName).Collection(usersCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := collection.InsertOne(ctx, u)
	if err != nil {
		cancel()
		return err
	}

	cancel()

	return nil
}

// FindUser finds a user based on a UUID
func (m *Mongo) FindUser(u uuid.UUID) (user.User, error) {
	var result user.User
	collection := m.Client.Database(dbName).Collection(usersCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	err := collection.FindOne(ctx, bson.M{"_id": u.Bytes()}).Decode(&result)
	if err != nil {
		cancel()
		return result, err
	}

	cancel()

	return result, nil
}

// UpdateUser finds a user based on a UUID and updates it
// according to the data passed.
func (m *Mongo) UpdateUser(u user.User, uuid uuid.UUID) error {
	filter := bson.M{"_id": uuid.Bytes()}
	update := bson.D{
		{"$set", bson.D{
			{"name", u.Name},
		}},
	}

	collection := m.Client.Database(dbName).Collection(usersCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		cancel()
		return err
	}
	fmt.Println(res)
	cancel()

	return nil
}

// DeleteUser deletes a user based on a UUID
func (m *Mongo) DeleteUser(u uuid.UUID) error {
	collection := m.Client.Database(dbName).Collection(usersCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := collection.DeleteOne(ctx, bson.M{"_id": u.Bytes()})
	if err != nil {
		cancel()
		return err
	}

	cancel()

	return nil
}
