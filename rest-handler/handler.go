package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/thegoodparticle/challan-data-svc/controller"
	grpcclient "github.com/thegoodparticle/challan-data-svc/grpc-client"
	"github.com/thegoodparticle/challan-data-svc/models"
	"github.com/thegoodparticle/challan-data-svc/store"
	HttpStatus "github.com/thegoodparticle/challan-data-svc/utils"
)

type Handler struct {
	Controller controller.Interface
	grpcSvc    *grpcclient.GRPCClient
}

func NewHandler(repository store.Interface, grpcClient *grpcclient.GRPCClient) *Handler {
	return &Handler{
		Controller: controller.NewController(repository),
		grpcSvc:    grpcClient,
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	if chi.URLParam(r, "RegID") != "" {
		h.getOne(w, r)
	} else {
		h.getAll(w, r)
	}
}

func (h *Handler) getOne(w http.ResponseWriter, r *http.Request) {
	ID := chi.URLParam(r, "RegID")

	challanInfo, err := h.Controller.ListOne(ID)
	if err != nil {
		HttpStatus.StatusInternalServerError(w, r, err)
		return
	}

	response := models.ChallanResponse{
		VehicleRegNumber: challanInfo.VehicleRegNumber,
		TotalFine:        challanInfo.TotalFine,
		Violations:       challanInfo.Violations,
		CreatedAt:        challanInfo.CreatedAt,
		UpdatedAt:        challanInfo.UpdatedAt,
	}

	vehicleInfo := h.grpcSvc.GetVehicleDetailsByRegistrationNumber(ID)
	if vehicleInfo != nil {
		response.VehicleModel = vehicleInfo.VehicleModel
		response.VehicleCompany = vehicleInfo.Company
		response.RegistrationDate = vehicleInfo.RegistrationDate.AsTime()
	}

	driverInfo := h.grpcSvc.GetOwnerDetailsByLicenseNumber(vehicleInfo.OwnerLicenseNumber)
	if driverInfo != nil {
		response.LicenseNumber = driverInfo.LicenseNumber
		response.LicenseOwnerName = driverInfo.DriverName
		response.LicenseDate = driverInfo.LicenseDate.AsTime()
	}

	HttpStatus.StatusOK(w, r, response)
}

func (h *Handler) getAll(w http.ResponseWriter, r *http.Request) {
	response, err := h.Controller.ListAll()
	if err != nil {
		HttpStatus.StatusInternalServerError(w, r, err)
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

	ID, err := h.Controller.Create(vehicleViolationsBody)
	if err != nil {
		HttpStatus.StatusInternalServerError(w, r, err)
		return
	}

	HttpStatus.StatusOK(w, r, map[string]interface{}{"registration_id": ID})
}

func (h *Handler) Put(w http.ResponseWriter, r *http.Request) {
	ID := chi.URLParam(r, "RegID")

	vehicleViolationsBody, err := h.getBodyAndValidate(r)
	if err != nil {
		HttpStatus.StatusBadRequest(w, r, err)
		return
	}

	if err := h.Controller.Update(ID, vehicleViolationsBody); err != nil {
		HttpStatus.StatusInternalServerError(w, r, err)
		return
	}

	HttpStatus.StatusNoContent(w, r)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	ID := chi.URLParam(r, "RegID")

	if err := h.Controller.Remove(ID); err != nil {
		HttpStatus.StatusInternalServerError(w, r, err)
		return
	}

	HttpStatus.StatusNoContent(w, r)
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	HttpStatus.StatusOK(w, r, "UP")
}

func (h *Handler) Options(w http.ResponseWriter, r *http.Request) {
	HttpStatus.StatusNoContent(w, r)
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
