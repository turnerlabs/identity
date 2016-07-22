package membership

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/turnerlabs/identity/identitymgmt"
)

//NewSimpleAccount -
func NewSimpleAccount(url string, authKey string) identitymgmt.Membership {
	if len(url) == 0 || len(authKey) == 0 {
		return nil
	}

	store := newGithubAccount(url, authKey)
	return &store
}

func newGithubAccount(url string, authKey string) githubAccount {
	s := githubAccount{
		url:     url,
		authKey: authKey,
	}
	return s
}

type githubAccount struct {
	url     string
	authKey string
}

func (s *githubAccount) Members(org string) ([]identitymgmt.GithubMember, error) {
	allMembers := []identitymgmt.GithubMember{}
	fullURL := s.url + "/orgs/" + org + "/members"
	client := &http.Client{}
	req, err := http.NewRequest("GET", fullURL, nil)
	token := "token " + s.authKey
	req.Header.Set("Authorization", token)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return allMembers, err
	}

	defer resp.Body.Close()

	status := resp.StatusCode
	if status != 200 {
		return allMembers, errors.New("Invalid Status Code: " + resp.Status)
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return allMembers, errors.New("Error in ReadAll")
	}

	json.Unmarshal(responseBody, &allMembers)

	// if chaccounts == nil {
	// 	return chaccounts, errors.New("No records found")
	// }

	return allMembers, nil
}
