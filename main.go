package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/crackatoa/jinada/database"
)

func main() {
	var project string
	var domain string
	var ip string
	var new bool
	var remove bool
	var list string
	var dummy bool
	var subdomain string
	var input string
	var add bool

	flag.StringVar(&project, "p", "", "project name")
	flag.StringVar(&domain, "d", "", "domain name")
	flag.StringVar(&ip, "ip", "", "ip target")
	flag.BoolVar(&new, "new", false, "create new project")
	flag.BoolVar(&remove, "remove", false, "delete project")
	flag.BoolVar(&add, "add", false, "add subdomain to domain")
	flag.StringVar(&list, "list", "", "list asset")
	flag.BoolVar(&dummy, "dummy", false, "insert dummy data")
	flag.StringVar(&subdomain, "subdomain", "", "add subdomain")
	flag.StringVar(&input, "i", "", "read input from file")
	flag.Parse()

	database.Init()
	//flag list action
	if len(list) > 0 {
		switch list {
		case "project":
			database.ListProject()
		case "domain":
			database.ListDomain()
		case "subdomain":
			if len(domain) > 0 {
				database.ListSubdomainDomain(domain)
			} else {
				database.ListSubdomain()

			}
		case "ip":
			database.ListIp()
		}
	}

	if new == true {
		if len(domain) > 0 {
			database.InsertDomainToProject(project, domain)
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

	//input from stdin

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		if add {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				if len(domain) > 0 {
					database.InsertSubdomainToDomain(domain, scanner.Text())
				} else {
					fmt.Println("Domain not set")
				}
			}

			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}

		} else {
			fmt.Println("flag -add not set")
		}
	}

	if dummy == true {
		//database.InsertDummyProjects()
		//database.InsertIpToProject(project, ip)
		database.InsertSubdomainToDomain(domain, subdomain)
	}
}
