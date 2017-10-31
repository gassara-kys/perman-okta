package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/ldap.v2"
)

const (
	sizeLimit = 0
	timeLimit = 0
	typeOnly  = false
)

func main() {

	host := os.Getenv("LDAP_HOST")
	baseDn := os.Getenv("BASE_DN")
	filter := os.Getenv("FILTER_STRING")

	log.Println("start ... ")
	log.Printf("LDAP_HOST: %s\n", host)
	log.Printf("LDAP_BASE_DN: %s\n", baseDn)
	log.Printf("FILTER_STRING: %s\n", filter)

	ldapConn, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", host, 389))
	if err != nil {
		log.Fatalf("connerction Error... err: %+v", err)
	}
	defer ldapConn.Close()

	// ldapsearch
	searchRequest := ldap.NewSearchRequest(
		baseDn,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		sizeLimit,
		timeLimit,
		typeOnly,
		filter,
		[]string{"dn", "uid", "email", "employeeNumber", "description"},
		nil,
	)
	result, err := ldapConn.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}

	// output ldap entries
	log.Println("Let's show ldap entries ...")
	for _, entry := range result.Entries {
		fmt.Printf("dn: %s\n", entry.DN)
		fmt.Printf("uid: %s\n", entry.GetAttributeValue("uid"))
		fmt.Printf("email: %s\n", entry.GetAttributeValue("email"))
		fmt.Printf("employeeNumber: %s\n", entry.GetAttributeValue("employeeNumber"))
		descriptions := entry.GetAttributeValues("description")
		for idx, desc := range descriptions {
			fmt.Printf("description_%02d: %s\n", idx, desc)
		}
	}
}
