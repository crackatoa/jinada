package main

import (
	"flag"
	"jinada/database"
)

func main() {
	var project string
	flag.StringVar(&project, "p", "", "project name")

	var domain string
	flag.StringVar(&domain, "d", "", "domain name")

	var new bool
	flag.BoolVar(&new, "new", false, "create new project")

	var list string
	flag.StringVar(&list, "list", "", "list asset")

	var dummy bool
	flag.BoolVar(&dummy, "dummy", false, "insert dummy data")

	var addsubdomain string
	flag.StringVar(&addsubdomain, "addsubdomain", "", "add subdomain")

	flag.Parse()

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
		database.InsertProject(project)
	}

	if dummy == true {
		database.InsertDummy()
	}

	if len(domain) > 0 {

	}

}
