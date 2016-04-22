package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const roleArn = ""
const roleSessionName = ""

func TestAssumeRoleFail(t *testing.T) {
	stsaws := NewRole(region, nil, nil, nil)
	ret1, ret2, ret3, err := stsaws.AssumeRole(roleArn, roleSessionName)

	assert.Nil(t, ret1)
	assert.Nil(t, ret2)
	assert.Nil(t, ret3)
	assert.NotNil(t, err)
}
