package store

import (
	"errors"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/google/uuid"
	"github.com/thegoodparticle/challan-data-svc/internal/models"
)

type Store struct {
	connection *dynamodb.DynamoDB
}

func NewController(connection *dynamodb.DynamoDB) *Store {
	return &Store{connection: connection}
}

func (s *Store) ListAllChallansForVehicleNumber(challanReq *models.ChallanInfo) (enitites []models.ChallanInfo, err error) {
	log.Printf("Get Challans for Vehicle Number: %s", challanReq.VehicleRegNumber)

	response, err := s.findAllEntriesByID(challanReq.PrimaryKey(),
		challanReq.VehicleRegNumber, challanReq.TableName(), challanReq.ProjectionFields())
	if err != nil {
		return nil, err
	}

	log.Printf("Response - %+v", response)

	var challansInfo []models.ChallanInfo

	_ = dynamodbattribute.UnmarshalListOfMaps(response, &challansInfo)

	return challansInfo, nil
}

func (s *Store) Create(entity *models.ChallanInfo) (string, error) {
	entity.CreatedAt = time.Now()
	entity.ChallanID = uuid.New().String()

	_, err := s.createOrUpdate(entity.GetMap(), entity.TableName(), entity.PrimaryKey(), entity.SortKey())
	return entity.VehicleRegNumber, err
}

func (s *Store) findAllEntriesByID(pk, pkValue, tableName string, projections []string) ([]map[string]*dynamodb.AttributeValue, error) {
	keyBuilder := expression.Key(pk).Equal(expression.Value(pkValue))

	expressionBuilder := expression.NewBuilder().WithKeyCondition(keyBuilder)

	var projectionBuilder expression.ProjectionBuilder
	for _, keyName := range projections {
		projectionBuilder = projectionBuilder.AddNames(expression.Name(keyName))
	}

	expressionBuilder = expressionBuilder.WithProjection(projectionBuilder)

	expr, _ := expressionBuilder.Build()

	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String(tableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      expr.Projection(),
		Limit:                     aws.Int64(100),
	}

	resp, err := s.connection.Query(queryInput)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if resp == nil {
		log.Printf("no data found for id %s", pkValue)
		return nil, errors.New("no data found")
	}

	return resp.Items, nil
}

func (s *Store) createOrUpdate(entity interface{}, tableName, primaryKey, sortKey string) (response *dynamodb.PutItemOutput, err error) {
	entityParsed, err := dynamodbattribute.MarshalMap(entity)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.PutItemInput{
		Item:                     entityParsed,
		TableName:                aws.String(tableName),
		ConditionExpression:      aws.String("attribute_not_exists(#P) AND attribute_not_exists(#S)"),
		ExpressionAttributeNames: map[string]*string{"#P": &primaryKey, "#S": &sortKey},
	}
	return s.connection.PutItem(input)
}
