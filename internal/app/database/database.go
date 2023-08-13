package database

import (
	"context"
	"log"
	"time"

	"github.com/tarkue/tolpi-backend/config"
	vk "github.com/tarkue/tolpi-backend/internal/app/VK"
	"github.com/tarkue/tolpi-backend/internal/app/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	client *mongo.Client
}

func New() *DB {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	credential := options.Credential{
		AuthMechanism: "SCRAM-SHA-256",
		Username:      config.DataBaseUserName,
		Password:      config.DataBasePassword,
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.DataBaseUri).SetAuth(credential))
	if err != nil {
		log.Fatal(err)
	}
	return &DB{
		client: client,
	}
}

func (db *DB) CreateUser(input *model.NewUser, userId string) *model.User {
	collection := db.client.Database("Tolpi").Collection("Users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var findResult *model.User
	res := collection.FindOne(ctx, bson.M{"userid": userId})

	if res.Err() == nil {
		res.Decode(&findResult)
		if findResult.Avatar != input.Avatar ||
			findResult.FirstName != input.FirstName ||
			findResult.LastName != input.LastName {
			_, err := collection.UpdateOne(
				ctx, bson.M{"userid": userId},
				bson.M{
					"$set": bson.M{
						"avatar":    input.Avatar,
						"firstname": input.FirstName,
						"lastname":  input.LastName,
					},
				},
			)
			if err != nil {
				log.Fatal(err)
			}
		}

		res = collection.FindOne(ctx, bson.M{"userid": userId})
		res.Decode(&findResult)
		findResult.Status = vk.GetUserStatus(userId)
		findResult.Tolpies = db.GetUserTolpiesList(userId)

		return findResult

	}

	user := &model.User{
		Avatar:      input.Avatar,
		UserID:      userId,
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		TrackerList: []string{},
		Status:      vk.GetUserStatus(userId),
		Tolpies:     db.GetUserTolpiesList(userId),
	}

	_, err := collection.InsertOne(ctx, user)

	if err != nil {
		log.Fatal(err)
	}

	return user
}

func (db *DB) CreateTolpi(input *model.NewTolpi, userId string) *model.Tolpi {
	collection := db.client.Database("Tolpi").Collection("Tolpies")
	collectionUsers := db.client.Database("Tolpi").Collection("Users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	findRes := collectionUsers.FindOne(ctx, bson.M{"userid": userId})

	user := model.User{}
	findRes.Decode(&user)

	tolpi := &model.Tolpi{
		Text:      input.Text,
		Timestamp: int(time.Now().Unix()),
		User:      &user,
		Country:   input.Country,
	}
	_, err := collection.InsertOne(ctx, tolpi)

	if err != nil {
		log.Fatal(err)
	}
	return tolpi
}

func (db *DB) UpdateUserCountry(userID string, country string) *model.User {
	collection := db.client.Database("Tolpi").Collection("Users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.UpdateOne(ctx, bson.M{"userid": userID}, bson.M{"$set": bson.M{"country": country}})
	if err != nil {
		log.Fatal(err)
	}

	findRes := collection.FindOne(ctx, bson.M{"userid": userID})

	user := &model.User{}
	findRes.Decode(&user)

	return user
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
	user.Status = vk.GetUserStatus(userID)
	user.Tolpies = db.GetUserTolpiesList(userID)

	return &user
}

func (db *DB) GetLastTolpies(country string) []*model.Tolpi {
	collection := db.client.Database("Tolpi").Collection("Tolpies")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	opts := options.Find().SetLimit(80).SetSort(bson.M{"$natural": -1})

	var tolpies []*model.Tolpi
	res, err := collection.Find(ctx, bson.M{"country": country}, opts)

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

func (db *DB) GetSubscribes(userID string) []*model.User {
	collection := db.client.Database("Tolpi").Collection("Users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := collection.Find(ctx, bson.M{"trackerlist": userID})
	if err != nil {
		log.Fatal(err)
		return nil
	}

	var users []*model.User
	for res.Next(ctx) {
		var user *model.User

		err = res.Decode(&user)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}

	return users
}

func (db *DB) GetUserTolpiesList(userID string) []*model.Tolpi {
	collection := db.client.Database("Tolpi").Collection("Tolpies")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := collection.Find(ctx, bson.M{"user.userid": userID})
	if err != nil {
		log.Fatal(err)
		return nil
	}

	var tolpies []*model.Tolpi
	for res.Next(ctx) {
		var tolpi *model.Tolpi

		err = res.Decode(&tolpi)
		if err != nil {
			log.Fatal(err)
		}
		tolpies = append(tolpies, tolpi)
	}

	log.Print(tolpies)

	return tolpies
}
