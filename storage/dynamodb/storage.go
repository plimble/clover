package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/plimble/clover/oauth2"
	"go.uber.org/zap"
)

type DynamoDB struct {
	db     *dynamodb.DynamoDB
	logger *zap.Logger
}

func New(id, secret, region string, logger *zap.Logger) (*DynamoDB, error) {
	config := aws.NewConfig()
	config.WithCredentials(credentials.NewStaticCredentials(id, secret, ""))
	config.WithRegion(region)
	sess, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}

	return &DynamoDB{dynamodb.New(sess), logger}, nil
}

func (s *DynamoDB) GetClient(id string) (*oauth2.Client, error) {
	res, err := s.db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("oauth_client"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	})

	if err != nil {
		s.logger.Error("unable to get client", zap.Error(err), zap.String("client_id", id))
		return nil, oauth2.ServerError("unable to get client", err)
	}

	if len(res.Item) == 0 {
		return nil, oauth2.NotFound(err)
	}

	c := &oauth2.Client{}
	if err = dynamodbattribute.UnmarshalMap(res.Item, c); err != nil {
		s.logger.Error("error unmarshalMap client", zap.Error(err))
		return nil, oauth2.ServerError("unable to get client", err)
	}

	return c, nil
}

func (s *DynamoDB) GetClientWithSecret(id, secret string) (*oauth2.Client, error) {
	client, err := s.GetClient(id)
	if err != nil {
		return nil, err
	}

	if client.Secret != secret {
		return nil, oauth2.NotFound(err)
	}

	return client, nil
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
		s.logger.Error("unable to get refreshtoken", zap.Error(err), zap.String("refreshtoken", refreshToken))
		return nil, oauth2.ServerError("unable to get refreshtoken", err)
	}

	if len(res.Item) == 0 {
		return nil, oauth2.NotFound(err)
	}

	at := &oauth2.RefreshToken{}
	if err = dynamodbattribute.UnmarshalMap(res.Item, at); err != nil {
		s.logger.Error("error UnmarshalMap refreshtoken", zap.Error(err))
		return nil, oauth2.ServerError("unable to get refreshtoken", err)
	}

	return at, nil
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
		s.logger.Error("unable to get code", zap.Error(err), zap.String("code", code))
		return nil, oauth2.ServerError("unable to get authorize code", err)
	}

	if len(res.Item) == 0 {
		return nil, oauth2.NotFound(err)
	}

	at := &oauth2.AuthorizeCode{}
	if err = dynamodbattribute.UnmarshalMap(res.Item, at); err != nil {
		s.logger.Error("error UnmarshalMap code", zap.Error(err))
		return nil, oauth2.ServerError("unable to get authorize code", err)
	}

	return at, nil
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
		s.logger.Error("unable to get accesstoken", zap.Error(err), zap.String("acesstoken", accessToken))
		return nil, oauth2.ServerError("unable to get accesstoken", err)
	}

	if len(res.Item) == 0 {
		return nil, oauth2.NotFound(err)
	}

	at := &oauth2.AccessToken{}
	if err = dynamodbattribute.UnmarshalMap(res.Item, at); err != nil {
		s.logger.Error("error UnmarshalMap aceesstoken", zap.Error(err))
		return nil, oauth2.ServerError("unable to get accesstoken", err)
	}

	return at, nil
}

func (s *DynamoDB) SaveAccessToken(accessToken *oauth2.AccessToken) error {
	data, err := dynamodbattribute.MarshalMap(accessToken)
	if err != nil {
		s.logger.Error("error MarshalMap aceesstoken", zap.Error(err))
		return oauth2.ServerError("unable to save accesstoken", err)
	}

	if _, err = s.db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("oauth_accesstoken"),
		Item:      data,
	}); err != nil {
		s.logger.Error("unable to save accesstoken", zap.Error(err), zap.Any("acesstoken", accessToken))
		return oauth2.ServerError("unable to save accesstoken", err)
	}

	return nil
}

func (s *DynamoDB) SaveRefreshToken(refreshToken *oauth2.RefreshToken) error {
	data, err := dynamodbattribute.MarshalMap(refreshToken)
	if err != nil {
		s.logger.Error("error MarshalMap refreshtoken", zap.Error(err))
		return oauth2.ServerError("unable to save refreshtoken", err)
	}

	if _, err = s.db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("oauth_refreshtoken"),
		Item:      data,
	}); err != nil {
		s.logger.Error("unable to save refreshtoken", zap.Error(err), zap.Any("refreshtoken", refreshToken))
		return oauth2.ServerError("unable to save refreshtoken", err)
	}

	return err
}

func (s *DynamoDB) SaveAuthorizeCode(authCode *oauth2.AuthorizeCode) error {
	data, err := dynamodbattribute.MarshalMap(authCode)
	if err != nil {
		s.logger.Error("error MarshalMap code", zap.Error(err))
		return oauth2.ServerError("unable to save authorize code", err)
	}

	if _, err = s.db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("oauth_authcode"),
		Item:      data,
	}); err != nil {
		s.logger.Error("unable to save code", zap.Error(err), zap.Any("code", authCode))
		return oauth2.ServerError("unable to save authorize code", err)
	}

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
	if err != nil {
		s.logger.Error("unable to get scopes", zap.Error(err), zap.Any("scopes", scopes))
		return false, oauth2.ServerError("unable to get scopes", err)
	}

	if _, ok := res.Responses["oauth_scope"]; ok && len(res.Responses["oauth_scope"]) == len(scopes) {
		return true, nil
	}

	return false, nil
}

func (s *DynamoDB) RevokeRefreshToken(refreshToken string) error {
	if _, err := s.db.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String("oauth_refreshtoken"),
		Key: map[string]*dynamodb.AttributeValue{
			"rt": {
				S: aws.String(refreshToken),
			},
		},
	}); err != nil {
		s.logger.Error("unable to revoke refreshtoken", zap.Error(err), zap.Any("refreshtoken", refreshToken))
		return oauth2.ServerError("unable to revoke refreshtoken", err)
	}

	return nil
}

func (s *DynamoDB) RevokeAccessToken(accessToken string) error {
	if _, err := s.db.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String("oauth_accesstoken"),
		Key: map[string]*dynamodb.AttributeValue{
			"at": {
				S: aws.String(accessToken),
			},
		},
	}); err != nil {
		s.logger.Error("unable to revoke accesstoken", zap.Error(err), zap.Any("accesstoken", accessToken))
		return oauth2.ServerError("unable to revoke accesstoken", err)
	}

	return nil
}
