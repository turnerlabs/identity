package internalAccount

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
)

func createAccount(accountID string) Account {

	var accountItems []AccountItem

	var newAccountItem = AccountItem{
		Key:   "TestKey",
		Value: "TestValue",
	}

	accountItems = append(accountItems, newAccountItem)

	var newAccount = &Account{
		AccountID:    accountID,
		AccountItems: accountItems,
	}

	return *newAccount
}

func TestAccount(t *testing.T) {
	generatedAccountID := uuid.New()
	account := createAccount(generatedAccountID)

	store := NewAccountStore("us-east-1", true)
	err := store.Add(account)
	assert.Nil(t, err)

	acc, err := store.GetByAccountID(generatedAccountID)
	assert.Nil(t, err)
	assert.EqualValues(t, acc.AccountID, generatedAccountID, "Accounts are not equal")

	err = store.Delete(generatedAccountID)
	assert.Nil(t, err)
}

// test for a non existent accountID
func TestNicknameGetByAccountIDMissingRow(t *testing.T) {
	var accountID = "12348"

	store := NewAccountStore("us-east-1", true)
	acc, err := store.GetByAccountID(accountID)
	assert.Nil(t, acc)
	assert.Error(t, err)
	assert.Equal(t, err, ErrNoRowFound)
}

func TestJSONMarshall(t *testing.T) {
	generatedAccountID := uuid.New()
	account := createAccount(generatedAccountID)

	st, err := json.Marshal(account)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(st))

}

func TestJSONUnMarshall(t *testing.T) {
	accJSON := []byte(`{"accountID":"123456789","accountItems":[{"key":"test1","value":"test2"}]}`)
	var acc Account
	err := json.Unmarshal(accJSON, &acc)
	fmt.Println(acc.AccountID)
	if err != nil {
		fmt.Println("error")
		fmt.Println(err)
	}

	fmt.Println(acc)

}
