package grpcclient

import (
	"context"
	"log"
	"time"

	pb "github.com/thegoodparticle/challan-data-svc/vehicledata"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	host       string
	portNumber string
}

func New(host, portNumber string) *GRPCClient {
	return &GRPCClient{
		host:       host,
		portNumber: portNumber,
	}
}

func (c *GRPCClient) initClientConnection() *grpc.ClientConn {
	conn, err := grpc.Dial(c.host+":"+c.portNumber, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("fail to dial: %v", err)
		return nil
	}

	return conn
}

func (c *GRPCClient) GetVehicleDetailsByRegistrationNumber(registrationNumber string) *pb.VehicleInfo {
	conn := c.initClientConnection()
	if conn == nil {
		return nil
	}

	defer conn.Close()

	client := pb.NewVehicleDataClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	vehicleInfo, err := client.GetVehicleDataByRegistration(ctx, &pb.RegistrationRequest{RegistrationNumber: registrationNumber})
	if err != nil {
		log.Printf("client.GetVehicleDetailsByRegistrationNumber failed: %v", err)
		return nil
	}

	return vehicleInfo
}

func (c *GRPCClient) GetOwnerDetailsByID(ownerID string) *pb.OwnerInfo {
	conn := c.initClientConnection()
	if conn == nil {
		return nil
	}

	defer conn.Close()

	client := pb.NewVehicleDataClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	driverInfo, err := client.GetOwnerDataByID(ctx, &pb.OwnerRequest{OwnerID: ownerID})
	if err != nil {
		log.Printf("client.GetOwnerDetailsByID failed: %v", err)
		return nil
	}

	return driverInfo
}
