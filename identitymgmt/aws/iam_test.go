package aws

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const region = "us-east-1"
const group1 = "test1"
const group2 = "test2"
const user = "test1@test.com"

func TestCreateGroup1(t *testing.T) {
	iamaws := NewIdentity(region, nil, nil, nil)
	err := iamaws.CreateGroup(group1)
	if err != nil {
		fmt.Println(err.Error())
	}

	assert.Nil(t, err)
}

func TestCreateGroup2(t *testing.T) {
	iamaws := NewIdentity(region, nil, nil, nil)
	err := iamaws.CreateGroup(group2)
	if err != nil {
		fmt.Println(err.Error())
	}

	assert.Nil(t, err)
}

func TestCreateUser(t *testing.T) {
	iamaws := NewIdentity(region, nil, nil, nil)
	err := iamaws.CreateUser(user)
	if err != nil {
		fmt.Println(err.Error())
	}

	assert.Nil(t, err)
}

func TestListUsers(t *testing.T) {
	iamaws := NewIdentity(region, nil, nil, nil)
	resp, err := iamaws.ListUsers()
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(resp)
	assert.Nil(t, err)
	assert.NotEmpty(t, len(resp))
}

func TestEmailableUsers(t *testing.T) {
	iamaws := NewIdentity(region, nil, nil, nil)
	resp, err := iamaws.EmailableUsers()
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(resp)
	assert.Nil(t, err)
	assert.NotEmpty(t, len(resp))
}

func TestCreateAccessKey1(t *testing.T) {
	iamaws := NewIdentity(region, nil, nil, nil)
	err := iamaws.CreateAccessKey(user)
	if err != nil {
		fmt.Println(err.Error())
	}

	assert.Nil(t, err)
}

func TestCreateAccessKey2(t *testing.T) {
	iamaws := NewIdentity(region, nil, nil, nil)
	err := iamaws.CreateAccessKey(user)
	if err != nil {
		fmt.Println(err.Error())
	}

	assert.Nil(t, err)
}

func TestListAccessKeys(t *testing.T) {
	iamaws := NewIdentity(region, nil, nil, nil)
	resp, err := iamaws.ListAccessKeys(user)
	if err != nil {
		fmt.Println(err.Error())
	}

	assert.Nil(t, err)
	assert.Equal(t, 2, len(resp))
}

func TestAddUserToGroup1(t *testing.T) {
	iamaws := NewIdentity(region, nil, nil, nil)
	err := iamaws.AddUserToGroup(user, group1)
	if err != nil {
		fmt.Println(err.Error())
	}

	assert.Nil(t, err)
}

func TestAddUserToGroup2(t *testing.T) {
	iamaws := NewIdentity(region, nil, nil, nil)
	err := iamaws.AddUserToGroup(user, group2)
	if err != nil {
		fmt.Println(err.Error())
	}

	assert.Nil(t, err)
}

func TestListGroupsForUser(t *testing.T) {
	iamaws := NewIdentity(region, nil, nil, nil)
	resp, err := iamaws.ListGroupsForUser(user)
	if err != nil {
		fmt.Println(err.Error())
	}

	assert.Nil(t, err)
	assert.Equal(t, 2, len(resp))
}

func TestRemoveUserFromGroups(t *testing.T) {
	iamaws := NewIdentity(region, nil, nil, nil)
	resp, err := iamaws.ListGroupsForUser(user)
	if err != nil {
		fmt.Println(err.Error())
	}

	assert.Nil(t, err)
	assert.Equal(t, 2, len(resp))

	for i := 0; i < len(resp); i++ {
		err := iamaws.RemoveUserFromGroup(user, resp[i])
		if err != nil {
			fmt.Println(err.Error())
		}

		assert.Nil(t, err)
	}
}

func TestRemoveAccessKeysFromUser(t *testing.T) {
	iamaws := NewIdentity(region, nil, nil, nil)
	resp, err := iamaws.ListAccessKeys(user)
	if err != nil {
		fmt.Println(err.Error())
	}

	assert.Nil(t, err)
	assert.Equal(t, 2, len(resp))

	for i := 0; i < len(resp); i++ {
		err := iamaws.DeleteAccessKey(user, resp[i])
		if err != nil {
			fmt.Println(err.Error())
		}
		assert.Nil(t, err)
	}
}

func TestDeleteUser(t *testing.T) {
	iamaws := NewIdentity(region, nil, nil, nil)
	err := iamaws.DeleteUser(user)
	if err != nil {
		fmt.Println(err.Error())
	}

	assert.Nil(t, err)
}

func TestDeleteGroup1(t *testing.T) {
	iamaws := NewIdentity(region, nil, nil, nil)
	err := iamaws.DeleteGroup(group1)
	if err != nil {
		fmt.Println(err.Error())
	}

	assert.Nil(t, err)
}

func TestDeleteGroup2(t *testing.T) {
	iamaws := NewIdentity(region, nil, nil, nil)
	err := iamaws.DeleteGroup(group2)
	if err != nil {
		fmt.Println(err.Error())
	}

	assert.Nil(t, err)
}
