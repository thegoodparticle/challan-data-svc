package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/joho/godotenv"
	grpcclient "github.com/thegoodparticle/challan-data-svc/grpc-client"
	"github.com/thegoodparticle/challan-data-svc/migration"
	"github.com/thegoodparticle/challan-data-svc/routes"
	"github.com/thegoodparticle/challan-data-svc/store"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error while loadinf env file - %s", err)
	}

	httpPortNumber := os.Getenv("HTTP_PORT")
	httpTimeout := os.Getenv("HTTP_TIMEOUT")

	connection := GetConnection()
	if connection == nil {
		log.Fatal("could not connect to dynamodb")
	}
	repository := store.NewAdapter(connection)

	log.Printf("Waiting service starting.... %+v ", connection)

	errors := Migrate(connection)
	if len(errors) > 0 {
		for _, err := range errors {
			if err != nil {
				log.Panic("Error on migrate: ", err)
			}
		}
	}

	if err := checkTables(connection); err != nil {
		log.Panic("", err)
	}

	repository.LoadDataIntoTables()

	grpcServerHost := os.Getenv("GRPC_SERVER_HOST")
	grpcServerPort := os.Getenv("GRPC_SERVER_PORT")
	grpcClientObj := grpcclient.New(grpcServerHost, grpcServerPort)

	port := fmt.Sprintf(":%v", httpPortNumber)
	timeout, _ := strconv.Atoi(httpTimeout)
	router := routes.NewRouter(timeout).SetRouters(repository, grpcClientObj)
	log.Print("HTTP Service running on port ", httpPortNumber)

	server := http.ListenAndServe(port, router)
	log.Fatal(server)
}

func Migrate(connection *dynamodb.DynamoDB) []error {
	var errors []error

	callMigrateAndAppendError(&errors, connection, migration.NewMigration())

	return errors
}

func callMigrateAndAppendError(errors *[]error, connection *dynamodb.DynamoDB, migration *migration.Migration) {
	err := migration.Migrate(connection)
	if err != nil {
		*errors = append(*errors, err)
	}
}

func checkTables(connection *dynamodb.DynamoDB) error {
	response, err := connection.ListTables(&dynamodb.ListTablesInput{})
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

func GetConnection() *dynamodb.DynamoDB {
	endpoint := os.Getenv("DYNAMO_ENDPOINT")

	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String("us-west-2"),
		Endpoint: aws.String(endpoint),
	})
	if err != nil {
		log.Println(err)
		return nil
	}

	// Create DynamoDB client
	return dynamodb.New(sess)
}
