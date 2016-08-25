package main

import (
	"flag"
	"fmt"

	"github.com/turnerlabs/identity/identitymgmt/aws"
	"github.com/turnerlabs/identity/identitymgmt/internalValidation"
)

var url string
var authKey string
var boolA *bool
var boolC *bool
var boolE *bool

func init() {
	boolA = flag.Bool("a", false, "to get all users based on the credentials key")
	boolC = flag.Bool("c", false, "to get all the non active users based on the credentials key")
	boolE = flag.Bool("e", false, "to get all the emailable users based on the credentials key")

	flag.StringVar(&url, "url", "", "url is required(authorization endpoint)")
	flag.StringVar(&authKey, "authKey", "", "authKey is required(authorization key on header)")
}

func main() {
	flag.Parse()

	var accessKey *string
	var secretKey *string
	var sessionToken *string
	var region string

	region = "us-east-1"

	iam := aws.NewIdentity(region, accessKey, secretKey, sessionToken)

	if *boolE == true {
		users, err := iam.EmailableUsers()
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(users)
	}

	if *boolA == true {
		users, err := iam.ListUsers()
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(users)
	}

	if *boolC == true {
		if url == "" {
			fmt.Println("intURL - url is required(authorization endpoint)")
			return
		}

		if authKey == "" {
			fmt.Println("intAuthKey - authKey is required(authorization key on header)")
			return
		}

		users, err := iam.ListUsers()
		if err != nil {
			fmt.Println(err.Error())
		}

		for a := 0; a < len(users); a++ {
			internal := internalValidation.NewSimpleIdentity(url, authKey)
			isValid, err := internal.Validate(users[a])
			if err != nil && err.Error() != "No records found" {
				continue
			}
			if isValid == false {
				fmt.Println(users[a])
			}
		}
	}
}
