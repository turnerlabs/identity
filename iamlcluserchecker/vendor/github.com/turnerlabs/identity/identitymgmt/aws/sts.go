package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/turnerlabs/identity/identitymgmt"
)

//Initilaization

// NewRole -
func NewRole(region string, accessKey *string, secretKey *string, token *string) identitymgmt.SimpleRole {
	store := newSTSAWS(region, accessKey, secretKey, token)
	return &store
}

func newSTSAWS(region string, accessKey *string, secretKey *string, token *string) stsaws {
	var awsConfig *aws.Config
	if accessKey == nil || secretKey == nil || token == nil {
		awsConfig = aws.NewConfig().WithRegion(region)
	} else {
		awsConfig = aws.NewConfig().WithRegion(region)
		awsConfig.WithCredentials(credentials.NewStaticCredentials(*accessKey, *secretKey, *token))
	}

	awsConfig.WithCredentialsChainVerboseErrors(true)
	session := session.New(awsConfig)

	s := stsaws{
		svc: sts.New(session),
	}
	return s
}

type stsaws struct {
	svc *sts.STS
}

func (s *stsaws) AssumeRole(roleArn string, roleSessionName string) (*string, *string, *string, error) {
	params := &sts.AssumeRoleInput{
		RoleArn:         aws.String(roleArn),
		RoleSessionName: aws.String(roleSessionName),
	}
	resp, err := s.svc.AssumeRole(params)
	if err != nil {
		//		fmt.Println(err.Error())
		return nil, nil, nil, err
	}

	return resp.Credentials.AccessKeyId, resp.Credentials.SecretAccessKey, resp.Credentials.SessionToken, err
}
