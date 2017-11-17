package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// OktaClient OktaAPI client
type OktaClient struct {
	FQDN   string
	APIKEY string
}

// OktaUser Response
type OktaUser struct {
	ID              string    `json:"id"`
	Status          string    `json:"status"`
	Created         time.Time `json:"created"`
	Activated       time.Time `json:"activated"`
	StatusChanged   time.Time `json:"statusChanged"`
	LastLogin       time.Time `json:"lastLogin"`
	LastUpdated     time.Time `json:"lastUpdated"`
	PasswordChanged time.Time `json:"passwordChanged"`
	UserProfile     `json:"profile"`
}

// UserProfile OktaUser Profile
type UserProfile struct {
	LastName    string      `json:"lastName"`
	SecondEmail interface{} `json:"secondEmail"`
	MobilePhone interface{} `json:"mobilePhone"`
	Email       string      `json:"email"`
	Login       string      `json:"login"` // must email type format
	FirstName   string      `json:"firstName"`
}

// OktaGroup Response
type OktaGroup struct {
	ID                    string    `json:"id"`
	Created               time.Time `json:"created"`
	LastUpdated           time.Time `json:"lastUpdated"`
	LastMembershipUpdated time.Time `json:"lastMembershipUpdated"`
	Type                  string    `json:"type"`
	ObjectClass           []string  `json:"objectClass"`
	GroupProfile          `json:"profile"`
}

// GroupProfile OktaGroup Profile
type GroupProfile struct {
	Name        string      `json:"name"`
	Description interface{} `json:"description"`
}

// CreateUserRequest request body for create
type CreateUserRequest struct {
	UserProfile `json:"profile"`
}

// CreateGroupRequest request body for create
type CreateGroupRequest struct {
	GroupProfile `json:"profile"`
}

// GetUserWithLogin Get User with Login API
func (okta OktaClient) GetUserWithLogin(login string) (*OktaUser, error) {

	req, _ := http.NewRequest("GET", "https://"+okta.FQDN+"/api/v1/users/"+login, nil)
	okta.setHeader(req)

	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	} else if res.StatusCode == http.StatusNotFound {
		log.Printf("Not Found user: http status %d", res.StatusCode)
		return &OktaUser{}, nil
	} else if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unable to get this url :http status %d", res.StatusCode)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	// 取得したjsonを構造体へデコード
	var oktaUser OktaUser
	if err := json.Unmarshal(body, &oktaUser); err != nil {
		return nil, err
	}

	return &oktaUser, nil
}

// CreateUser Create Activated User without Credentials
func (okta OktaClient) CreateUser(profile *UserProfile) (*OktaUser, error) {

	createReq := CreateUserRequest{}
	createReq.UserProfile = *profile
	jsonBytes, err := json.Marshal(createReq)
	if err != nil {
		return nil, err
	}
	req, _ := http.NewRequest(
		"POST",
		"https://"+okta.FQDN+"/api/v1/users?activate=true",
		bytes.NewBuffer(jsonBytes),
	)
	okta.setHeader(req)

	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"Could not create user: http status %d: body %s ",
			res.StatusCode,
			body,
		)
	}
	// 取得したjsonを構造体へデコード
	oktaUser := OktaUser{}
	if err := json.Unmarshal(body, &oktaUser); err != nil {
		return nil, err
	}

	return &oktaUser, nil
}

// DeleteUser Delete User API
func (okta OktaClient) DeleteUser(id string) error {

	// deactivate user
	req, _ := http.NewRequest("POST", "https://"+okta.FQDN+"/api/v1/users/"+id+"/lifecycle/deactivate", nil)
	okta.setHeader(req)

	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return err
	} else if res.StatusCode == http.StatusNotFound {
		log.Printf("Not Found user: http status %d: user id %s ", res.StatusCode, id)
		return nil
	} else if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Could not deactivate user :http status %d", res.StatusCode)
	}
	log.Printf("Deactivated user: %s", id)

	// delete user
	req, _ = http.NewRequest("DELETE", "https://"+okta.FQDN+"/api/v1/users/"+id, nil)
	okta.setHeader(req)

	client = new(http.Client)
	res, err = client.Do(req)
	if err != nil {
		return err
	} else if res.StatusCode == http.StatusNotFound {
		log.Printf("Not Found user: http status %d: user id %s ", res.StatusCode, id)
		return nil
	} else if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Could not delete user : http status %d", res.StatusCode)
	}
	log.Printf("Deleted user: %s", id)

	return nil
}

