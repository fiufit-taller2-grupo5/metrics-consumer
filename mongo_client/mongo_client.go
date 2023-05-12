package mongo_client

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	client *mongo.Client
}

func NewMongoClient() (*MongoClient, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(buildMongoURI()).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}

	return &MongoClient{
		client: client,
	}, nil
}

func (mongoClient *MongoClient) Disconnect() {
	err := mongoClient.client.Disconnect(context.TODO())
	if err != nil {
		println("Failed disconnecting Mongo client: " + err.Error())
	}
}

func (mongoClient *MongoClient) InsertJSONDocument(jsonStr *string, collectionName string) error {
	var jsonObj map[string]interface{}
	err := json.Unmarshal([]byte(*jsonStr), &jsonObj)
	if err != nil {
		return fmt.Errorf("failed to parse JSON: %v", err)
	}

	collection := mongoClient.client.Database("fiufit").Collection(collectionName)
	_, err = collection.InsertOne(context.Background(), jsonObj)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func buildMongoURI() string {
	const MongoUriTemplate = "mongodb+srv://%s:%s@fiufit.zdkdc6u.mongodb.net/?retryWrites=true&w=majority"
	username := "fiufitmetricscron"
	password := "Q7Re0TXRSyJsfUfY"

	return fmt.Sprintf(MongoUriTemplate, username, password)
}
