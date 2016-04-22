package internalValidation

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var url string
var authKey string

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

func TestInvalidHost(t *testing.T) {
	emailAddress := "test@turner.com"
	validation := NewSimpleIdentity("http://notaplace1356etyu.com/", authKey)
	isValid, err := validation.Validate(emailAddress)

	assert.Equal(t, false, isValid)
	assert.NotNil(t, err)
}

func TestInvalidStatus(t *testing.T) {
	emailAddress := "test@turner.com"
	validation := NewSimpleIdentity("http://cnn.com/", authKey)
	isValid, err := validation.Validate(emailAddress)

	assert.Equal(t, false, isValid)
	assert.Equal(t, errors.New("Invalid Status Code: 404 Not Found"), err)
}

func TestValidUser(t *testing.T) {
	emailAddress := "don.browning@turner.com"
	validation := NewSimpleIdentity(url, authKey)
	isValid, err := validation.Validate(emailAddress)

	assert.Equal(t, true, isValid)
	assert.Nil(t, err)
}

func TestInvalidUser(t *testing.T) {
	emailAddress := "test@turner.com"
	validation := NewSimpleIdentity(url, authKey)
	isValid, err := validation.Validate(emailAddress)

	assert.Equal(t, false, isValid)
	assert.Equal(t, errors.New("No records found"), err)
}