// SearchGroups Search Groups API
func (okta OktaClient) SearchGroups(name string) (*OktaGroup, error) {

	req, _ := http.NewRequest("GET", "https://"+okta.FQDN+"/api/v1/groups?q="+name, nil)
	okta.setHeader(req)

	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	} else if res.StatusCode == http.StatusNotFound {
		log.Printf("Not Found Group: http status %d: Group is %s ", res.StatusCode, name)
		return &OktaGroup{}, nil
	} else if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unable to get this url :http status %d", res.StatusCode)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	// 取得したjsonを構造体へデコード
	var oktaGroups []OktaGroup
	if err := json.Unmarshal(body, &oktaGroups); err != nil {
		return nil, err
	}
	for _, oktaGroup := range oktaGroups {
		// Check Exactly Match Group name
		if oktaGroup.GroupProfile.Name == name {
			return &oktaGroup, nil
		}
	}

	return &OktaGroup{}, nil

}

// AddGroup  Add Group API
func (okta OktaClient) AddGroup(profile *GroupProfile) (*OktaGroup, error) {
	createReq := CreateGroupRequest{}
	createReq.GroupProfile = *profile
	jsonBytes, err := json.Marshal(createReq)
	if err != nil {
		return nil, err
	}
	req, _ := http.NewRequest("POST", "https://"+okta.FQDN+"/api/v1/groups", bytes.NewBuffer(jsonBytes))
	okta.setHeader(req)
	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"Could not create group: http status %d: body %s ",
			res.StatusCode,
			body,
		)
	}
	// 取得したjsonを構造体へデコード
	oktaGroup := OktaGroup{}
	if err := json.Unmarshal(body, &oktaGroup); err != nil {
		return nil, err
	}

	return &oktaGroup, nil
}

// RemoveGroup Call Remove Group API (only OKTA_GROUP type)
func (okta OktaClient) RemoveGroup(id string) error {

	req, _ := http.NewRequest("DELETE", "https://"+okta.FQDN+"/api/v1/groups/"+id, nil)
	okta.setHeader(req)

	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return err
	} else if res.StatusCode == http.StatusNotFound {
		log.Printf("Not Found Group: http status %d: id %s ", res.StatusCode, id)
		return nil
	} else if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Could not delete Group : http status %d", res.StatusCode)
	}
	log.Printf("Deleted Group: %s", id)

	return nil
}

// AddUserToGroup Call Add User to Group API (only OKTA_GROUP type)
func (okta OktaClient) AddUserToGroup(gid, uid string) error {

	req, _ := http.NewRequest("PUT", "https://"+okta.FQDN+"/api/v1/groups/"+gid+"/users/"+uid, nil)
	okta.setHeader(req)

	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return err
	} else if res.StatusCode == http.StatusNotFound {
		log.Printf("Not Found Group or User: http status %d: gid=%s uid=%s", res.StatusCode, gid, uid)
	} else if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Could not Add user to Group : http status %d", res.StatusCode)
	}
	log.Printf("Add User(%s) to Group(%s)", uid, gid)

	return nil
}

// RemoveUserFromGroup Call Remove User From Group API (only OKTA_GROUP type group)
func (okta OktaClient) RemoveUserFromGroup(gid, uid string) error {

	req, _ := http.NewRequest("DELETE", "https://"+okta.FQDN+"/api/v1/groups/"+gid+"/users/"+uid, nil)
	okta.setHeader(req)

	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return err
	} else if res.StatusCode == http.StatusNotFound {
		log.Printf("Not Found Group or User: http status %d: gid=%s uid=%s", res.StatusCode, gid, uid)
	} else if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Could not Remove user to Group : http status %d", res.StatusCode)
	}
	log.Printf("Remove User(%s) to Group(%s)", uid, gid)

	return nil
}

// set Common HTTP Header
func (okta OktaClient) setHeader(req *http.Request) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "SSWS "+okta.APIKEY)

}
