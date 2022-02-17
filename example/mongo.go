package example

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = func() context.Context {
	return context.Background()
}()

type student struct {
	Name  string `bson:"name"`
	Grade int    `bson:"Grade"`
}

type another struct {
	Name  string `bson:"name"`
	Kelas int    `bson: "Grade"`
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

	//return client.Database("jinada_dev"), nil
	return client.Database("jinada"), nil
}

func Insert() {
	db, err := connect()
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = db.Collection("student").InsertOne(ctx, student{"Wick", 2})
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = db.Collection("student").InsertOne(ctx, another{"Ethan", 2})
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Insert success!")
}

func InsertProd(project string, domain string) {
	db, err := connect()
	if err != nil {
		log.Fatal(err.Error())
	}

	data := bson.D{
		{"project", project},
		{"domain", domain},
		{"date", time.Now()},
		{"subdomains", bson.A{
			bson.D{
				{"subdomain", "test"},
				{"title", "test"},
				{"response", "hallo"},
			},
			bson.D{
				{"subdomain", "test2"},
				{"title", "test"},
				{"response", "hallo"},
			},
		}},
	}

	_, err = db.Collection("workspace").InsertOne(ctx, data)
	if err != nil {
		log.Fatal(err)
	}

}

func DeleteProject(project string) {
	db, err := connect()
	if err != nil {
		log.Fatal(err.Error())
	}

	var selector = bson.M{"project": project}
	_, err = db.Collection("workspace").DeleteOne(ctx, selector)
	if err != nil {
		log.Fatal(err.Error())
	} else {
		fmt.Println("project removed")
	}

}

func UpdateOne(domain string, newdomain string) {
	db, err := connect()
	if err != nil {
		log.Fatal(err.Error())
	}

	var selector = bson.M{"domain": domain}
	_, err = db.Collection("workspace").UpdateOne(ctx, selector, bson.M{"$set": bson.M{"domain": newdomain}})
	if err != nil {
		log.Fatal(err.Error())
	} else {
		fmt.Println("Data Updated")
	}

}

func FindAll() {
	db, err := connect()
	if err != nil {
		log.Fatal(err)
	}

	cursor, err := db.Collection("workspace").Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	var project []bson.M
	if err = cursor.All(ctx, &project); err != nil {
		log.Fatal(err)
	}

	fmt.Println(project)
}

func FindOne() {
	db, err := connect()
	if err != nil {
		log.Fatal(err)
	}

	var project bson.M
	if err = db.Collection("workspace").FindOne(ctx, bson.M{}).Decode(&project); err != nil {
		log.Fatal(err)
	}
	fmt.Println(project)
}

func FindOneDomain(domain string) {
	db, err := connect()
	if err != nil {
		log.Fatal(err)
	}

	var project bson.M
	if err = db.Collection("workspace").FindOne(ctx, bson.M{"domain": domain}).Decode(&project); err != nil {
		log.Fatal(err)
	}
	fmt.Println(project)
}

func CheckDomainExist(domain string) bool {
	db, err := connect()
	if err != nil {
		log.Fatal()
	}
	var result bson.M
	err = db.Collection("workspace").FindOne(ctx, bson.M{"domain": domain}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
			log.Fatal(err)
		}
	}
	return true
}

func FindDomain(domain string) {
	db, err := connect()
	if err != nil {
		log.Fatal(err)
	}

	filter, err := db.Collection("workspace").Find(ctx, bson.M{"domain": domain})
	if err != nil {
		log.Fatal(err)
	}
	var domainexist []bson.M
	if err = filter.All(ctx, &domainexist); err != nil {
		log.Fatal(err)
	}

	fmt.Println(domainexist)
}
