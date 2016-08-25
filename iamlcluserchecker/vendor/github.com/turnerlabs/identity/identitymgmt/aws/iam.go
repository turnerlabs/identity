package aws

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"

	"github.com/turnerlabs/identity/identitymgmt"
)

//Initilaization

// NewIdentity -
func NewIdentity(region string, accessKey *string, secretKey *string, token *string) identitymgmt.Identity {
	store := newIAMAWS(region, accessKey, secretKey, token)
	return &store
}

func newIAMAWS(region string, accessKey *string, secretKey *string, token *string) iamaws {
	var awsConfig *aws.Config
	if accessKey == nil || secretKey == nil || token == nil {
		awsConfig = aws.NewConfig().WithRegion(region)
	} else {
		awsConfig = aws.NewConfig().WithRegion(region)
		awsConfig.WithCredentials(credentials.NewStaticCredentials(*accessKey, *secretKey, *token))
	}

	awsConfig.WithCredentialsChainVerboseErrors(true)
	session := session.New(awsConfig)

	s := iamaws{
		svc: iam.New(session),
	}
	return s
}

type iamaws struct {
	svc *iam.IAM
}

// User Related

//CreateUser -
func (s *iamaws) CreateUser(user string) error {
	params := &iam.CreateUserInput{
		UserName: aws.String(user),
	}
	_, err := s.svc.CreateUser(params)
	if err != nil {
		fmt.Println(err.Error())
	}
	return err
}

//DeleteUser -
func (s *iamaws) DeleteUser(user string) error {
	params := &iam.DeleteUserInput{
		UserName: aws.String(user),
	}
	_, err := s.svc.DeleteUser(params)
	if err != nil {
		fmt.Println(err.Error())
	}
	return err
}

//EmailableUsers -
func (s *iamaws) EmailableUsers() ([]string, error) {
	users, err := s.ListUsers()

	var emailableUsers []string
	for i := 0; i < len(users); i++ {
		if strings.Contains(users[i], "@") {
			emailableUsers = append(emailableUsers, users[i])
		}
	}
	return emailableUsers, err
}

//ListUsers -
func (s *iamaws) ListUsers() ([]string, error) {
	params := &iam.ListUsersInput{
		MaxItems: aws.Int64(320),
	}

	resp, err := s.svc.ListUsers(params)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	if len(resp.Users) > 0 {
		var users []string
		for i := 0; i < len(resp.Users); i++ {
			users = append(users, *resp.Users[i].UserName)
		}
		return users, err
	}

	return nil, err
}

// Group Related

//CreateGroup -
func (s *iamaws) CreateGroup(group string) error {
	params := &iam.CreateGroupInput{
		GroupName: aws.String(group),
	}
	_, err := s.svc.CreateGroup(params)
	if err != nil {
		fmt.Println(err.Error())
	}
	return err
}

//DeleteGroup -
func (s *iamaws) DeleteGroup(group string) error {
	params := &iam.DeleteGroupInput{
		GroupName: aws.String(group),
	}
	_, err := s.svc.DeleteGroup(params)
	if err != nil {
		fmt.Println(err.Error())
	}
	return err
}

//ListGroupsForUser -
func (s *iamaws) ListGroupsForUser(user string) ([]string, error) {
	params := &iam.ListGroupsForUserInput{
		UserName: aws.String(user),
	}
	resp, err := s.svc.ListGroupsForUser(params)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	if len(resp.Groups) > 0 {
		var groups []string
		for i := 0; i < len(resp.Groups); i++ {
			groups = append(groups, *resp.Groups[i].GroupName)
		}
		return groups, err
	}

	return nil, err
}

// Access Key Related

//ListAccessKeys -
func (s *iamaws) ListAccessKeys(user string) ([]string, error) {
	params := &iam.ListAccessKeysInput{
		UserName: aws.String(user),
	}
	resp, err := s.svc.ListAccessKeys(params)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	if len(resp.AccessKeyMetadata) > 0 {
		var accesskeys []string
		for i := 0; i < len(resp.AccessKeyMetadata); i++ {
			accesskeys = append(accesskeys, *resp.AccessKeyMetadata[i].AccessKeyId)
		}
		return accesskeys, err
	}

	return nil, err
}

//CreateAccessKey -
func (s *iamaws) CreateAccessKey(user string) error {
	params := &iam.CreateAccessKeyInput{
		UserName: aws.String(user),
	}
	_, err := s.svc.CreateAccessKey(params)

	if err != nil {
		fmt.Println(err.Error())
	}

	return err
}

//DeleteAccessKey -
func (s *iamaws) DeleteAccessKey(user string, key string) error {
	params := &iam.DeleteAccessKeyInput{
		UserName:    aws.String(user),
		AccessKeyId: aws.String(key),
	}
	_, err := s.svc.DeleteAccessKey(params)

	if err != nil {
		fmt.Println(err.Error())
	}

	return err
}

// Combination Related

//AddUserToGroup -
func (s *iamaws) AddUserToGroup(user string, group string) error {
	params := &iam.AddUserToGroupInput{
		UserName:  aws.String(user),
		GroupName: aws.String(group),
	}
	_, err := s.svc.AddUserToGroup(params)

	if err != nil {
		fmt.Println(err.Error())
	}

	return err
}

//RemoveUserFromGroup -
func (s *iamaws) RemoveUserFromGroup(user string, group string) error {
	params := &iam.RemoveUserFromGroupInput{
		UserName:  aws.String(user),
		GroupName: aws.String(group),
	}
	_, err := s.svc.RemoveUserFromGroup(params)

	if err != nil {
		fmt.Println(err.Error())
	}

	return err
}
