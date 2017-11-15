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
	Profile         `json:"profile"`
}

// Profile OktaUser Profile
type Profile struct {
	LastName    string      `json:"lastName"`
	SecondEmail interface{} `json:"secondEmail"`
	MobilePhone interface{} `json:"mobilePhone"`
	Email       string      `json:"email"`
	Login       string      `json:"login"` // must email type format
	FirstName   string      `json:"firstName"`
}

// CreateUserRequest request body for create
type CreateUserRequest struct {
	Profile `json:"profile"`
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
		return nil, fmt.Errorf("Unable to get this url : http status %d", res.StatusCode)
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
func (okta OktaClient) CreateUser(profile *Profile) (*OktaUser, error) {

	createReq := CreateUserRequest{}
	createReq.Profile = *profile
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
			"Could not create user: http status %d: http body :%s ",
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

// set Common HTTP Header
func (okta OktaClient) setHeader(req *http.Request) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "SSWS "+okta.APIKEY)

}
