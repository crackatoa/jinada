package main

import (
	"bufio"
	"flag"
	"fmt"
	"jinada/database"
	"log"
	"os"
)

func main() {
	var project string
	flag.StringVar(&project, "p", "", "project name")

	var domain string
	flag.StringVar(&domain, "d", "", "domain name")

	var new bool
	flag.BoolVar(&new, "new", false, "create new project")

	var remove bool
	flag.BoolVar(&remove, "remove", false, "delete project")

	var list string
	flag.StringVar(&list, "list", "", "list asset")

	var dummy bool
	flag.BoolVar(&dummy, "dummy", false, "insert dummy data")

	var addsubdomain string
	flag.StringVar(&addsubdomain, "addsubdomain", "", "add subdomain")

	// handler
	var input string
	flag.StringVar(&input, "i", "", "read input from file")

	flag.Parse()

	//input from stdin

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if len(domain) > 0 {
			database.InsertSubdomain(domain, scanner.Text())
		} else {
			fmt.Println("Domain not set")
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	//flag list action
	if len(list) > 0 {
		switch list {
		case "project":
			database.List("project")
		case "domain":
			database.List("domain")
		case "subdomain":
			database.FindAllSubdomain()

		}
	}

	//flag addsubdomain action
	if len(addsubdomain) > 0 {
		database.InsertSubdomain(domain, addsubdomain)
	}

	if new == true {
		if len(domain) > 0 {
			database.InsertProjectDomain(project, domain)
		} else {
			fmt.Println("Domain not set")
		}
	}

	if remove == true {
		if len(project) > 0 {
			database.DeleteProject(project)
		} else {
			fmt.Println("Project not set")
		}
	}

	if dummy == true {
		database.InsertDummy()
	}

	//handler

}
