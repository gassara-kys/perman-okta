package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
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
	Profile         struct {
		LastName    string      `json:"lastName"`
		SecondEmail interface{} `json:"secondEmail"`
		MobilePhone interface{} `json:"mobilePhone"`
		Email       string      `json:"email"`
		Login       string      `json:"login"`
		FirstName   string      `json:"firstName"`
	} `json:"profile"`
	Credentials struct {
		Password struct {
		} `json:"password"`
		RecoveryQuestion struct {
			Question string `json:"question"`
		} `json:"recovery_question"`
		Provider struct {
			Type string `json:"type"`
			Name string `json:"name"`
		} `json:"provider"`
	} `json:"credentials"`
	Links struct {
		Suspend struct {
			Href   string `json:"href"`
			Method string `json:"method"`
		} `json:"suspend"`
		ResetPassword struct {
			Href   string `json:"href"`
			Method string `json:"method"`
		} `json:"resetPassword"`
		ExpirePassword struct {
			Href   string `json:"href"`
			Method string `json:"method"`
		} `json:"expirePassword"`
		ForgotPassword struct {
			Href   string `json:"href"`
			Method string `json:"method"`
		} `json:"forgotPassword"`
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		ChangeRecoveryQuestion struct {
			Href   string `json:"href"`
			Method string `json:"method"`
		} `json:"changeRecoveryQuestion"`
		Deactivate struct {
			Href   string `json:"href"`
			Method string `json:"method"`
		} `json:"deactivate"`
		ChangePassword struct {
			Href   string `json:"href"`
			Method string `json:"method"`
		} `json:"changePassword"`
	} `json:"_links"`
}

// GetUser
func (okta OktaClient) GetUser(q string) (*[]OktaUser, error) {

	req, _ := http.NewRequest("GET", "https://"+okta.FQDN+"/api/v1/users?q="+q, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "SSWS "+okta.APIKEY)
	dump, _ := httputil.DumpRequestOut(req, true)
	log.Printf("%s", dump)

	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	} else if res.StatusCode == 404 {
		return &[]OktaUser{}, nil
	} else if res.StatusCode != 200 {
		return nil, fmt.Errorf("Unable to get this url : http status %d", res.StatusCode)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	// 取得したjsonを構造体へデコード
	var oktaUser []OktaUser
	if err := json.Unmarshal(body, &oktaUser); err != nil {
		return nil, err
	}

	return &oktaUser, nil
}
