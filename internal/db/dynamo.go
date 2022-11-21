package db

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/eduardohoraciosanto/users/internal/errors"
	"github.com/eduardohoraciosanto/users/internal/logger"
)

type dynamo struct {
	client    *dynamodb.Client
	tableName string
	log       logger.Logger
}

func NewDynamoDB(tableName string, dynamoClient *dynamodb.Client, log logger.Logger) DB {
	return &dynamo{
		client:    dynamoClient,
		tableName: tableName,
		log:       log,
	}
}

func (d *dynamo) Set(ctx context.Context, key string, data interface{}) error {
	l := d.log.WithField("key", key).WithField("data", data)

	item, err := attributevalue.MarshalMap(data)
	if err != nil {
		l.WithError(err).Error(ctx, "unable to marshall data into dynamo attributes")
		return err
	}
	_, err = d.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(d.tableName), Item: item,
	})
	if err != nil {
		l.WithError(err).Error(ctx, "unable to set item on DB")
		return err
	}
	return nil
}

func (d *dynamo) Get(ctx context.Context, key string, here interface{}) error {
	l := d.log.WithField("key", key)

	mKey, err := attributevalue.Marshal(key)
	if err != nil {
		l.WithError(err).Error(ctx, "marshall key attributevalue")
		return err
	}
	response, err := d.client.GetItem(ctx, &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{"email": mKey}, TableName: aws.String(d.tableName),
	})
	if err != nil {
		l.WithError(err).Error(ctx, "unable to get item from DB")
		return err
	} else {
		if response.Item == nil {
			l.WithError(err).Error(ctx, "element not found on DB")
			return errors.New(errors.DBErrorNotFoundCode, "element not found on DB")
		}
		err = attributevalue.UnmarshalMap(response.Item, here)
		if err != nil {
			l.WithError(err).Error(ctx, "unable to unmarshal item from DB")
			return err
		}
	}
	return nil
}

func (d *dynamo) Delete(ctx context.Context, key string) error {
	mKey, err := attributevalue.Marshal(key)
	if err != nil {
		return err
	}

	_, err = d.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(d.tableName), Key: map[string]types.AttributeValue{"email": mKey},
	})
	if err != nil {
		d.log.WithField("key", key).WithError(err).Error(ctx, "unable to delete item from DB")
		return err
	}
	return err
}
func (d *dynamo) Alive(ctx context.Context) bool {
	exists := true
	_, err := d.client.DescribeTable(
		ctx, &dynamodb.DescribeTableInput{TableName: aws.String(d.tableName)},
	)
	if err != nil {
		d.log.WithField("table_name", d.tableName).WithError(err).Error(ctx, "unable to determine table existence")
		exists = false
	}

	return exists
}
