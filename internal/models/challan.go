package models

import (
	"encoding/json"
	"errors"
	"io"
	"time"
)

type ChallanInfo struct {
	VehicleRegNumber string    `json:"Vehicle Number,omitempty"`
	ChallanID        string    `json:"Challan ID"`
	UnitName         string    `json:"Unit Name"`
	Date             string    `json:"Date"`
	Time             string    `json:"Time"`
	PlaceOfViolation string    `json:"Place of Violation"`
	PSLimits         string    `json:"PS Limits"`
	Violation        string    `json:"Violation"`
	FineAmount       int       `json:"Fine Amount"`
	CreatedAt        time.Time `json:"Created At,omitempty"`
}

type ChallanResponse struct {
	VehicleRegNumber       string        `json:"Vehicle Number"`
	VehicleModel           string        `json:"Vehicle Model"`
	VehicleCompany         string        `json:"Vehicle Company"`
	RegistrationDate       time.Time     `json:"Vehicle Registration Date"`
	RegistrationExpiryDate time.Time     `json:"Vehicle Registration Expiry Date"`
	OwnerFirstName         string        `json:"Owner First Name"`
	OwnerLastName          string        `json:"Owner Last Name"`
	OwnerAddress           string        `json:"Owner Address"`
	OwnerMobileNumber      string        `json:"Owner Mobile Number"`
	ChallanInfo            []ChallanInfo `json:"Challans"`
}

func InterfaceToModel(data interface{}) (instance *ChallanInfo, err error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return instance, err
	}

	return instance, json.Unmarshal(bytes, &instance)
}

func (p *ChallanInfo) PrimaryKey() string {
	return "Vehicle Number"
}

func (p *ChallanInfo) SortKey() string {
	return "Challan ID"
}

func (p *ChallanInfo) TableName() string {
	return "vehicle-challans"
}

func (p *ChallanInfo) ProjectionFields() []string {
	return []string{"Challan ID", "Unit Name", "Date", "Time", "Place of Violation", "PS Limits", "Violation", "Fine Amount"}
}

func (p *ChallanInfo) Bytes() ([]byte, error) {
	return json.Marshal(p)
}

func (p *ChallanInfo) GetMap() map[string]interface{} {
	return map[string]interface{}{
		"Vehicle Number":     p.VehicleRegNumber,
		"Challan ID":         p.ChallanID,
		"Unit Name":          p.UnitName,
		"Date":               p.Date,
		"Time":               p.Time,
		"Place of Violation": p.PlaceOfViolation,
		"PS Limits":          p.PSLimits,
		"Violation":          p.Violation,
		"Fine Amount":        p.FineAmount,
		"Created At":         p.CreatedAt,
	}
}

func ConvertIoReaderToStruct(data io.Reader, model interface{}) (interface{}, error) {
	if data == nil {
		return nil, errors.New("body is invalid")
	}
	return model, json.NewDecoder(data).Decode(model)
}
