package db_client

import (
	pkgMongo "github.com/404nffff/go_pkg/mongo"

	"go.mongodb.org/mongo-driver/mongo"
)

func MongoLocal() *mongo.Database {

	return pkgMongo.NewClient("Local")
}
