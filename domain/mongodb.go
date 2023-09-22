package domain

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStore struct {
	Vehicles IVehicle
}

func ConnectToDatabase(uri string, dbname string) *MongoStore {
	if uri == "" {
		log.Fatal("[EROR] Not found 'MONGODB_URI' environment variable.")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("[EROR] Couldn't connect to mongo db.")
	}

	db := client.Database(dbname)

	return &MongoStore{
		Vehicles: InitVehicleStore(db.Collection("vehicles")),
	}
}
