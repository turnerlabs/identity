package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/support"

	"github.com/turnerlabs/identity/identitymgmt"
)

//Initilaization

// NewServiceLevels -
func NewServiceLevels(region string, accessKey *string, secretKey *string, token *string) identitymgmt.Support {
	store := newSupportAWS(region, accessKey, secretKey, token)
	return &store
}

func newSupportAWS(region string, accessKey *string, secretKey *string, token *string) supportaws {
	var awsConfig *aws.Config
	if accessKey == nil || secretKey == nil || token == nil {
		awsConfig = aws.NewConfig().WithRegion(region)
	} else {
		awsConfig = aws.NewConfig().WithRegion(region)
		awsConfig.WithCredentials(credentials.NewStaticCredentials(*accessKey, *secretKey, *token))
	}

	awsConfig.WithCredentialsChainVerboseErrors(true)
	session := session.New(awsConfig)

	s := supportaws{
		svc: support.New(session),
	}
	return s
}

type supportaws struct {
	svc *support.Support
}

func (s *supportaws) ListServiceLevels(checkid string) ([]identitymgmt.ServiceLimit, error) {
	params := &support.DescribeTrustedAdvisorCheckResultInput{
		CheckId:  aws.String(checkid), // Required
		Language: aws.String("en"),
	}
	resp, err := s.svc.DescribeTrustedAdvisorCheckResult(params)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	if len(resp.Result.FlaggedResources) > 0 {
		var serviceLimits []identitymgmt.ServiceLimit
		for i := 0; i < len(resp.Result.FlaggedResources); i++ {
			serviceLimit := new(identitymgmt.ServiceLimit)
			serviceLimit.Region = *resp.Result.FlaggedResources[i].Metadata[0]
			serviceLimit.ServiceName = *resp.Result.FlaggedResources[i].Metadata[1]
			serviceLimit.ServiceItem = *resp.Result.FlaggedResources[i].Metadata[2]
			serviceLimit.Max = *resp.Result.FlaggedResources[i].Metadata[3]
			if resp.Result.FlaggedResources[i].Metadata[4] != nil {
				serviceLimit.Current = *resp.Result.FlaggedResources[i].Metadata[4]
			} else {
				serviceLimit.Current = "0"
			}
			serviceLimit.Color = *resp.Result.FlaggedResources[i].Metadata[5]

			serviceLimits = append(serviceLimits, *serviceLimit)
		}
		return serviceLimits, err
	}

	return nil, err

}
