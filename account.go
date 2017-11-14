package main

import (
	"encoding/json"
	"io/ioutil"
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

// Account Permanアカウント
type Account struct {
	Dn             string   `json:"dn"`
	UID            string   `json:"uid"`
	Email          string   `json:"email"`
	EmployeeNumber string   `json:"employeeNumber"`
	Descriptions   []string `json:"descriptions"`
}

// Convert ldapsearchの結果をAccount型に変換します。
func (a Account) ConvertFromLdap(entries []*ldap.Entry) *[]Account {
	accounts := []Account{}
	for _, entry := range entries {
		var account = Account{}
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
func (a Account) OutJSON(fileNm string, accounts *[]Account) (err error) {

	jsonBytes, err := json.MarshalIndent(accounts, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fileNm, jsonBytes, 0644)

}

// LoadJSON JSONから読み込みます
func (a Account) LoadJSON(fileNm string) (accouts *[]Account, err error) {

	data, err := ioutil.ReadFile(fileNm)
	if err != nil {
		return
	}
	var accounts []Account
	err = json.Unmarshal(data, &accounts)
	if err != nil {
		return
	}
	return &accounts, nil
}

// Diff 差分をチェックして作成、修正、削除が必要なLdapAccountを返します。
func (a Account) Diff(old, new *[]Account) (result map[string][]Account, err error) {

	result = make(map[string][]Account)
	commonData := make(map[string]Account)

	// oldとnewで同一DNの更新をチェック、差分がある場合は更新リストに追加
	for _, oldData := range *old {
		for _, newData := range *new {
			if oldData.Dn != newData.Dn {
				continue
			}
			commonData[newData.Dn] = newData
			if !reflect.DeepEqual(oldData, newData) {
				result[UpdateKey] = append(result[UpdateKey], newData)
			}
			break
		}
	}
	// old側にしか存在しないデータは削除リストに追加
	for _, data := range *old {
		if _, ok := commonData[data.Dn]; !ok {
			result[DeleteKey] = append(result[DeleteKey], data)
		}
	}
	// new側にしか存在しないデータは新規作成リストに追加
	for _, data := range *new {
		if _, ok := commonData[data.Dn]; !ok {
			result[CreateKey] = append(result[CreateKey], data)
		}
	}

	return result, nil
}
