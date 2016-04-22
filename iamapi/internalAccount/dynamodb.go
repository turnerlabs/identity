package internalAccount

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var (
	accountsTableName               = aws.String("iamaccountextra")
	conditionalCheckFailedException = "ConditionalCheckFailedException"
)

//dynamodbStore
type dynamodbStore struct {
	dynamodbService        *dynamodb.DynamoDB
	reportConsumedCapacity bool
	consumedCapacity       *string
}

// accountObject -
type accountObject struct {
	AccountID    string `json:"AccountID"`
	AccountItems []AccountItem
}

func toAccountObject(n *Account) *accountObject {
	return &accountObject{
		AccountID:    n.AccountID,
		AccountItems: n.AccountItems,
	}
}

func (n *accountObject) toAccount() *Account {
	return &Account{
		AccountID:    n.AccountID,
		AccountItems: n.AccountItems,
	}
}

type accountKey struct {
	AccountID string
}

func (s *dynamodbStore) toAccountKey(accountID string) *accountKey {
	return &accountKey{
		AccountID: accountID,
	}
}

// NewAccountStore -
func NewAccountStore(region string, reportConsumedCapacity bool) AccountStore {
	store := newAccountStore(region, reportConsumedCapacity)
	return &store
}

func newAccountStore(region string, reportConsumedCapacity bool) dynamodbStore {
	return newStoreFromDynamoDB(dynamodb.New(session.New(aws.NewConfig().WithRegion(region))), reportConsumedCapacity)
}

func newStoreFromDynamoDB(ddb *dynamodb.DynamoDB, reportConsumedCapacity bool) dynamodbStore {
	s := dynamodbStore{
		dynamodbService:        ddb,
		reportConsumedCapacity: reportConsumedCapacity,
	}

	if reportConsumedCapacity {
		s.consumedCapacity = aws.String(dynamodb.ReturnConsumedCapacityIndexes)
	} else {
		s.consumedCapacity = aws.String(dynamodb.ReturnConsumedCapacityNone)
	}
	return s
}

func (s *dynamodbStore) executeDynamodbQuery(tableName, indexName *string, consistentRead bool, keyConditionExpression string,
	expressionAttributeValues map[string]*dynamodb.AttributeValue, lastEvaluatedKey map[string]*dynamodb.AttributeValue, projectionExpression *string,
	expressionAttributeNames map[string]*string) (*dynamodb.QueryOutput, error) {

	input := &dynamodb.QueryInput{
		TableName:                 tableName,
		IndexName:                 indexName,
		ConsistentRead:            aws.Bool(consistentRead),
		KeyConditionExpression:    aws.String(keyConditionExpression),
		ExpressionAttributeValues: expressionAttributeValues,
		ReturnConsumedCapacity:    s.consumedCapacity,
		ExclusiveStartKey:         lastEvaluatedKey,
		ProjectionExpression:      projectionExpression,
		ExpressionAttributeNames:  expressionAttributeNames,
	}

	return s.dynamodbService.Query(input)
}

func (s *dynamodbStore) executeDynamodbGetItem(tableName *string, key map[string]*dynamodb.AttributeValue) (*dynamodb.GetItemOutput, error) {
	input := &dynamodb.GetItemInput{
		TableName: tableName,
		Key:       key,
		ReturnConsumedCapacity: s.consumedCapacity,
	}

	return s.dynamodbService.GetItem(input)
}

func (s *dynamodbStore) executeDynamodbPutItem(tableName *string, item map[string]*dynamodb.AttributeValue) (*dynamodb.PutItemOutput, error) {
	input := &dynamodb.PutItemInput{
		TableName: tableName,
		Item:      item,
		ReturnConsumedCapacity: s.consumedCapacity,
	}

	return s.dynamodbService.PutItem(input)
}

//executeDynamodbDeleteItem - helper function for dynamodb DeleteItem function
func (s *dynamodbStore) executeDynamodbDeleteItem(tableName *string, key map[string]*dynamodb.AttributeValue) (*dynamodb.DeleteItemOutput, error) {
	input := &dynamodb.DeleteItemInput{
		TableName: tableName,
		Key:       key,
		ReturnConsumedCapacity: s.consumedCapacity,
	}

	return s.dynamodbService.DeleteItem(input)
}

// Add - add Account
func (s *dynamodbStore) Add(entry Account) error {
	if len(entry.AccountID) == 0 {
		return ErrMissingAccountID
	}

	// convert public to internal object then convert to dynamo map
	accountObj := toAccountObject(&entry)
	item, err := dynamodbattribute.ConvertToMap(*accountObj)
	if err != nil {
		return err
	}

	// calls PutItem
	_, err = s.executeDynamodbPutItem(accountsTableName, item)
	if err != nil {
		return err
	}

	return nil
}

// Update - update Account
func (s *dynamodbStore) Update(entry Account) error {
	if len(entry.AccountID) == 0 {
		return ErrMissingAccountID
	}

	// convert public to internal object then convert to dynamo map
	accountObj := toAccountObject(&entry)
	item, err := dynamodbattribute.ConvertToMap(*accountObj)
	if err != nil {
		return err
	}

	// calls PutItem
	_, err = s.executeDynamodbPutItem(accountsTableName, item)
	if err != nil {
		return err
	}

	return nil
}

// Delete - delete nicknameEntry by nickname
func (s *dynamodbStore) Delete(accountID string) error {
	if len(accountID) == 0 {
		return ErrMissingAccountID
	}

	current := s.toAccountKey(accountID)
	key, err := dynamodbattribute.ConvertToMap(*current)
	if err != nil {
		return err
	}

	_, err = s.executeDynamodbDeleteItem(accountsTableName, key)
	if err != nil {
		return err
	}

	return nil
}

// GetByAccountID -
func (s *dynamodbStore) GetByAccountID(accountID string) (*Account, error) {
	if len(accountID) == 0 {
		return nil, ErrMissingAccountID
	}

	current := s.toAccountKey(accountID)
	key, err := dynamodbattribute.ConvertToMap(*current)
	if err != nil {
		return nil, err
	}

	output, err := s.executeDynamodbGetItem(accountsTableName, key)
	if err != nil {
		return nil, err
	}

	if len(output.Item) == 0 {
		return nil, ErrNoRowFound
	}

	item := accountObject{}
	err = dynamodbattribute.ConvertFromMap(output.Item, &item)
	if err != nil {
		return nil, err
	}

	var account = item.toAccount()

	return account, nil
}
