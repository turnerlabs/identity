package internalValidation

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/smithatlanta/cloudthings/nonprovider/identitymgmt"
)

//NewSimpleIdentity -
func NewSimpleIdentity(url string, authKey string) identitymgmt.SimpleIdentity {
	if len(url) == 0 || len(authKey) == 0 {
		return nil
	}

	store := newValidation(url, authKey)
	return &store
}

func newValidation(url string, authKey string) validation {
	s := validation{
		url:     url,
		authKey: authKey,
	}
	return s
}

type validation struct {
	url     string
	authKey string
}

//Account -
type Account struct {
	AccountEnabled bool `json:"accountEnabled"`
}

//Validate -
func (s *validation) Validate(user string) (bool, error) {
	if len(user) == 0 {
		return false, errors.New("User was not passed in")
	}

	fullURL := s.url + user
	client := &http.Client{}
	req, err := http.NewRequest("GET", fullURL, nil)
	req.Header.Add("Authorization", s.authKey)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return false, err
	}

	defer resp.Body.Close()

	status := resp.StatusCode
	if status != 200 {
		return false, errors.New("Invalid Status Code: " + resp.Status)
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, errors.New("Error in ReadAll")
	}

	var accounts []Account
	json.Unmarshal(responseBody, &accounts)

	if len(accounts) == 0 {
		return false, errors.New("No records found")
	}

	return accounts[0].AccountEnabled, nil
}
