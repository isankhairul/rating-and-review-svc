package database

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

var (
	DbMongo *mongo.Database
)

func NewMongo() (*mongo.Database, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=%s", viper.GetString("database.username"), viper.GetString("database.password"), viper.GetString("database.hostname"), viper.GetString("database.port"), viper.GetString("database.dbname"))
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	// Create index
	_, err = client.Database(viper.GetString("database.dbname")).Collection("ratingTypesNumCol").Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "type", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		return nil, err
	}

	return client.Database(viper.GetString("database.dbname")), nil
}
