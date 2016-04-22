package account

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/smithatlanta/cloudthings/nonprovider/identitymgmt"
	"github.com/stretchr/testify/assert"
)

var url string
var authKey string
var accountID string

func CToGoString(c []byte) string {
	n := -1
	for i, b := range c {
		if b == 0 {
			break
		}
		n = i
	}
	return string(c[:n+1])
}

func TestMain(m *testing.M) {
	flag.StringVar(&url, "url", "", "url for authorization")
	flag.StringVar(&authKey, "authKey", "", "authorization key")
	flag.Parse()

	if len(url) == 0 || len(authKey) == 0 {
		fmt.Println("Missing url or authKey parameter.")
		os.Exit(0)
	}

	os.Exit(m.Run())
}

func TestMarshall(t *testing.T) {
	var test BasicCloudhealthAccount
	var x1 identitymgmt.BasicAccount
	x1.ID = 1
	x1.Name = "Name1"
	test.AWSAccounts = append(test.AWSAccounts, x1)
	var x2 identitymgmt.BasicAccount
	x2.ID = 2
	x2.Name = "Name1"
	test.AWSAccounts = append(test.AWSAccounts, x2)

	b, err := json.Marshal(test)
	assert.JSONEq(t, `{"aws_accounts":[{"id":1,"name":"Name1"},{"id":2,"name":"Name1"}]}`, CToGoString(b[:]))
	assert.Nil(t, err)
}

func TestUnMarshall(t *testing.T) {
	file, e := ioutil.ReadFile("./test.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	var chAccount BasicCloudhealthAccount
	json.Unmarshal(file, &chAccount)

	assert.Equal(t, 1, chAccount.AWSAccounts[0].ID)
	assert.Equal(t, "Name1", chAccount.AWSAccounts[0].Name)
}

func TestGetAccounts(t *testing.T) {
	account := NewSimpleAccount(url, authKey)
	accounts, err := account.GetAccounts()

	lenAccounts := len(accounts)

	fmt.Println(accounts)
	assert.NotEqual(t, 0, lenAccounts)
	assert.Nil(t, err)
}

func TestGetAccount(t *testing.T) {

	accountID := "343597384152"
	account := NewSimpleAccount(url, authKey)
	chaccount, err := account.GetAccount(accountID)

	assert.NotNil(t, chaccount)
	assert.Nil(t, err)
}
