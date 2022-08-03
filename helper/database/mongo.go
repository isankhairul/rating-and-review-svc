package database

import (
	"context"
	"fmt"
	"go-klikdokter/helper/config"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	DbMongo *mongo.Database
)

func NewMongo() (*mongo.Database, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=%s", viper.GetString("database.username"), viper.GetString("database.password"), viper.GetString("database.hostname"), viper.GetString("database.port"), viper.GetString("database.dbname"))
	uri := fmt.Sprintf("%s://%s:%s@%s/%s?authSource=admin", config.GetConfigString(viper.GetString("database.uri")), config.GetConfigString(viper.GetString("database.username")), config.GetConfigString(viper.GetString("database.password")), config.GetConfigString(viper.GetString("database.hostname")), config.GetConfigString(viper.GetString("database.dbname")))
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	err = CreateIndex(client, "ratingTypesNumCol", "type", true)
	if err != nil {
		return nil, err
	}
	err = CreateIndex(client, "ratingTypesLikertCol", "type", true)
	if err != nil {
		return nil, err
	}
	err = CreateIndex(client, "ratingsCol", "name", true)
	if err != nil {
		return nil, err
	}
	err = CreateIndex(client, "ratingSubCol", "source_trans_id", true)
	if err != nil {
		return nil, err
	}
	err = CreateIndex(client, "ratingSubCol", "rating_id", false)
	if err != nil {
		return nil, err
	}
	err = CreateIndexRatingsCol(client)
	if err != nil {
		return nil, err
	}

	return client.Database(config.GetConfigString(viper.GetString("database.dbname"))), nil
}

func CreateIndex(client *mongo.Client, collection, col string, isUnique bool) error {
	_, err := client.Database(config.GetConfigString(viper.GetString("database.dbname"))).Collection(collection).Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: col, Value: 1}},
			Options: options.Index().SetUnique(isUnique),
		},
	)
	return err
}

func CreateIndexRatingsCol(client *mongo.Client) error {
	_, err := client.Database(config.GetConfigString(viper.GetString("database.dbname"))).Collection("ratingsCol").Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.D{{Key: "source_uid", Value: 1}, {Key: "source_type", Value: 1}, {Key: "update_at", Value: 1}},
		},
	)
	return err
}
