package models

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type ChallanInfo struct {
	VehicleRegNumber string                   `json:"vehicle_registration_number"`
	TotalFine        float32                  `json:"fine_amount"`
	Violations       []map[string]interface{} `json:"violations,omitempty"`
	CreatedAt        time.Time                `json:"createdAt"`
	UpdatedAt        time.Time                `json:"updatedAt"`
}

type ChallanResponse struct {
	VehicleRegNumber string                   `json:"vehicle_registration_number"`
	TotalFine        float32                  `json:"fine_amount"`
	Violations       []map[string]interface{} `json:"violations,omitempty"`
	CreatedAt        time.Time                `json:"createdAt"`
	UpdatedAt        time.Time                `json:"updatedAt"`
	VehicleModel     string                   `json:"vehicle_model"`
	VehicleCompany   string                   `json:"vehicle_company"`
	RegistrationDate time.Time                `json:"vehicle_registration_date"`
	LicenseDate      time.Time                `json:"license_date"`
	LicenseNumber    string                   `json:"owner_license_number"`
	LicenseOwnerName string                   `json:"owner_name"`
}

func InterfaceToModel(data interface{}) (instance *ChallanInfo, err error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return instance, err
	}

	return instance, json.Unmarshal(bytes, &instance)
}

func (p *ChallanInfo) GetFilterId() map[string]interface{} {
	return map[string]interface{}{"vehicle_registration_number": p.VehicleRegNumber}
}

func (p *ChallanInfo) TableName() string {
	return "vehicle-violations"
}

func (p *ChallanInfo) Bytes() ([]byte, error) {
	return json.Marshal(p)
}

func (p *ChallanInfo) GetMap() map[string]interface{} {
	return map[string]interface{}{
		"vehicle_registration_number": p.VehicleRegNumber,
		"fine_amount":                 p.TotalFine,
		"violations":                  p.Violations,
		"createdAt":                   p.CreatedAt.Format("2006-01-02T15:04:05-0700"),
		"updatedAt":                   p.UpdatedAt.Format("2006-01-02T15:04:05-0700"),
	}
}

func ParseDynamoAtributeToStruct(response map[string]*dynamodb.AttributeValue) (p ChallanInfo, err error) {
	if response == nil || (response != nil && len(response) == 0) {
		return p, errors.New("Item not found")
	}

	log.Printf("Response: %+v", response)

	for key, value := range response {
		if key == "vehicle_registration_number" {
			p.VehicleRegNumber = *value.S
		}
		if key == "fine_amount" {
			var fineAmount float32
			err := dynamodbattribute.Unmarshal(value, &fineAmount)
			if err != nil {
				return p, err
			}
			p.TotalFine = fineAmount
		}
		if key == "violations" {
			var violations []map[string]interface{}
			err := dynamodbattribute.Unmarshal(value, &violations)
			if err != nil {
				return p, err
			}
			p.Violations = violations
		}
		if key == "createdAt" {
			p.CreatedAt, err = time.Parse("2006-01-02T15:04:05-0700", *value.S)
		}
		if key == "updatedAt" {
			p.UpdatedAt, err = time.Parse("2006-01-02T15:04:05-0700", *value.S)
		}
		if err != nil {
			return p, err
		}
	}

	return p, nil
}

func ConvertIoReaderToStruct(data io.Reader, model interface{}) (interface{}, error) {
	if data == nil {
		return nil, errors.New("body is invalid")
	}
	return model, json.NewDecoder(data).Decode(model)
}
