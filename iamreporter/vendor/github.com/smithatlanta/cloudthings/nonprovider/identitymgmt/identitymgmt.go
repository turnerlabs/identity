package identitymgmt

//Identity -
type Identity interface {
	//CreateUser
	CreateUser(user string) error
	//DeleteUser
	DeleteUser(user string) error
	//ListUsers
	ListUsers() ([]string, error)
	//CreateGroup
	CreateGroup(group string) error
	//DeleteGroup
	DeleteGroup(group string) error
	//ListGroupsForUser
	ListGroupsForUser(user string) ([]string, error)
	//CreateAccessKey
	CreateAccessKey(user string) error
	//DeleteAccessKey
	DeleteAccessKey(user string, key string) error
	//ListAccessKeys
	ListAccessKeys(user string) ([]string, error)
	//AddUserToGroup
	AddUserToGroup(user string, group string) error
	//RemoveUserFromGroup
	RemoveUserFromGroup(user string, group string) error
}

//SimpleIdentity -
type SimpleIdentity interface {
	//Validate
	Validate(user string) (bool, error)
}

//SimpleRole -
type SimpleRole interface {
	//AssumeRole
	AssumeRole(roleArn string, roleSessionName string) (*string, *string, *string, error)
}
