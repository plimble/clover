package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/plimble/clover/oauth2"
)

type DynamoDB struct {
	db *dynamodb.DynamoDB
}

func NewDynamoDB(id, secret, region string) (*DynamoDB, error) {
	config := aws.NewConfig()
	config.WithCredentials(credentials.NewStaticCredentials(id, secret, ""))
	config.WithRegion(region)
	sess, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}

	return &DynamoDB{dynamodb.New(sess)}, nil
}

func (s *DynamoDB) GetClientWithSecret(id, secret string) (*oauth2.Client, error) {
	return nil, nil
}

func (s *DynamoDB) GetRefreshToken(refreshToken string) (*oauth2.RefreshToken, error) {
	// output := s.db.GetItem(&dynamodb.GetItemInput{
	// 	Key: nil,
	// })

	return nil, nil
}

func (s *DynamoDB) GetAuthorizeCode(code string) (*oauth2.AuthorizeCode, error) {
	return nil, nil
}

func (s *DynamoDB) GetAccessToken(accessToken string) (*oauth2.AccessToken, error) {
	return nil, nil
}

func (s *DynamoDB) SaveAccessToken(accessToken *oauth2.AccessToken) error {
	data, err := dynamodbattribute.MarshalMap(accessToken)
	if err != nil {
		return err
	}

	s.db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("oauth2_accesstoken"),
		Item:      data,
	})

	return nil
}

func (s *DynamoDB) SaveRefreshToken(refreshToken *oauth2.RefreshToken) error {
	return nil
}

func (s *DynamoDB) SaveAuthorizeCode(authCode *oauth2.AuthorizeCode) error {
	return nil
}

func (s *DynamoDB) IsAvailableScope(scopes []string) (bool, error) {
	return false, nil
}

func (s *DynamoDB) RevokeRefreshToken(refreshToken string) error {
	return nil
}

func (s *DynamoDB) RevokeAccessToken(accessToken string) error {
	return nil
}
