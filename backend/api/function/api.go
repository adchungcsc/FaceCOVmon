package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)


var newSession, _ = session.NewSession()
var db = dynamodb.New(newSession)

const(
	//ObservationsTable Name of dynamoDB table storing observations
	ObservationsTable = "Observations"
)

func getTenantData(cameraID *string, date *string) (*ReturnedData, error) {
	//TEMPORARILY HARDCODED TO THE SINGLE CAMERA AND TODAY UNTIL I CAN FIGURE OUT WHY MY QUERIES DON'T GET THINGS THAT EXIST
	result, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(ObservationsTable),
		Key: map[string]*dynamodb.AttributeValue{
			"cameraID": {
				S: aws.String("test"),
			},
			"date": {
				S: aws.String("today"),
			},
		},
	})
	if result.Item == nil {
		return nil, err
	}
	fmt.Println("succed get")

	data := &ReturnedData{}

	fmt.Println("Trying unmarshal")
	err = dynamodbattribute.UnmarshalMap(result.Item, &data)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}
	fmt.Println("unmarshal work ")

	fmt.Println(data.CameraID + " " + data.Date + " " + data.ImproperFaceCovering + " " + data.NoFaceCovering + " " + data.Total )

	//err = dynamodbattribute.UnmarshalMap(result.Item, data)
	//if err != nil {
	//	return nil, err
	//}

	return data, nil
}

// Add a record to DynamoDB.
func uploadData(cameraID string, image string) error {
	//Insert the new tuple with
	input := &dynamodb.PutItemInput{
		TableName: aws.String(ObservationsTable),
		Item: map[string]*dynamodb.AttributeValue{
			"cameraID": {
				S: aws.String(cameraID),
			},
			"date": {
				S: aws.String("today"),
			},
		},
	}
	_, err := db.PutItem(input)
	return err
}

