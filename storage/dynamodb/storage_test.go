// +build integration

package dynamodb

import (
	"testing"

	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/plimble/clover/oauth2"
	"github.com/stretchr/testify/require"
)

var db *DynamoDB

func setup() *DynamoDB {
	if db == nil {
		db, _ = New(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), os.Getenv("AWS_REGION"))
	}

	return db
}

func TestCRUDClient(t *testing.T) {
	setup()
	expc := &oauth2.Client{
		ID:     "c1",
		Name:   "client",
		Secret: "secret",
	}

	data, err := dynamodbattribute.MarshalMap(expc)
	require.NoError(t, err)

	_, err = db.db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("oauth_client"),
		Item:      data,
	})

	c, err := db.GetClientWithSecret(expc.ID, expc.Secret)
	require.NoError(t, err)
	require.Equal(t, expc, c)

	c, err = db.GetClient(expc.ID)
	require.NoError(t, err)
	require.Equal(t, expc, c)

	_, err = db.db.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String("oauth_client"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(c.ID),
			},
		},
	})
	require.NoError(t, err)

	c, err = db.GetClientWithSecret(c.ID, c.Secret)
	require.Equal(t, oauth2.DbNotFoundError(nil), err)
	require.Nil(t, c)

}

func TestCRUDAccessToken(t *testing.T) {
	setup()
	expat := &oauth2.AccessToken{
		AccessToken: "aaaaa",
		ClientID:    "cccccccc",
		UserID:      "uuuuuuuu",
		Expired:     123123123123,
		ExpiresIn:   3699,
		Scopes:      []string{"s1", "s2"},
	}
	err := db.SaveAccessToken(expat)
	require.NoError(t, err)

	at, err := db.GetAccessToken("aaaaa")
	require.NoError(t, err)
	require.EqualValues(t, *expat, *at)

	err = db.RevokeAccessToken(at.AccessToken)
	require.NoError(t, err)

	at, err = db.GetAccessToken(at.AccessToken)
	require.Equal(t, oauth2.DbNotFoundError(nil), err)
	require.Nil(t, at)
}

func TestCRUDRefreshToken(t *testing.T) {
	setup()
	exprt := &oauth2.RefreshToken{
		RefreshToken: "aaaaa",
		ClientID:     "cccccccc",
		UserID:       "uuuuuuuu",
		Expired:      123123123123,
		Scopes:       []string{"s1", "s2"},
	}
	err := db.SaveRefreshToken(exprt)
	require.NoError(t, err)

	rt, err := db.GetRefreshToken("aaaaa")
	require.NoError(t, err)
	require.Equal(t, exprt, rt)

	err = db.RevokeRefreshToken(rt.RefreshToken)
	require.NoError(t, err)

	rt, err = db.GetRefreshToken(rt.RefreshToken)
	require.Equal(t, oauth2.DbNotFoundError(nil), err)
	require.Nil(t, rt)
}

func TestCRUDAuthCode(t *testing.T) {
	setup()
	expc := &oauth2.AuthorizeCode{
		Code:     "aaaaa",
		ClientID: "cccccccc",
		UserID:   "uuuuuuuu",
		Expired:  123123123123,
		Scopes:   []string{"s1", "s2"},
	}
	err := db.SaveAuthorizeCode(expc)
	require.NoError(t, err)

	c, err := db.GetAuthorizeCode("aaaaa")
	require.NoError(t, err)
	require.Equal(t, expc, c)

	_, err = db.db.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String("oauth_authcode"),
		Key: map[string]*dynamodb.AttributeValue{
			"c": {
				S: aws.String(c.Code),
			},
		},
	})
	require.NoError(t, err)

	c, err = db.GetAuthorizeCode(c.Code)
	require.Equal(t, oauth2.DbNotFoundError(nil), err)
	require.Nil(t, c)
}

func TestIsAvailableScope(t *testing.T) {
	setup()
	db.db.BatchWriteItem(&dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			"oauth_scope": {
				{
					PutRequest: &dynamodb.PutRequest{
						Item: map[string]*dynamodb.AttributeValue{
							"id": {
								S: aws.String("s1"),
							},
						},
					},
				},
				{
					PutRequest: &dynamodb.PutRequest{
						Item: map[string]*dynamodb.AttributeValue{
							"id": {
								S: aws.String("s2"),
							},
						},
					},
				},
			},
		},
	})

	ok, err := db.IsAvailableScope([]string{"s1", "s2"})
	require.True(t, ok)
	require.NoError(t, err)

	ok, err = db.IsAvailableScope([]string{"s2", "s3"})
	require.False(t, ok)
	require.NoError(t, err)
}
