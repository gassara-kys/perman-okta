package main

import (
	"encoding/json"
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
		Host:       os.Getenv("LDAP_HOST"),
		BaseDn:     os.Getenv("BASE_DN"),
		Filter:     os.Getenv("FILTER_STRING"),
		Attributes: []string{"dn", "uid", "email", "employeeNumber", "description"},
		SizeLimit:  noSizeLimit,
		TimeLimit:  noTimeLimit,
		TypeOnly:   noTypeOnly,
	}
	result, err := ldapClient.Search()
	if err != nil {
		log.Fatal(err)
	}
	// get ldap datas
	var account = Account{}
	serverData := account.ConvertFromLdap(result.Entries)
	if len(*serverData) == 0 {
		log.Fatal("LDAP Server Account is 0...") // LDAPサーバーのアカウント0件は異常終了にする
	}
	localData, err := account.LoadJSON(fileNm)
	if err != nil {
		log.Fatal(err)
	}

	// diff
	diff, err1 := account.Diff(localData, serverData)
	if err1 != nil {
		log.Fatal(err)
	}

	for _, data := range diff[CreateKey] {
		log.Printf("[%s]%s", CreateKey, data.Dn)
	}
	for _, data := range diff[UpdateKey] {
		log.Printf("[%s]%s", UpdateKey, data.Dn)
	}
	for _, data := range diff[DeleteKey] {
		log.Printf("[%s]%s", DeleteKey, data.Dn)
	}

	// output JSON file
	if err := account.OutJSON(fileNm, serverData); err != nil {
		log.Fatal(err)
	}

	// Okta API
	var oktaClient = OktaClient{
		FQDN:   os.Getenv("OKTA_FQDN"),
		APIKEY: os.Getenv("OKTA_APIKEY"),
	}
	// oktaUser, oktaErr := oktaClient.GetUser("azuma_kazuhiro@cyberagent.co.jp")
	oktaUser, oktaErr := oktaClient.GetUser("no_user@cyberagent.co.jp")
	if oktaErr != nil {
		log.Fatal(oktaErr)
	}

	jsonByte, _ := json.MarshalIndent(oktaUser, "", "  ")
	log.Printf("user count: %d", len(*oktaUser))
	log.Printf("%s", jsonByte)

}
