package database

import (
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var resultFindSubdomain struct {
	ID       string   `bson:"_id"`
	Sudomain []string `bson:"subdomain"`
}

var resultFindIp struct {
	ID string   `bson:"_id"`
	IP []string `bson:"ip"`
}

func InsertDomainToProject(project string, domain string) {

	data := DomainField{
		Domain:     domain,
		Subdomains: []SubdomainField{},
	}

	opts := options.Update().SetUpsert(true)
	filter := bson.M{"project": project}
	update := bson.D{{"$addToSet", bson.M{"domains": data}}}
	result, err := projectsCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Fatal(err)
	}
	if result.ModifiedCount == 1 {
		fmt.Println("New Domain", domain, "added")
	} else if result.UpsertedCount == 1 {
		fmt.Println("Create new project", project, "\nNew domain", domain)
	} else {
		fmt.Println("Domain", domain, "exist")
	}
}

func InsertIpToProject(project string, ip string) {

	data := IPField{
		IP: ip,
	}

	opts := options.Update().SetUpsert(true)
	filter := bson.M{"project": project}
	update := bson.D{{"$addToSet", bson.M{"ips": data}}}
	result, err := projectsCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Fatal(err)
	}
	if result.ModifiedCount == 1 {
		fmt.Println("New IP", ip, "added")
	} else if result.UpsertedCount == 1 {
		fmt.Println("Create new project", project, "\nNew domain", ip)
	} else {
		fmt.Println("IP", ip, "exist")
	}
}

func InsertSubdomainToDomain(domain string, subdomain string) {
	data := SubdomainField{
		Subdomain: subdomain,
	}

	opts := options.Update().SetUpsert(true)
	filter := bson.M{"domains.domain": domain}
	update := bson.D{{"$addToSet", bson.M{"domains.$.subdomains": data}}}
	result, err := projectsCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Fatal(err)
	}

	if result.ModifiedCount == 1 {
		fmt.Println("New Subdomain", subdomain, "added to", domain)
	} else {
		fmt.Println(subdomain, "already exist")
	}
}

func InsertSubdomainToDomainMany(domain string, subdomain []string) {
	for _, s := range subdomain {
		data := SubdomainField{
			Subdomain: s,
		}

		opts := options.Update().SetUpsert(true)
		filter := bson.M{"domains.domain": domain}
		update := bson.D{{"$addToSet", bson.M{"domains.$.subdomains": data}}}
		result, err := projectsCollection.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			log.Fatal(err)
		}

		if result.ModifiedCount == 1 {
			fmt.Println("New Subdomain", s, "added to", domain)
		} else {
			fmt.Println(s, "already exist")
		}
	}
}

func InsertDummyProjects() {
	data := Projects{
		ID:         primitive.NewObjectID(),
		Project:    "project new",
		LastUpdate: time.Now(),
		Domains: []DomainField{
			{
				Domain: "test.com",
				Subdomains: []SubdomainField{
					{
						Subdomain: "a.test.com",
					},
				},
			},
		},
		IPs: []IPField{
			{
				IP: "10.10.10.1",
			},
		},
	}

	result, err := projectsCollection.InsertOne(ctx, data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("result :", result.InsertedID)
}

func DeleteProject(project string) {
	var remove = bson.M{"project": project}
	result, err := projectsCollection.DeleteOne(ctx, remove)
	if err != nil {
		log.Fatal(err)
	}
	if result.DeletedCount == 1 {
		fmt.Println("Project", project, "removed")
	} else {
		fmt.Println("Project", project, "not exist")
	}
}

func ListProject() {
	opts := options.Find().SetProjection(bson.M{"project": 1})
	cursor, err := projectsCollection.Find(ctx, bson.M{}, opts)
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

func ListDomain() {
	unwind1 := bson.D{{"$unwind", "$domains"}}
	group := bson.D{{"$group", bson.D{{"_id", "$domains.domain"}}}}

	cursor, err := projectsCollection.Aggregate(ctx, mongo.Pipeline{unwind1, group})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		if err = cursor.Decode(&resultFindSubdomain); err != nil {
			log.Fatal(err)
		}
		fmt.Println(resultFindSubdomain.ID)
	}
}

func ListSubdomain() {
	unwind1 := bson.D{{"$unwind", "$domains"}}
	unwind2 := bson.D{{"$unwind", "$domains.subdomains"}}
	group := bson.D{{"$group", bson.D{{"_id", "$domains.domain"}, {"subdomain", bson.D{{"$push", "$domains.subdomains.subdomain"}}}}}}

	cursor, err := projectsCollection.Aggregate(ctx, mongo.Pipeline{unwind1, unwind2, group})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		if err = cursor.Decode(&resultFindSubdomain); err != nil {
			log.Fatal(err)
		}
		for _, j := range resultFindSubdomain.Sudomain {
			fmt.Println(resultFindSubdomain.ID, j)
		}
	}
}

func ListSubdomainDomain(domain string) {

	unwind1 := bson.D{{"$unwind", "$domains"}}
	unwind2 := bson.D{{"$unwind", "$domains.subdomains"}}
	match := bson.D{{"$match", bson.D{{"domains.domain", domain}}}}
	group := bson.D{{"$group", bson.D{{"_id", "$domains.domain"}, {"subdomain", bson.D{{"$push", "$domains.subdomains.subdomain"}}}}}}

	cursor, err := projectsCollection.Aggregate(ctx, mongo.Pipeline{unwind1, unwind2, match, group})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		if err = cursor.Decode(&resultFindSubdomain); err != nil {
			log.Fatal(err)
		}
		for _, j := range resultFindSubdomain.Sudomain {
			fmt.Println(j)
		}
	}
}

func ListIp() {

	unwind1 := bson.D{{"$unwind", "$ips"}}
	group := bson.D{{"$group", bson.D{{"_id", "$project"}, {"ip", bson.D{{"$push", "$ips.ip"}}}}}}

	cursor, err := projectsCollection.Aggregate(ctx, mongo.Pipeline{unwind1, group})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		if err = cursor.Decode(&resultFindIp); err != nil {
			log.Fatal(err)
		}
		for _, i := range resultFindIp.IP {
			fmt.Println(resultFindIp.ID, i)
		}
	}
}
