package main

import (
	"testing"

	ldap "gopkg.in/ldap.v2"
)

const testFileNm = "tmp/test.json"

var testAccounts = []LdapAccount{
	{Dn: "uid=aaa_user,dc=example,dc=com", UID: "aaa_user", Email: "aaa_user@example.com", EmployeeNumber: "EMP_NO001", Descriptions: []string{"desc_aaa"}},
	{Dn: "uid=bbb_user,dc=example,dc=com", UID: "bbb_user", Email: "bbb_user@example.com", EmployeeNumber: "EMP_NO002", Descriptions: []string{"desc bbb"}},
	{Dn: "uid=ccc_user,dc=example,dc=com", UID: "ccc_user", Email: "ccc_user@example.com", EmployeeNumber: "EMP_NO003", Descriptions: []string{"desc ccc"}},
}

func TestJSON(t *testing.T) {

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
	var ldapAccount = LdapAccount{}

	// OutJSON
	err := ldapAccount.OutJSON(testFileNm, entries)
	if err != nil {
		t.Errorf("ldapAccount.OutJSON exec failed: file: %s, data: %v", testFileNm, entries)
	}

	// LoadJSON
	accouts, err := ldapAccount.LoadJSON(testFileNm)
	if err != nil {
		t.Errorf("ldapAccount.LoadJSON exec failed: file: %s", testFileNm)
	}

	for idx, data := range *accouts {
		if testAccounts[idx].Dn != data.Dn ||
			testAccounts[idx].UID != data.UID ||
			testAccounts[idx].Email != data.Email ||
			testAccounts[idx].EmployeeNumber != data.EmployeeNumber {
			t.Errorf("JSON data not match: testAccount: %s, jsonAccout: %s", testAccounts[idx].Dn, data.Dn)
		}
	}
}
