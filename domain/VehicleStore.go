package domain

import (
	"context"
	"mtrack/device_server/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type IVehicle interface {
	GetVehicle(imei string) (types.Vehicle, error)

	UpdateLocationInfo(imei, tripId string, tm int64) error
	UpdateDeviceStatus(imei string, status bool) error
	UpdateCurrentStatus(imei string, eng bool, gsm, volt int64, acc, gps bool) error
}

type VehicleStore struct {
	v_col *mongo.Collection
	t_col *mongo.Collection
}

func InitVehicleStore(v_col, t_col *mongo.Collection) IVehicle {
	return &VehicleStore{
		v_col: v_col,
		t_col: t_col,
	}
}

// GETTERS
func (s *VehicleStore) GetVehicle(imei string) (types.Vehicle, error) {
	var result types.Vehicle
	err := s.v_col.FindOne(context.Background(), bson.M{"imei": imei}).Decode(&result)
	return result, err
}

// SETTERS
func (s *VehicleStore) UpdateLocationInfo(imei, tripId string, tm int64) error {
	filter := bson.D{{Key: "imei", Value: imei}}
	update := bson.D{
		{Key: "lastTripID", Value: tripId},
		{Key: "lastGpsUpdate", Value: tm},
	}
	_, err := s.v_col.UpdateOne(context.Background(), filter, bson.D{{Key: "$set", Value: update}})
	return err
}
func (s *VehicleStore) UpdateDeviceStatus(imei string, status bool) error {
	filter := bson.D{{Key: "imei", Value: imei}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "device", Value: status}}}}

	_, err := s.v_col.UpdateOne(context.Background(), filter, update)
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

	_, err := s.v_col.UpdateOne(context.Background(), filter, update)
	return err
}
