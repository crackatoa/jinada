package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.Background()
var workspaceCollection *mongo.Collection
var projectsCollection *mongo.Collection

type Project struct {
	ID         primitive.ObjectID `bson:"_id"`
	Project    string             `bson:"project"`
	Domain     string             `bson:"domain"`
	Date       time.Time          `bson:"date"`
	Subdomains []SubdomainField   `bson:"subdomains"`
}

type Projects struct {
	ID         primitive.ObjectID `bson:"_id"`
	Project    string             `bson:"project"`
	LastUpdate time.Time          `bson:"last_update"`
	Domains    []DomainField      `bson: "domains"`
	IPs        []IPField          `bson:"ips"`
}

type DomainField struct {
	Domain     string           `bson:"domain"`
	Subdomains []SubdomainField `bson:"subdomains"`
}

type SubdomainField struct {
	Subdomain string `bson:"subdomain"`
}

type IPField struct {
	IP string `bson:"ip"`
}

func Init() {
	projectsCollection = connectCollection("projects")
}

func connect() (*mongo.Database, error) {
	clientOptions := options.Client()
	clientOptions.SetRetryWrites(true)
	clientOptions.ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return client.Database("jinada"), nil
}

func connectCollection(collection string) *mongo.Collection {
	db, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	return db.Collection(collection)
}
