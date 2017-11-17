package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	noSizeLimit = 0
	noTimeLimit = 0
	noTypeOnly  = false
	fileNm      = "tmp/ldap_accounts.json"
)

func main() {
	EnvLoad(".env")
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

	var jsonByte []byte
	var oktaErr error
	var oktaUser *OktaUser
	var oktaGroup *OktaGroup

	// Create User
	oktaUser, oktaErr = oktaClient.CreateUser(
		&UserProfile{
			LastName:    "OgasawaraTest006",
			SecondEmail: "",
			MobilePhone: "",
			Email:       "ogasawara_kiyoshi+006@cyberagent.co.jp",
			Login:       "ogasawara_kiyoshi+006@cyberagent.co.jp",
			FirstName:   "OgasawaraTest006",
		},
	)
	if oktaErr != nil {
		log.Fatal(oktaErr)
	}
	jsonByte, _ = json.MarshalIndent(oktaUser, "", "  ")
	log.Printf("create user:\n%s", jsonByte)

	// Get User
	oktaUser, oktaErr = oktaClient.GetUserWithLogin("ogasawara_kiyoshi+006@cyberagent.co.jp")
	if oktaErr != nil {
		log.Fatal(oktaErr)
	}
	jsonByte, _ = json.MarshalIndent(oktaUser, "", "  ")
	log.Printf("find user:\n%s", jsonByte)

	// Create Group
	oktaGroup, oktaErr = oktaClient.AddGroup(
		&GroupProfile{
			Name:        "oga-test-group-006",
			Description: "oga-test-group-006",
		},
	)
	if oktaErr != nil {
		log.Fatal(oktaErr)
	}
	jsonByte, _ = json.MarshalIndent(oktaGroup, "", "  ")
	log.Printf("create group:\n%s", jsonByte)

	// Search Group
	oktaGroup, oktaErr = oktaClient.SearchGroups("oga-test-group-006")
	if oktaErr != nil {
		log.Fatal(oktaErr)
	}
	jsonByte, _ = json.MarshalIndent(oktaGroup, "", "  ")
	log.Printf("find group:\n%s", jsonByte)

	// Add member to Group
	oktaErr = oktaClient.AddUserToGroup(oktaGroup.ID, oktaUser.ID)
	if oktaErr != nil {
		log.Fatal(oktaErr)
	}

	// Remove member from Group
	oktaErr = oktaClient.RemoveUserFromGroup(oktaGroup.ID, oktaUser.ID)
	if oktaErr != nil {
		log.Fatal(oktaErr)
	}

	// Delete Group
	// oktaErr = oktaClient.RemoveGroup(oktaGroup.ID)

	// Delete User
	// oktaErr = oktaClient.DeleteUser(oktaUser.ID)
	if oktaErr != nil {
		log.Fatal(oktaErr)
	}

}

// EnvLoad .env load
func EnvLoad(envFile string) {
	if envFile == "" {
		envFile = ".env" // default
	}
	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatalf("Error loading %s file", envFile)
	}
}
