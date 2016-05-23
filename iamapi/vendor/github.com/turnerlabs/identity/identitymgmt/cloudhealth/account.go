package account

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/turnerlabs/identity/identitymgmt"
)

var allAccounts []identitymgmt.BasicAccount

//BasicCloudhealthAccount -
type BasicCloudhealthAccount struct {
	AWSAccounts []identitymgmt.BasicAccount `json:"aws_accounts"`
}

//NewSimpleAccount -
func NewSimpleAccount(url string, authKey string) identitymgmt.SimpleAccount {
	if len(url) == 0 || len(authKey) == 0 {
		return nil
	}

	store := newchAccount(url, authKey)
	return &store
}

func newchAccount(url string, authKey string) chAccount {
	s := chAccount{
		url:     url,
		authKey: authKey,
	}
	return s
}

type chAccount struct {
	url     string
	authKey string
}

//GetAccounts -
func (s *chAccount) GetAccounts() ([]identitymgmt.BasicAccount, error) {
	allAccounts := []identitymgmt.BasicAccount{}
	var chaccounts BasicCloudhealthAccount
	pg := 1

	for {
		fullURL := s.url + "/v1/aws_accounts?api_key=" + s.authKey + "&page=" + strconv.Itoa(pg)
		client := &http.Client{}
		req, err := http.NewRequest("GET", fullURL, nil)
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err.Error())
			return chaccounts.AWSAccounts, err
		}

		defer resp.Body.Close()

		status := resp.StatusCode
		if status != 200 {
			return chaccounts.AWSAccounts, errors.New("Invalid Status Code: " + resp.Status)
		}

		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return chaccounts.AWSAccounts, errors.New("Error in ReadAll")
		}

		json.Unmarshal(responseBody, &chaccounts)

		accoutLen := len(chaccounts.AWSAccounts)

		if accoutLen == 0 {
			break
		} else {
			for ctr := 0; ctr < accoutLen; ctr++ {
				allAccounts = append(allAccounts, chaccounts.AWSAccounts[ctr])
			}
			pg++
		}
	}

	return allAccounts, nil
}

//GetAccount -
func (s *chAccount) GetAccount(id string) (identitymgmt.CloudhealthAccount, error) {
	var chaccounts identitymgmt.CloudhealthAccount

	fullURL := s.url + "/v1/aws_accounts/" + string(id) + "?api_key=" + s.authKey
	client := &http.Client{}
	req, err := http.NewRequest("GET", fullURL, nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return chaccounts, err
	}

	defer resp.Body.Close()

	status := resp.StatusCode
	if status != 200 {
		return chaccounts, errors.New("Invalid Status Code: " + resp.Status)
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return chaccounts, errors.New("Error in ReadAll")
	}

	json.Unmarshal(responseBody, &chaccounts)

	// if chaccounts == nil {
	// 	return chaccounts, errors.New("No records found")
	// }

	return chaccounts, nil
}
