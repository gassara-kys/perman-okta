package main

import (
	"log"
	"os"
)

const (
	noSizeLimit = 0
	noTimeLimit = 0
	noTypeOnly  = false
	fileNm      = "tmp/ldap_accounts.json"
)

func main() {
	// ldapsearch
	ldapClient := LdapClient{
		os.Getenv("LDAP_HOST"),
		os.Getenv("BASE_DN"),
		os.Getenv("FILTER_STRING"),
		[]string{"dn", "uid", "email", "employeeNumber", "description"},
		noSizeLimit,
		noTimeLimit,
		noTypeOnly,
	}
	result, err := ldapClient.Search()
	if err != nil {
		log.Fatal(err)
	}

	// output ldap data to json file
	var ldapAccount = LdapAccount{}
	if err := ldapAccount.OutJSON(fileNm, result.Entries); err != nil {
		log.Fatal(err)
	}

	// load from json
	accouts, err := ldapAccount.LoadJSON(fileNm)
	if err != nil {
		log.Fatal(err)
	}

	for _, data := range *accouts {
		log.Printf("%s", data.Dn)
	}

}
