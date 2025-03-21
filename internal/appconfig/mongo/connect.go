package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"reup/config"
	pkgCrt "reup/pkg/encrypter"
	"reup/pkg/mongo"
)

const (
	connectTimeout = 10 * time.Second
)

// Connect connects to the database
func Connect(mongoConfig config.MongoConfig, encrypter pkgCrt.Encrypter) (mongo.Client, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), connectTimeout)
	defer cancelFunc()
	fmt.Println("mongoConfig.ENCODED_URI:", mongoConfig.ENCODED_URI)

	uri, err := encrypter.Decrypt(mongoConfig.ENCODED_URI)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt mongo uri: %w", err)
	}

	fmt.Println("==========================uri:", uri)

	opts := mongo.NewClientOptions().
		ApplyURI(uri)

	if mongoConfig.ENABLE_MONITOR {
		opts.SetMonitor(mongo.CommandMonitor{
			Started: func(ctx context.Context, e *mongo.CommandStartedEvent) {
				log.Printf("MongoDB command started: %v", e.Command)
			},
			Succeeded: func(ctx context.Context, e *mongo.CommandSucceededEvent) {
				log.Printf("MongoDB command succeeded: %v", e.Reply)
			},
			Failure: func(ctx context.Context, e *mongo.CommandFailedEvent) {
				log.Printf("MongoDB command failed: %v", e.Failure)
			},
		})
	}

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	err = client.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping to DB: %w", err)
	}

	log.Println("Connected to MongoDB!")

	return client, nil
}

// Disconnect disconnects from the database.
func Disconnect(client mongo.Client) {
	if client == nil {
		return
	}

	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connection to MongoDB closed.")
}
