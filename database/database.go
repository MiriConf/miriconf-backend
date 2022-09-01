package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoConnection() *mongo.Client {
  mongoURI := os.Getenv("MONGO_URI")
	
  client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
  if err != nil {
    log.Fatal(err)
  }

  ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
  
  err = client.Connect(ctx)
  if err != nil {
    log.Fatal(err)
  }
  
  defer client.Disconnect(ctx)
  
  return client
}
