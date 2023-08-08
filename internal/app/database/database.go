package database

import (
	"context"
	"log"
	"time"

	"github.com/tarkue/tolpi-backend/config"
	"github.com/tarkue/tolpi-backend/internal/app/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	client *mongo.Client
}

func New() *DB {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.DataBaseUri))
	if err != nil {
		log.Fatal(err)
	}
	return &DB{
		client: client,
	}
}

func (db *DB) CreateUser(input *model.NewUser) *model.User {
	collection := db.client.Database("Tolpi").Collection("Users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	return &model.User{
		ID:          res.InsertedID.(primitive.ObjectID).Hex(),
		UserID:      input.UserID,
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		TrackerList: []string{},
	}
}

func (db *DB) CreateTolpi(input *model.NewTolpi) *model.Tolpi {
	collection := db.client.Database("Tolpi").Collection("Tolpies")
	collectionUsers := db.client.Database("Tolpi").Collection("Users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	findRes := collectionUsers.FindOne(ctx, bson.M{"userid": input.UserID})

	user := model.User{}
	findRes.Decode(&user)

	tolpi := &model.Tolpi{
		Text:      input.Text,
		Timestamp: int(time.Now().Unix()),
		User:      &user,
	}
	_, err := collection.InsertOne(ctx, tolpi)

	if err != nil {
		log.Fatal(err)
	}
	return tolpi
}

func (db *DB) UpdateUserTrackers(userID string, trackers []string) {
	collection := db.client.Database("Tolpi").Collection("Users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.UpdateOne(ctx, bson.M{"userid": userID}, bson.M{"$set": bson.M{"trackerlist": trackers}})
	if err != nil {
		log.Fatal(err)
	}
}
func (db *DB) FindUserById(userID string) *model.User {
	collection := db.client.Database("Tolpi").Collection("Users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res := collection.FindOne(ctx, bson.M{"userid": userID})

	user := model.User{}
	res.Decode(&user)

	return &user
}

func (db *DB) GetLastTolpies() []*model.Tolpi {
	collection := db.client.Database("Tolpi").Collection("Tolpies")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	opts := options.Find().SetLimit(80).SetSort(bson.M{"$natural": -1})

	var tolpies []*model.Tolpi
	res, err := collection.Find(ctx, bson.M{}, opts)

	if err != nil {
		log.Fatal(err)
	}
	for res.Next(ctx) {
		var tolpi *model.Tolpi

		err = res.Decode(&tolpi)
		if err != nil {
			log.Fatal(err)
		}
		tolpies = append(tolpies, tolpi)
	}

	return tolpies
}
