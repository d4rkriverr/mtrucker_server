package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type IVehicle interface {
	UpdateDeviceStatus(imei string, status bool) error
	UpdateCurrentStatus(imei string, eng bool, gsm, volt int64, acc, gps bool) error
}

type VehicleStore struct {
	col *mongo.Collection
}

func InitVehicleStore(col *mongo.Collection) IVehicle {
	return &VehicleStore{
		col: col,
	}
}

// HANDLERS
func (s *VehicleStore) UpdateDeviceStatus(imei string, status bool) error {
	filter := bson.D{{Key: "imei", Value: imei}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "device", Value: status}}}}

	_, err := s.col.UpdateOne(context.Background(), filter, update)
	return err
}
func (s *VehicleStore) UpdateCurrentStatus(imei string, eng bool, gsm, volt int64, acc, gps bool) error {
	filter := bson.D{{Key: "imei", Value: imei}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "engine", Value: eng},
		{Key: "gsmSignle", Value: int(gsm)},
		{Key: "voltageLevel", Value: int(volt)},
		{Key: "accStatus", Value: acc},
		{Key: "gpsStatus", Value: gps},
	}}}

	_, err := s.col.UpdateOne(context.Background(), filter, update)
	return err
}
