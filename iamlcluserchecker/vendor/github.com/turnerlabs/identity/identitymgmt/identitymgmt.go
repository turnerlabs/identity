package identitymgmt

// Cloudhealth

//CloudhealthAccount -
type CloudhealthAccount struct {
	BasicAccount
	OtherCloudhealthAccount
	CloudhealthStatus     `json:"status"`
	CloudhealthBilling    `json:"billing"`
	CloudhealthCloudtrail `json:"cloudtrail"`
	CloudhealthAWSConfig  `json:"aws_config"`
	CloudhealthCloudwatch `json:"cloudwatch"`
	Groups                []CloudhealthGroup `json:"groups"`
}

//OtherCloudhealthAccount -
type OtherCloudhealthAccount struct {
	AmazonName  string `json:"amazon_name"`
	OwnerID     string `json:"owner_id"`
	Region      string `json:"region"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	AccountType string `json:"account_type"`
	VPCOnly     bool   `json:"vpc_only"`
	ClusterName string `json:"cluster_name"`
}

//CloudhealthStatus -
type CloudhealthStatus struct {
	Level      string `json:"level"`
	LastUpdate string `json:"last_update"`
}

//CloudhealthBilling -
type CloudhealthBilling struct {
	Bucket         string `json:"bucket"`
	IsConsolidated bool   `json:"is_consolidated"`
}

//CloudhealthCloudtrail -
type CloudhealthCloudtrail struct {
	Enabled bool   `json:"enabled"`
	Bucket  string `json:"bucket"`
}

//CloudhealthAWSConfig -
type CloudhealthAWSConfig struct {
	Enabled bool `json:"enabled"`
}

//CloudhealthCloudwatch -
type CloudhealthCloudwatch struct {
	Enabled bool `json:"enabled"`
}

//CloudhealthGroup -
type CloudhealthGroup struct {
	Name  string `json:"name"`
	Group string `json:"group"`
}

// AWS

//ServiceLimit -
type ServiceLimit struct {
	Region      string
	ServiceName string
	ServiceItem string
	Max         string
	Current     string
	Color       string
}

//Support -
type Support interface {
	ListServiceLevels(checkid string) ([]ServiceLimit, error)
}

//Identity -
type Identity interface {
	//CreateUser
	CreateUser(user string) error
	//DeleteUser
	DeleteUser(user string) error
	//ListUsers
	ListUsers() ([]string, error)
	//EmailableUsers
	EmailableUsers() ([]string, error)
	//CreateGroup
	CreateGroup(group string) error
	//DeleteGroup
	DeleteGroup(group string) error
	//ListGroupsForUser
	ListGroupsForUser(user string) ([]string, error)
	//CreateAccessKey
	CreateAccessKey(user string) error
	//DeleteAccessKey
	DeleteAccessKey(user, key string) error
	//ListAccessKeys
	ListAccessKeys(user string) ([]string, error)
	//AddUserToGroup
	AddUserToGroup(user, group string) error
	//RemoveUserFromGroup
	RemoveUserFromGroup(user, group string) error
}

// Github

//Membership -
type Membership interface {
	//Members
	Members(org string) ([]GithubMember, error)
}

//GithubMember -
type GithubMember struct {
	Login string `json:"login"`
	Type  string `json:"type"`
}

//BasicAccount -
type BasicAccount struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

//SimpleIdentity -
type SimpleIdentity interface {
	//Validate
	Validate(user string) (bool, error)
}

//SimpleRole -
type SimpleRole interface {
	//AssumeRole
	AssumeRole(roleArn, roleSessionName string) (*string, *string, *string, error)
}

//SimpleAccount -
type SimpleAccount interface {
	//GetAccounts -
	GetAccounts() ([]BasicAccount, error)
	// //GetAccount -
	GetAccount(id string) (CloudhealthAccount, error)
}
