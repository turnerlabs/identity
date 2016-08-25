package main

import (
	"fmt"
	"os"

	"github.com/turnerlabs/identity/identitymgmt/aws"
	"github.com/turnerlabs/identity/identitymgmt/internalValidation"
)

// var intURL string
// var intAuthKey string
//
// func init() {
// 	flag.StringVar(&intURL, "intURL", "", "url is required(authorization endpoint)")
// 	flag.StringVar(&intAuthKey, "intauthKey", "", "authKey is required(authorization key on header)")
// }

func main() {
	//	flag.Parse()

	var accessKey *string
	var secretKey *string
	var sessionToken *string
	var region string

	region = "us-east-1"

	//AWS user scan showing name, groups
	iam := aws.NewIdentity(region, accessKey, secretKey, sessionToken)

	switch os.Args[1] {
	case "-e":
		users, err := iam.EmailableUsers()
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(users)
	case "-a":
		users, err := iam.ListUsers()
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(users)
	case "-c":
		users, err := iam.ListUsers()
		if err != nil {
			fmt.Println(err.Error())
		}
		var userList []string
		var intURL string
		var intAuthKey string

		for a := 0; a < len(users); a++ {
			internal := internalValidation.NewSimpleIdentity(intURL, intAuthKey)
			isValid, err := internal.Validate(users[a])
			if err != nil && err.Error() != "No records found" {
				continue
			}
			if isValid == false {
				userList = append(userList, users[a])
			}
		}

	}
}
