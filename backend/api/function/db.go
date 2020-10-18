package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
	"os"
	"strings"
)


var newSession, _ = session.NewSession()
var db = dynamodb.New(newSession)

const(
	//LocationsTable Name of dynamoDB table storing locations
	LocationsTable = "Locations"
)

func getDestination(locationName string) (*Location, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(LocationsTable),
		Key: map[string]*dynamodb.AttributeValue{
			"LocationName": {
				S: aws.String(locationName),
			},
		},
	}

	result, err := db.GetItem(input)
	if err != nil {
		return nil, err
	}
	if result.Item == nil {
		return nil, nil
	}

	location := new(Location)
	err = dynamodbattribute.UnmarshalMap(result.Item, location)
	if err != nil {
		return nil, err
	}

	return location, nil
}
// generatePrimaryKey Generate a primary key using UUID with removed dashes.
func generatePrimaryKey() string{
	uuidWithHyphen := uuid.New()
	primaryKey := strings.Replace(uuidWithHyphen.String(), "-", "", -1)
	return primaryKey
}

// Add a Location record to DynamoDB.
func insertDestination(location *Location) error {
	var uuid = generatePrimaryKey()
	//Insert the new tuple with a generated location ID.
	location.LocationID = uuid
	attributeValues, err := dynamodbattribute.MarshalMap(location)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	input := &dynamodb.PutItemInput{
		TableName: aws.String(LocationsTable),
		Item:      attributeValues,
	}
	_, err = db.PutItem(input)
	return err
}

