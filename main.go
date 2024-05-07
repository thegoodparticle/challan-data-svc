package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	eventconsumer "github.com/thegoodparticle/challan-data-svc/event-consumer"
	grpcclient "github.com/thegoodparticle/challan-data-svc/grpc-client"
	"github.com/thegoodparticle/challan-data-svc/internal/store"
	handler "github.com/thegoodparticle/challan-data-svc/rest-handler"
)

func main() {

	// load env variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error while loadinf env file - %s", err)
	}

	endpoint := os.Getenv("DYNAMO_ENDPOINT")
	region := os.Getenv("DYNAMO_REGION")

	// get dynamoDB connection
	connection := getConnection(endpoint, region)
	if connection == nil {
		log.Fatal("could not connect to dynamodb")
	}

	log.Printf("Waiting service starting.... %+v ", connection)

	// setup data store - interface for dynamoDB operations
	dataStore := store.NewController(connection)

	err = dataStore.Migrate()
	if err != nil {
		log.Panic("Error on migrate: ", err)
	}

	if err := dataStore.CheckTables(); err != nil {
		log.Panic("Error on Tables Check", err)
	}

	// load few pre-defined data into DB tables
	dataStore.LoadDataIntoTables()

	// get grpc server details from env, setup the client
	grpcServerHost := os.Getenv("GRPC_SERVER_HOST")
	grpcServerPort := os.Getenv("GRPC_SERVER_PORT")
	grpcClientObj := grpcclient.New(grpcServerHost, grpcServerPort)

	// extract http port number for rest handlers
	httpPortNumber := os.Getenv("HTTP_PORT")

	// setup rest api handlers and routes for defining the endpoints
	port := fmt.Sprintf(":%v", httpPortNumber)

	h := handler.NewHandler(dataStore, grpcClientObj)

	router := chi.NewRouter()

	router.Route("/challan-info", func(route chi.Router) {
		route.Get("/{RegID}", h.GetChallanResponseForRegistrationID)
	})

	log.Print("HTTP Service running on port ", httpPortNumber)

	go func() {
		server := http.ListenAndServe(port, router)
		log.Print(server)
	}()

	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	kafkaTopic := os.Getenv("KAFKA_TOPIC")

	brokersList := strings.Split(kafkaBrokers, ",")

	eventconsumer.SetupConsumer(brokersList, kafkaTopic, h)
}

func getConnection(endpoint, region string) *dynamodb.DynamoDB {
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String(region),
		Endpoint: aws.String(endpoint),
	})
	if err != nil {
		log.Println(err)
		return nil
	}

	// Create DynamoDB client
	return dynamodb.New(sess)
}
