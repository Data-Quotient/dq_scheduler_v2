package config

import (
	"context"
	"dq_scheduler_v2/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	db *mongo.Database
}

func NewConfig(mongoURI string) (*Config, error) {
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	db := client.Database("scheduler_db")
	return &Config{db: db}, nil
}

func (c *Config) LoadSchedulers() ([]model.Scheduler, error) {
	collection := c.db.Collection("schedulers")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var schedulers []model.Scheduler
	for cursor.Next(context.Background()) {
		var scheduler model.Scheduler
		err := cursor.Decode(&scheduler)
		if err != nil {
			return nil, err
		}
		schedulers = append(schedulers, scheduler)
	}
	return schedulers, nil
}

func (c *Config) SaveScheduler(scheduler model.Scheduler) error {
	collection := c.db.Collection("schedulers")
	_, err := collection.InsertOne(context.Background(), scheduler)
	return err
}

func (c *Config) UpdateScheduler(scheduler model.Scheduler) error {
	collection := c.db.Collection("schedulers")
	_, err := collection.ReplaceOne(context.Background(), bson.M{"id": scheduler.ID}, scheduler)
	return err
}

func (c *Config) DeleteScheduler(schedulerID string) error {
	collection := c.db.Collection("schedulers")
	_, err := collection.DeleteOne(context.Background(), bson.M{"id": schedulerID})
	return err
}
