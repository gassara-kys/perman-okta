package main

import (
	"testing"

	ldap "gopkg.in/ldap.v2"
)

const testFileNm = "tmp/test.json"

var testAccounts = []Account{
	{
		Dn:             "uid=aaa_user,dc=example,dc=com",
		UID:            "aaa_user",
		Email:          "aaa_user@example.com",
		EmployeeNumber: "EMP_NO001",
		Descriptions:   []string{"desc_aaa"},
	},
	{
		Dn:             "uid=bbb_user,dc=example,dc=com",
		UID:            "bbb_user",
		Email:          "bbb_user@example.com",
		EmployeeNumber: "EMP_NO002",
		Descriptions:   []string{"desc bbb"},
	},
	{
		Dn:             "uid=ccc_user,dc=example,dc=com",
		UID:            "ccc_user",
		Email:          "ccc_user@example.com",
		EmployeeNumber: "EMP_NO003",
		Descriptions:   []string{"desc ccc"},
	},
}

func TestConvertAndJSON(t *testing.T) {

	var entries []*ldap.Entry
	for _, data := range testAccounts {
		entry := ldap.Entry{
			DN: data.Dn,
			Attributes: []*ldap.EntryAttribute{
				{"uid", []string{data.UID}, [][]byte{}},
				{"email", []string{data.Email}, [][]byte{}},
				{"employeeNumber", []string{data.EmployeeNumber}, [][]byte{}},
				{"description", data.Descriptions, [][]byte{}},
			},
		}
		entries = append(entries, &entry)
	}

	var account = Account{}

	// Convert
	convertData := account.ConvertFromLdap(entries)

	// OutJSON
	err := account.OutJSON(testFileNm, convertData)
	if err != nil {
		t.Errorf("account.OutJSON exec failed: file: %s, data: %v", testFileNm, entries)
	}

	// LoadJSON
	accounts, err := account.LoadJSON(testFileNm)
	if err != nil {
		t.Errorf("account.LoadJSON exec failed: file: %s", testFileNm)
	}

	for idx, data := range *accounts {
		if testAccounts[idx].Dn != data.Dn ||
			testAccounts[idx].UID != data.UID ||
			testAccounts[idx].Email != data.Email ||
			testAccounts[idx].EmployeeNumber != data.EmployeeNumber {
			t.Errorf("JSON data not match: testAccount: %s, jsonAccout: %s", testAccounts[idx].Dn, data.Dn)
		}
	}
}

func TestDiff(t *testing.T) {
	var account = Account{}

	// Pattern: [No Diff]
	var testOld = testAccounts
	var testNew = testAccounts

	result, err := account.Diff(&testOld, &testNew)
	if err != nil {
		t.Errorf("account.Diff [No Diff]exec fatal:%v", err)
	}
	if len(result[UpdateKey]) != 0 {
		t.Errorf("account.Diff [No Diff]update list exists: %v", len(result[UpdateKey]))
	}
	if len(result[CreateKey]) != 0 {
		t.Errorf("account.Diff [No Diff]create list exists: %v", len(result[CreateKey]))
	}
	if len(result[DeleteKey]) != 0 {
		t.Errorf("account.Diff [No Diff]delete list exists: %v", len(result[DeleteKey]))
	}

	// Pattern: [Modify]
	testOld = []Account{
		{
			Dn:             "uid=aaa_user,dc=example,dc=com",
			UID:            "aaa_user",
			Email:          "aaa_user@example.com",
			EmployeeNumber: "EMP_NO001",
			Descriptions:   []string{"desc_aaa"},
		},
		{
			Dn:             "uid=bbb_user,dc=example,dc=com",
			UID:            "xxx_user",             // mod
			Email:          "xxx_user@example.com", //mod
			EmployeeNumber: "EMP_NO002",
			Descriptions:   []string{"desc bbb"},
		},
		{
			Dn:             "uid=ccc_user,dc=example,dc=com",
			UID:            "ccc_user",
			Email:          "ccc_user@example.com",
			EmployeeNumber: "EMP_NO003",
			Descriptions:   []string{"desc 0001", "desc 0002"}, // mod
		},
	}
	result, err = account.Diff(&testOld, &testNew)
	if err != nil {
		t.Errorf("account.Diff [Modify]exec fatal:%v", err)
	}
	if len(result[UpdateKey]) != 2 {
		t.Errorf("account.Diff [Modify]update list count wrong: %d", len(result[UpdateKey]))
	}
	for _, data := range result[UpdateKey] {
		if data.Dn != "uid=bbb_user,dc=example,dc=com" && data.Dn != "uid=ccc_user,dc=example,dc=com" {
			t.Errorf("account.Diff [Modify]update data wrong: %v", data)
		}
	}

	// Pattern: Create Delete
	testOld = append(testOld, Account{
		Dn:             "uid=ddd_user,dc=example,dc=com",
		UID:            "ddd_user",
		Email:          "ddd_user@example.com",
		EmployeeNumber: "EMP_NO004",
		Descriptions:   []string{"desc_ddd"},
	})
	testNew = append(testNew, Account{
		Dn:             "uid=eee_user,dc=example,dc=com",
		UID:            "eee_user",
		Email:          "eee_user@example.com",
		EmployeeNumber: "EMP_NO005",
		Descriptions:   []string{"desc_eee"},
	})
	result, err = account.Diff(&testOld, &testNew)
	if err != nil {
		t.Errorf("account.Diff [Create Delete]exec fatal:%v", err)
	}
	if len(result[UpdateKey]) != 2 {
		t.Errorf("account.Diff [Create Delete]update list count wrong: %d", len(result[UpdateKey]))
	}
	if len(result[CreateKey]) != 1 {
		t.Errorf("account.Diff [Create Delete]create list count wrong: %d", len(result[CreateKey]))
	}
	if result[CreateKey][0].Dn != "uid=eee_user,dc=example,dc=com" {
		t.Errorf("account.Diff [Create Delete]create data wrong: %v", result[CreateKey][0])
	}
	if len(result[DeleteKey]) != 1 {
		t.Errorf("account.Diff [Create Delete]delete list count wrong: %d", len(result[DeleteKey]))
	}
	if result[DeleteKey][0].Dn != "uid=ddd_user,dc=example,dc=com" {
		t.Errorf("account.Diff [Create Delete]delete data wrong: %v", result[DeleteKey][0])
	}

}
