package dynamodb

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/patrickmn/go-cache"
	"github.com/plimble/clover/oauth2"
)

type DynamoDB struct {
	db    *dynamodb.DynamoDB
	cache *cache.Cache
}

func New(id, secret, region string) (*DynamoDB, error) {
	config := aws.NewConfig()
	config.WithCredentials(credentials.NewStaticCredentials(id, secret, ""))
	config.WithRegion(region)
	sess, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}

	cache := cache.New(cache.NoExpiration, 10*time.Minute)

	return &DynamoDB{dynamodb.New(sess), cache}, nil
}

func (s *DynamoDB) GetClientWithSecret(id, secret string) (*oauth2.Client, error) {
	res, err := s.db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("oauth_client"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
			"s": {
				S: aws.String(secret),
			},
		},
	})

	if err != nil {
		return nil, err
	}

	if len(res.Item) == 0 {
		return nil, oauth2.DbNotFoundError(err)
	}

	c := &oauth2.Client{}
	err = dynamodbattribute.UnmarshalMap(res.Item, c)

	return c, err
}

func (s *DynamoDB) GetRefreshToken(refreshToken string) (*oauth2.RefreshToken, error) {
	res, err := s.db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("oauth_refreshtoken"),
		Key: map[string]*dynamodb.AttributeValue{
			"rt": {
				S: aws.String(refreshToken),
			},
		},
	})

	if err != nil {
		return nil, err
	}

	if len(res.Item) == 0 {
		return nil, oauth2.DbNotFoundError(err)
	}

	at := &oauth2.RefreshToken{}
	err = dynamodbattribute.UnmarshalMap(res.Item, at)

	return at, err
}

func (s *DynamoDB) GetAuthorizeCode(code string) (*oauth2.AuthorizeCode, error) {
	res, err := s.db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("oauth_authcode"),
		Key: map[string]*dynamodb.AttributeValue{
			"c": {
				S: aws.String(code),
			},
		},
	})

	if err != nil {
		return nil, err
	}

	if len(res.Item) == 0 {
		return nil, oauth2.DbNotFoundError(err)
	}

	at := &oauth2.AuthorizeCode{}
	err = dynamodbattribute.UnmarshalMap(res.Item, at)

	return at, err
}

func (s *DynamoDB) GetAccessToken(accessToken string) (*oauth2.AccessToken, error) {
	res, err := s.db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("oauth_accesstoken"),
		Key: map[string]*dynamodb.AttributeValue{
			"at": {
				S: aws.String(accessToken),
			},
		},
	})

	if err != nil {
		return nil, err
	}

	if len(res.Item) == 0 {
		return nil, oauth2.DbNotFoundError(err)
	}

	at := &oauth2.AccessToken{}
	err = dynamodbattribute.UnmarshalMap(res.Item, at)

	return at, err
}

func (s *DynamoDB) SaveAccessToken(accessToken *oauth2.AccessToken) error {
	data, err := dynamodbattribute.MarshalMap(accessToken)
	if err != nil {
		return err
	}

	_, err = s.db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("oauth_accesstoken"),
		Item:      data,
	})

	return err
}

func (s *DynamoDB) SaveRefreshToken(refreshToken *oauth2.RefreshToken) error {
	data, err := dynamodbattribute.MarshalMap(refreshToken)
	if err != nil {
		return err
	}

	_, err = s.db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("oauth_refreshtoken"),
		Item:      data,
	})

	return err
}

func (s *DynamoDB) SaveAuthorizeCode(authCode *oauth2.AuthorizeCode) error {
	data, err := dynamodbattribute.MarshalMap(authCode)
	if err != nil {
		return err
	}

	_, err = s.db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("oauth_authcode"),
		Item:      data,
	})

	return err
}

func (s *DynamoDB) IsAvailableScope(scopes []string) (bool, error) {
	items := map[string]*dynamodb.KeysAndAttributes{
		"oauth_scope": {
			Keys: make([]map[string]*dynamodb.AttributeValue, len(scopes)),
		},
	}
	for i, scope := range scopes {
		items["oauth_scope"].Keys[i] = map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(scope),
			},
		}
	}

	res, err := s.db.BatchGetItem(&dynamodb.BatchGetItemInput{
		RequestItems: items,
	})

	if _, ok := res.Responses["oauth_scope"]; ok && len(res.Responses["oauth_scope"]) == len(scopes) {
		return true, nil
	}

	return false, err
}

func (s *DynamoDB) RevokeRefreshToken(refreshToken string) error {
	_, err := s.db.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String("oauth_refreshtoken"),
		Key: map[string]*dynamodb.AttributeValue{
			"rt": {
				S: aws.String(refreshToken),
			},
		},
	})

	return err
}

func (s *DynamoDB) RevokeAccessToken(accessToken string) error {
	_, err := s.db.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String("oauth_accesstoken"),
		Key: map[string]*dynamodb.AttributeValue{
			"at": {
				S: aws.String(accessToken),
			},
		},
	})

	return err
}
