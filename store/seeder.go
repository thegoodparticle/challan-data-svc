package store

import (
	"log"
	"time"

	"github.com/thegoodparticle/challan-data-svc/models"
)

var challans []models.ChallanInfo = []models.ChallanInfo{
	{
		VehicleRegNumber: "KA20AB1234",
		TotalFine:        1200,
		Violations: []map[string]interface{}{
			{
				"without-helmet": 500,
			},
			{
				"dangerous-driving": 500,
			},
			{
				"pollution-cert": 200,
			},
		},
		CreatedAt: time.Date(2011, 5, 17, 20, 34, 58, 651387237, time.UTC),
		UpdatedAt: time.Date(2024, 5, 17, 20, 34, 58, 651387237, time.UTC),
	},
	{
		VehicleRegNumber: "KA20CD5678",
		TotalFine:        1000,
		Violations: []map[string]interface{}{
			{
				"signal-jump": 1000,
			},
		},
		CreatedAt: time.Date(2023, 1, 12, 15, 20, 52, 561387237, time.UTC),
		UpdatedAt: time.Date(2023, 1, 12, 15, 20, 52, 561387237, time.UTC),
	},
}

func (db *Database) LoadDataIntoTables() {
	for _, violation := range challans {
		_, err := db.CreateOrUpdate(violation.GetMap(), violation.TableName())
		if err != nil {
			log.Printf("error while adding entry into the table. Error - %+v", err)
		}
	}
}
