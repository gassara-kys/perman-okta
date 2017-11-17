package main

import (
	"encoding/json"
	"os"
	"testing"
)

var tesUserProfile = UserProfile{
	LastName:    "Okta API Test",
	SecondEmail: "",
	MobilePhone: "",
	Email:       "test_okta_user@example.com",
	Login:       "test_okta_user@example.com",
	FirstName:   "Okta API Test",
}

var testGroupProfile = GroupProfile{
	Name:        "test_okta_group",
	Description: "test_okta_group",
}

func TestOkataClientAll(t *testing.T) {
	EnvLoad()
	// set Okta API
	var oktaClient = OktaClient{
		FQDN:   os.Getenv("OKTA_FQDN"),
		APIKEY: os.Getenv("OKTA_APIKEY"),
	}
	if oktaClient.FQDN == "" || oktaClient.APIKEY == "" {
		t.Error("set env okata client OKTA_FQDN and OKTA_APIKEY.")
		return
	}

	var jsonByte []byte
	var oktaErr error
	var oktaUser *OktaUser
	var oktaGroup *OktaGroup

	// Create User
	oktaUser, oktaErr = oktaClient.CreateUser(&tesUserProfile)
	if oktaErr != nil {
		t.Error(oktaErr)
	} else if oktaUser.LastName != tesUserProfile.LastName ||
		oktaUser.SecondEmail != tesUserProfile.SecondEmail ||
		oktaUser.MobilePhone != tesUserProfile.MobilePhone ||
		oktaUser.Email != tesUserProfile.Email ||
		oktaUser.Login != tesUserProfile.Login ||
		oktaUser.FirstName != tesUserProfile.FirstName {
		t.Errorf("CreateUser failed: %v", oktaUser)
	}
	jsonByte, _ = json.MarshalIndent(oktaUser, "", "  ")
	t.Logf("created user:\n%s", jsonByte)

	// Get User
	oktaUser, oktaErr = oktaClient.GetUserWithLogin(tesUserProfile.Login)
	if oktaErr != nil {
		t.Error(oktaErr)
	} else if oktaUser.LastName != tesUserProfile.LastName ||
		oktaUser.SecondEmail != tesUserProfile.SecondEmail ||
		oktaUser.MobilePhone != tesUserProfile.MobilePhone ||
		oktaUser.Email != tesUserProfile.Email ||
		oktaUser.Login != tesUserProfile.Login ||
		oktaUser.FirstName != tesUserProfile.FirstName {
		t.Errorf("GetUserWithLogin something wrong: %v", oktaUser)
	}
	jsonByte, _ = json.MarshalIndent(oktaUser, "", "  ")
	t.Logf("find user:\n%s", jsonByte)

	// Create Group
	oktaGroup, oktaErr = oktaClient.AddGroup(&testGroupProfile)
	if oktaErr != nil {
		t.Error(oktaErr)
	} else if oktaGroup.Name != testGroupProfile.Name ||
		oktaGroup.Description != testGroupProfile.Description {
		t.Errorf("AddGroup failed: %v", oktaGroup)
	}
	jsonByte, _ = json.MarshalIndent(oktaGroup, "", "  ")
	t.Logf("create group:\n%s", jsonByte)

	// Search Group
	oktaGroup, oktaErr = oktaClient.SearchGroups(testGroupProfile.Name)
	if oktaErr != nil {
		t.Error(oktaErr)
	} else if oktaGroup.Name != testGroupProfile.Name ||
		oktaGroup.Description != testGroupProfile.Description {
		t.Errorf("SearchGroups something wrong: %v", oktaGroup)
	}
	jsonByte, _ = json.MarshalIndent(oktaGroup, "", "  ")
	t.Logf("find group:\n%s", jsonByte)

	// Add member to Group
	oktaErr = oktaClient.AddUserToGroup(oktaGroup.ID, oktaUser.ID)
	if oktaErr != nil {
		t.Error(oktaErr)
	}

	// Remove member from Group
	oktaErr = oktaClient.RemoveUserFromGroup(oktaGroup.ID, oktaUser.ID)
	if oktaErr != nil {
		t.Error(oktaErr)
	}

	// Delete Group
	oktaErr = oktaClient.RemoveGroup(oktaGroup.ID)
	if oktaErr != nil {
		t.Error(oktaErr)
	}
	// Delete User
	oktaErr = oktaClient.DeleteUser(oktaUser.ID)
	if oktaErr != nil {
		t.Error(oktaErr)
	}
	return

}
