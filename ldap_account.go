package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"gopkg.in/ldap.v2"
)

// LdapAccount Permanアカウント
type LdapAccount struct {
	Dn             string   `json:"dn"`
	UID            string   `json:"uid"`
	Email          string   `json:"email"`
	EmployeeNumber string   `json:"employeeNumber"`
	Descriptions   []string `json:"descriptions"`
}

// OutJSON jsonファイルに吐き出します
func (ldapAccount LdapAccount) OutJSON(fileNm string, entries []*ldap.Entry) (err error) {
	var accounts = []LdapAccount{}
	for _, entry := range entries {
		var account = LdapAccount{}
		account.Dn = entry.DN
		account.UID = entry.GetAttributeValue("uid")
		account.Email = entry.GetAttributeValue("email")
		account.EmployeeNumber = entry.GetAttributeValue("employeeNumber")

		descriptions := entry.GetAttributeValues("description")
		for _, desc := range descriptions {
			account.Descriptions = append(account.Descriptions, desc)
		}
		accounts = append(accounts, account)
	}

	jsonBytes, err := json.MarshalIndent(accounts, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fileNm, jsonBytes, os.ModePerm)

}

// LoadJSON JSONから読み込みます
func (ldapAccount LdapAccount) LoadJSON(fileNm string) (accouts *[]LdapAccount, err error) {

	data, err := ioutil.ReadFile(fileNm)
	if err != nil {
		return
	}
	var accounts []LdapAccount
	err = json.Unmarshal(data, &accounts)
	if err != nil {
		return
	}
	return &accounts, nil
}
