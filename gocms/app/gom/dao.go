package gom

import (
	"context"
	"fmt"
	"time"

	"github.com/davidddw/gopj/gocms/app/conf"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoClient  *mongo.Client
	mongoContext context.Context
)

func init() {
	mongoContext, _ = context.WithTimeout(context.Background(), 30*time.Second)
	opts := &options.ClientOptions{}
	opts.SetAuth(options.Credential{
		AuthMechanism: "SCRAM-SHA-1",
		AuthSource:    conf.Conf.Db.DbName,
		Username:      conf.Conf.Db.Username,
		Password:      conf.Conf.Db.Password,
	}).ApplyURI(fmt.Sprintf("mongodb://%s:27017", conf.Conf.Db.Host))
	client, err := mongo.Connect(mongoContext, opts)
	if err != nil {
		fmt.Println(err)

	}
	mongoClient = client
}

func GetSession() (*mongo.Client, context.Context) {
	return mongoClient, mongoContext
}

func getCollect(client *mongo.Client, collectionName string) *mongo.Collection {
	return client.Database(conf.Conf.Db.DbName).Collection(collectionName)
}
