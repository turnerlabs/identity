package membership

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var url string
var authKey string
var org string

func TestMain(m *testing.M) {
	flag.StringVar(&url, "url", "", "url for github")
	flag.StringVar(&authKey, "authKey", "", "authorization key")
	flag.StringVar(&org, "org", "", "organization")
	flag.Parse()

	if len(url) == 0 {
		fmt.Println("Missing url parameter.")
		os.Exit(0)
	}

	if len(authKey) == 0 {
		fmt.Println("Missing authKey parameter.")
		os.Exit(0)
	}

	if len(org) == 0 {
		fmt.Println("Missing org parameter.")
		os.Exit(0)
	}

	os.Exit(m.Run())
}

func TestMembers(t *testing.T) {
	membership := NewSimpleAccount(url, authKey)
	members, err := membership.Members(org)

	lenAccounts := len(members)
	assert.NotEqual(t, 0, lenAccounts)
	assert.Nil(t, err)
}
