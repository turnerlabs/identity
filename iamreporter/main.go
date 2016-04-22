package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/smithatlanta/cloudthings/nonprovider/identitymgmt/aws"
	"github.com/smithatlanta/cloudthings/nonprovider/identitymgmt/internalValidation"
)

var url string
var authKey string
var region string
var account string
var role string

func init() {
	flag.StringVar(&url, "url", "", "url is required(authorization endpoint)")
	flag.StringVar(&authKey, "authKey", "", "authKey is required(authorization key on header)")
	flag.StringVar(&region, "region", "", "region is required(aws region)")
	flag.StringVar(&account, "account", "", "account to scan")
	flag.StringVar(&role, "role", "", "role name to use")
}

func main() {
	flag.Parse()

	if len(url) == 0 || len(authKey) == 0 || len(region) == 0 {
		flag.PrintDefaults()
		os.Exit(0)
	}

	//accounts to check
	accounts := []string{account}

	//Loop thru the accounts
	for i := 0; i < len(accounts); i++ {
		var accessKey *string
		var secretKey *string
		var sessionToken *string
		var err error
		//Assume Role in each account
		if len(accounts[i]) > 0 {
			sts := aws.NewRole(region, nil, nil, nil)
			accessKey, secretKey, sessionToken, err = sts.AssumeRole("arn:aws:iam::"+accounts[i]+":role/"+role, "IAMCHECKER")
			if err != nil {
				//fmt.Println(err.Error())
				os.Exit(0)
			}
		}

		//AWS user scan showing name, groups
		iam := aws.NewIdentity(region, accessKey, secretKey, sessionToken)
		users, err := iam.ListUsers()
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("-------------------------------------------------------------------------------")
		fmt.Println("---Users NOT valid based on internal checks for account: " + accounts[i] + "---")
		fmt.Println("-------------------------------------------------------------------------------")
		for a := 0; a < len(users); a++ {
			internal := internalValidation.NewSimpleIdentity(url, authKey)
			isValid, err := internal.Validate(users[a])
			if err != nil && err.Error() != "No records found" {
				continue
			}
			if isValid == false {
				groups, err := iam.ListGroupsForUser(users[a])
				if err != nil {
					fmt.Println(err.Error())
				}
				var inGroups string
				if groups != nil {
					for b := 0; b < len(groups); b++ {
						inGroups += groups[b]
						if b < (len(groups) - 1) {
							inGroups += ", "
						}
					}
				}
				if len(inGroups) > 0 {
					fmt.Println("User: " + users[a] + " Groups: " + inGroups)
				} else {
					fmt.Println("User: " + users[a])
				}
			}
		}
	}
}
