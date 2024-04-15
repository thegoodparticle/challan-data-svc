package migration

import (
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/thegoodparticle/challan-data-svc/models"
)

type Migration struct{}

func NewMigration() *Migration {
	return &Migration{}
}

func (r *Migration) Migrate(connection *dynamodb.DynamoDB) error {
	return r.createTable(connection)
}

func (r *Migration) createTable(connection *dynamodb.DynamoDB) error {
	table := &models.ChallanInfo{}

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("vehicle_registration_number"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("vehicle_registration_number"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(table.TableName()),
	}
	response, err := connection.CreateTable(input)
	if err != nil && (strings.Contains(err.Error(), "Table already exists") || strings.Contains(err.Error(), "Cannot create preexisting table")) {
		return nil
	}
	if response != nil && strings.Contains(response.GoString(), "TableStatus: \"CREATING\"") {
		time.Sleep(3 * time.Second)
		err = r.createTable(connection)
		if err != nil {
			return err
		}
	}
	return err
}
