package store

import (
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/thegoodparticle/challan-data-svc/internal/models"
)

func (s *Store) Migrate() error {
	return s.createTable(s.connection)
}

func (s *Store) CheckTables() error {
	response, err := s.connection.ListTables(&dynamodb.ListTablesInput{})
	if response != nil {
		if len(response.TableNames) == 0 {
			log.Print("Tables not found: ", nil)
		}

		for _, tableName := range response.TableNames {
			log.Print("Table found: ", *tableName)
		}
	}

	return err
}

func (s *Store) createTable(connection *dynamodb.DynamoDB) error {
	table := &models.ChallanInfo{}

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Vehicle Number"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("Challan ID"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Vehicle Number"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("Challan ID"),
				KeyType:       aws.String("RANGE"),
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
		err = s.createTable(connection)
		if err != nil {
			return err
		}
	}
	return err
}

var challans []models.ChallanInfo = []models.ChallanInfo{
	{
		VehicleRegNumber: "KA20AB1234",
		ChallanID:        "70c6cf47-43e8-4b98-acb2-f20a3d41d8a5",
		UnitName:         "Cyberabad",
		Date:             "02-Apr-2024",
		Time:             "15:33",
		PlaceOfViolation: "RAIDURGAM TR PS LIMITS",
		PSLimits:         "Raidurgam Tr PS",
		FineAmount:       100,
		Violation:        "Wrong Parking in the carriage way",
		CreatedAt:        time.Now(),
	},
	{
		VehicleRegNumber: "KA20AB1234",
		ChallanID:        "dbb6d85b-4956-420d-b672-245c4e86ad95",
		UnitName:         "Cyberabad",
		Date:             "16-Apr-2024",
		Time:             "21:33",
		PlaceOfViolation: "RAIDURGAM TR PS LIMITS",
		PSLimits:         "Raidurgam Tr PS",
		FineAmount:       2000,
		Violation:        "Driving vehicle in one way",
		CreatedAt:        time.Now(),
	},
	{
		VehicleRegNumber: "KA20CD5678",
		ChallanID:        "9c29f53f-d325-41a4-b9aa-f3f31cb72dd7",
		UnitName:         "Hyderabad",
		Date:             "11-Apr-2024",
		Time:             "09:33",
		PlaceOfViolation: "HYDER TR PS LIMITS",
		PSLimits:         "Hyder Tr PS",
		FineAmount:       1000,
		Violation:        "Riding without helmet",
		CreatedAt:        time.Now(),
	},
}

func (s *Store) LoadDataIntoTables() {
	for _, violation := range challans {
		_, err := s.createOrUpdate(violation.GetMap(), violation.TableName(), violation.PrimaryKey(), violation.SortKey())
		if err != nil {
			log.Printf("error while adding entry into the table. Error - %+v", err)
		}
	}
}
