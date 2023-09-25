package domain

import (
	"context"
	"fmt"
	"mtrack/device_server/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ITrip interface {
	CreateTrip(imei string, tm int64, data types.TripEvent) (string, error)
	AddEvent(id string, tm int64, data types.TripEvent) error
}

type TripStore struct {
	col *mongo.Collection
}

func InitTripStore(col *mongo.Collection) ITrip {
	return &TripStore{
		col: col,
	}
}

func (s *TripStore) CreateTrip(imei string, tm int64, data types.TripEvent) (string, error) {
	obj := types.Trip{
		ID:         primitive.NewObjectID(),
		IMEI:       imei,
		CreatedAt:  tm,
		LastUpdate: tm,
		Events:     append([]types.TripEvent{}, data),
	}

	_, err := s.col.InsertOne(context.Background(), obj)

	return obj.ID.Hex(), err
}

func (s *TripStore) AddEvent(id string, tm int64, data types.TripEvent) error {
	var result types.Trip
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	err = s.col.FindOne(context.Background(), bson.D{{Key: "_id", Value: objectId}}).Decode(&result)
	if err != nil {
		return err
	}
	arr := append(result.Events, data)
	fmt.Println("FOUND >> PASS TO UPDATE >> ", len(arr))
	filter := bson.D{{Key: "_id", Value: objectId}}
	update := bson.D{
		{Key: "lastupdate", Value: tm},
		{Key: "events", Value: arr},
	}
	_, err = s.col.UpdateOne(context.Background(), filter, bson.D{{Key: "$set", Value: update}})
	return err
}
