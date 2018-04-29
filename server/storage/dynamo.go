package storage

import (
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/dblooman/baffle/server/backends"
	"github.com/twinj/uuid"
)

func Put(d *dynamodb.DynamoDB, data backends.CreateSecret) error {

	params := &dynamodb.PutItemInput{
		TableName: aws.String("baffles"),
		Item: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(uuid.NewV4().String()),
			},
			"name": {
				S: aws.String(data.Name),
			},
			"fragment": {
				S: aws.String(data.Fragement),
			},
			"regex": {
				S: aws.String(data.Regex),
			},
			"path": {
				S: aws.String(data.Path),
			},
			"version": {
				S: aws.String(strconv.FormatInt(data.Version, 10)),
			},
		},
	}

	_, err := d.PutItem(params)
	if err != nil {
		return err
	}

	return nil
}
