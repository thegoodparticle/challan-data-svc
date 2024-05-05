package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	grpcclient "github.com/thegoodparticle/challan-data-svc/grpc-client"
	"github.com/thegoodparticle/challan-data-svc/internal/models"
	"github.com/thegoodparticle/challan-data-svc/internal/store"
	HttpStatus "github.com/thegoodparticle/challan-data-svc/internal/utils"
)

type Handler struct {
	controller *store.Store
	grpcSvc    *grpcclient.GRPCClient
}

func NewHandler(controller *store.Store, grpcClient *grpcclient.GRPCClient) *Handler {
	return &Handler{
		controller: controller,
		grpcSvc:    grpcClient,
	}
}

func (h *Handler) GetChallanResponseForRegistrationID(w http.ResponseWriter, r *http.Request) {
	ID := chi.URLParam(r, "RegID")
	if ID == "" {
		HttpStatus.StatusBadRequest(w, r, errors.New("vehicle registration number not provided"))
		return
	}

	challanInfo, err := h.controller.ListAllChallansForVehicleNumber(&models.ChallanInfo{VehicleRegNumber: ID})
	if err != nil {
		HttpStatus.StatusInternalServerError(w, r, err)
		return
	}

	response := models.ChallanResponse{
		VehicleRegNumber: ID,
		ChallanInfo:      challanInfo,
	}

	vehicleInfo := h.grpcSvc.GetVehicleDetailsByRegistrationNumber(ID)
	if vehicleInfo != nil {
		response.VehicleModel = vehicleInfo.VehicleModel
		response.VehicleCompany = vehicleInfo.Company
		response.RegistrationDate = vehicleInfo.RegistrationDate.AsTime()
		response.RegistrationExpiryDate = vehicleInfo.RegistrationExpiryDate.AsTime()
	}

	ownerInfo := h.grpcSvc.GetOwnerDetailsByID(vehicleInfo.OwnerID)
	if ownerInfo != nil {
		response.OwnerFirstName = ownerInfo.OwnerFirstName
		response.OwnerLastName = ownerInfo.OwnerLastName
		response.OwnerMobileNumber = ownerInfo.MobileNumber
		response.OwnerAddress = ownerInfo.Address
	}

	if response.VehicleModel == "" && response.OwnerFirstName == "" {
		HttpStatus.StatusNotFound(w, r, errors.New("vehicle details not found"))
		return
	}

	HttpStatus.StatusOK(w, r, response)
}

func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	vehicleViolationsBody, err := h.getBodyAndValidate(r)
	if err != nil {
		HttpStatus.StatusBadRequest(w, r, err)
		return
	}

	ID, err := h.controller.Create(vehicleViolationsBody)
	if err != nil {
		HttpStatus.StatusInternalServerError(w, r, err)
		return
	}

	HttpStatus.StatusOK(w, r, map[string]interface{}{"registration_id": ID})
}

func (h *Handler) getBodyAndValidate(r *http.Request) (*models.ChallanInfo, error) {
	vehicleViolationsBody := &models.ChallanInfo{}
	body, err := models.ConvertIoReaderToStruct(r.Body, vehicleViolationsBody)
	if err != nil {
		return &models.ChallanInfo{}, errors.New("body is required")
	}

	vehicleBodyParsed, err := models.InterfaceToModel(body)
	if err != nil {
		return &models.ChallanInfo{}, errors.New("error on convert body to model")
	}

	if vehicleBodyParsed.VehicleRegNumber == "" {
		return &models.ChallanInfo{}, errors.New("registration ID is required")
	}

	log.Printf("successful parse of request body. %+v", vehicleBodyParsed)

	return vehicleBodyParsed, nil
}
