package controller

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/thegoodparticle/challan-data-svc/models"
	"github.com/thegoodparticle/challan-data-svc/store"
)

type Controller struct {
	repository store.Interface
}

type Interface interface {
	ListOne(ID string) (entity models.ChallanInfo, err error)
	ListAll() (entities []models.ChallanInfo, err error)
	Create(entity *models.ChallanInfo) (string, error)
	Update(ID string, entity *models.ChallanInfo) error
	Remove(ID string) error
}

func NewController(repository store.Interface) Interface {
	return &Controller{repository: repository}
}

func (c *Controller) ListOne(id string) (entity models.ChallanInfo, err error) {
	entity.VehicleRegNumber = id
	log.Println(id)
	response, err := c.repository.FindOne(entity.GetFilterId(), entity.TableName())
	if err != nil {
		return entity, err
	}
	return models.ParseDynamoAtributeToStruct(response.Item)
}

func (c *Controller) ListAll() (entities []models.ChallanInfo, err error) {
	entities = []models.ChallanInfo{}
	var entity models.ChallanInfo

	filter := expression.Name("owner_name").NotEqual(expression.Value(""))
	condition, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		return entities, err
	}

	response, err := c.repository.FindAll(condition, entity.TableName())
	if err != nil {
		return entities, err
	}

	if response != nil {
		for _, value := range response.Items {
			entity, err := models.ParseDynamoAtributeToStruct(value)
			if err != nil {
				return entities, err
			}
			entities = append(entities, entity)
		}
	}

	return entities, nil
}

func (c *Controller) Create(entity *models.ChallanInfo) (string, error) {
	entity.CreatedAt = time.Now()
	entity.UpdatedAt = time.Now()
	_, err := c.repository.CreateOrUpdate(entity.GetMap(), entity.TableName())
	return entity.VehicleRegNumber, err
}

func (c *Controller) Update(id string, entity *models.ChallanInfo) error {
	found, err := c.ListOne(id)
	if err != nil {
		return err
	}

	found.VehicleRegNumber = id

	found.TotalFine = entity.TotalFine

	if entity.Violations != nil {
		found.Violations = entity.Violations
	}

	found.UpdatedAt = time.Now()
	_, err = c.repository.CreateOrUpdate(found.GetMap(), entity.TableName())
	return err
}

func (c *Controller) Remove(id string) error {
	entity, err := c.ListOne(id)
	if err != nil {
		return err
	}
	_, err = c.repository.Delete(entity.GetFilterId(), entity.TableName())
	return err
}
