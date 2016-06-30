package aws

import (
	"fmt"
	"testing"
)

const checkid = "eW7HH0l7J9"

func TestServiceLevels(t *testing.T) {
	supportaws := NewServiceLevels(region, nil, nil, nil)

	result, err := supportaws.ListServiceLevels(checkid)
	fmt.Println(err)
	fmt.Println(result)
}
