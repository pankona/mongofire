package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/sync/errgroup"
)

func doSomething(index int, loopCount int) error {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://mongo:27017"))

	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.Connect(ctx)

	if err != nil {
		return err
	}

	collection := client.Database("testing").Collection("numbers")

	for i := 0; i < loopCount; i++ {
		_, err = collection.UpdateOne(ctx, bson.M{"index": bson.M{"$eq": index}}, bson.M{"$set": bson.M{"value": i}})
		if err != nil {
			return err
		}
	}

	return nil
}

func setup(count int) error {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://mongo:27017"))

	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	err = client.Connect(ctx)

	if err != nil {
		return err
	}

	collection := client.Database("testing").Collection("numbers")

	err = collection.Drop(ctx)

	if err != nil {
		return err
	}

	for i := 0; i < count; i++ {
		_, err = collection.InsertOne(ctx, bson.M{"index": i, "value": 0})
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	connectionCount := 30
	loopCount := 1000

	setup(connectionCount)

	eg := errgroup.Group{}

	for i := 0; i < connectionCount; i++ {
		eg.Go(func() error {
			return doSomething(i, loopCount)
		})
	}

	err := eg.Wait()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done")
}
