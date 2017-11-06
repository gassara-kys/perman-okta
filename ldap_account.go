package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"

	"gopkg.in/ldap.v2"
)

const (
	// CreateKey for create data
	CreateKey = "CREATE"
	// UpdateKey for update data
	UpdateKey = "UPDATE"
	// DeleteKey for delete data
	DeleteKey = "DELETE"
)

// LdapAccount Permanアカウント
type LdapAccount struct {
	Dn             string   `json:"dn"`
	UID            string   `json:"uid"`
	Email          string   `json:"email"`
	EmployeeNumber string   `json:"employeeNumber"`
	Descriptions   []string `json:"descriptions"`
}

// Convert ldapsearchの結果をLdapAccount型に変換します。
func (ldapAccount LdapAccount) Convert(entries []*ldap.Entry) *[]LdapAccount {
	accounts := []LdapAccount{}
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
	return &accounts

}

// OutJSON jsonファイルに吐き出します
func (ldapAccount LdapAccount) OutJSON(fileNm string, accounts *[]LdapAccount) (err error) {

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

// Diff 差分をチェックして作成、修正、削除が必要なLdapAccountを返します。
func (ldapAccount LdapAccount) Diff(old, new *[]LdapAccount) (result map[string][]LdapAccount, err error) {

	result = make(map[string][]LdapAccount)
	commonLdapData := make(map[string]LdapAccount)

	// oldとnewで同一DNの更新をチェック、差分がある場合は更新リストに追加
	for _, oldData := range *old {
		for _, newData := range *new {
			if oldData.Dn != newData.Dn {
				continue
			}
			commonLdapData[newData.Dn] = newData
			if !reflect.DeepEqual(oldData, newData) {
				result[UpdateKey] = append(result[UpdateKey], newData)
			}
			break
		}
	}
	// old側にしか存在しないデータは削除リストに追加
	for _, data := range *old {
		if _, ok := commonLdapData[data.Dn]; !ok {
			result[DeleteKey] = append(result[DeleteKey], data)
		}
	}
	// new側にしか存在しないデータは新規作成リストに追加
	for _, data := range *new {
		if _, ok := commonLdapData[data.Dn]; !ok {
			result[CreateKey] = append(result[CreateKey], data)
		}
	}

	return result, nil
}
