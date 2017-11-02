package main

import (
	"fmt"
	"log"

	"gopkg.in/ldap.v2"
)

// LdapClient LDAPクライアント
type LdapClient struct {
	Host       string
	BaseDn     string
	Filter     string
	Attributes []string
	SizeLimit  int
	TimeLimit  int
	TypeOnly   bool
}

// Search ldapsearch
func (l LdapClient) Search() (result *ldap.SearchResult, err error) {

	ldapConn, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", l.Host, 389))
	if err != nil {
		log.Printf("connerction Error... err: %+v", err)
		return
	}
	defer ldapConn.Close()

	// ldapsearch
	searchRequest := ldap.NewSearchRequest(
		l.BaseDn,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		l.SizeLimit,
		l.TimeLimit,
		l.TypeOnly,
		l.Filter,
		[]string{"dn", "uid", "email", "employeeNumber", "description"},
		nil,
	)
	result, err = ldapConn.Search(searchRequest)
	if err != nil {
		log.Printf("ldap search Error... err: %+v", err)
		return
	}
	return
}
