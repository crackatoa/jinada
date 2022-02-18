package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = func() context.Context {
	return context.Background()
}()

type Project struct {
	ID         primitive.ObjectID `bson:"_id"`
	Project    string             `bson:"project"`
	Domain     string             `bson:"domain"`
	Date       time.Time          `bson:"date"`
	Subdomains []SubdomainField   `bson:"subdomains"`
}

type SubdomainField struct {
	Subdomain string `bson:"subdomain"`
	Title     string `bson:"title"`
}

func connect() (*mongo.Database, error) {
	clientOptions := options.Client()
	clientOptions.ApplyURI("mongodb://localhost:27017")
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

// list
func List(list string) {
	workspaceCollection := connectCollection("workspace")

	results, err := workspaceCollection.Distinct(ctx, list, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	for _, result := range results {
		fmt.Println(result)
	}
}

// Find query
func FindAllProject() {
	workspaceCollection := connectCollection("workspace")

	opts := options.Find().SetProjection(bson.M{"project": 1})

	cursor, err := workspaceCollection.Find(ctx, bson.D{}, opts)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var project bson.M
		if err = cursor.Decode(&project); err != nil {
			log.Fatal(err)
		}
		fmt.Println(project["project"])
	}
}

func FindAllDomain() {
	workspaceCollection := connectCollection("workspace")

	opts := options.Find().SetProjection(bson.M{"domain": 1})

	cursor, err := workspaceCollection.Find(ctx, bson.D{}, opts)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)
	var results []bson.D
	if err = cursor.All(ctx, &results); err != nil {
		log.Fatal(err)
	}

	for _, result := range results {
		fmt.Println(result[1])
	}

}

func FindAllSubdomain() {
	workspaceCollection := connectCollection("workspace")

	opts := options.Find().SetProjection(bson.M{"subdomains.subdomain": 1})

	cursor, err := workspaceCollection.Find(ctx, bson.D{}, opts)
	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(ctx)
	var result []Project
	if err = cursor.All(ctx, &result); err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
}

func FindScopeDomain(domain string) {

}

func CheckProjectExist(project string) {

}

func CheckDomainExist(domain string) {

}

// Insert query
func InsertProjectDomain(project string, domain string) {
	workspaceCollection := connectCollection("workspace")
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"domain": domain}
	update := bson.D{{"$set", bson.M{"project": project, "domain": domain}}}
	result, err := workspaceCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Fatal(err)
	}
	if result.UpsertedCount == 0 {
		log.Fatal("Project Exist")
	}
	fmt.Println("Project Created", result)

}

func InsertDomain(project string, domain string) {

}

func InsertSubdomain(domain string, subdomain string) {
	workspaceCollection := connectCollection("workspace")

	data := SubdomainField{
		Subdomain: subdomain,
	}
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"domain": domain}
	update := bson.D{{"$addToSet", bson.M{"subdomains": data}}}
	result, err := workspaceCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Fatal("Subdomain Exist")
	}
	if result.ModifiedCount == 1 {
		fmt.Println("New Subdomain", subdomain, "added")
	}
}

func InsertDummy() {
	workspaceCollection := connectCollection("workspace")
	data := Project{
		ID:      primitive.NewObjectID(),
		Project: "Test1",
		Domain:  "test1.com",
		Date:    time.Now(),
		Subdomains: []SubdomainField{
			{Subdomain: "a.test1.com", Title: "test 1"},
			{Subdomain: "b.test.com", Title: "test2"},
		},
	}

	result, err := workspaceCollection.InsertOne(ctx, data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)

}

// Delete query
func DeleteProject(project string) {
	workspaceCollection := connectCollection("workspace")
	var remove = bson.M{"project": project}
	result, err := workspaceCollection.DeleteOne(ctx, remove)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Removed :", result.DeletedCount)
}

func DeleteDomain(project string, domain string) {

}

func DeleteSubdomain(domain string, subdomain string) {

}

// Update
func UpdateProject(project string) {

}

func UpdateDomain(domain string) {

}

func UpdateSubdomain(subdomain string) {

}
